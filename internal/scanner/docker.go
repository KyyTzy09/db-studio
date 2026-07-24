package scanner

import (
	"bufio"
	"context"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"db-studio-go/internal/config"
)

// DockerComposeScanner scans docker-compose.yml or compose.yaml files in project path
type DockerComposeScanner struct{}

func NewDockerComposeScanner() *DockerComposeScanner {
	return &DockerComposeScanner{}
}

func (s *DockerComposeScanner) Name() string {
	return "Docker Compose Scanner"
}

func (s *DockerComposeScanner) Scan(ctx context.Context, projectPath string) ([]config.ConnectionConfig, error) {
	composeFiles := []string{"docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"}
	var detected []config.ConnectionConfig

	for _, filename := range composeFiles {
		filePath := filepath.Join(projectPath, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}

		conns, err := parseDockerComposeFile(filePath, projectPath, filename)
		if err == nil && len(conns) > 0 {
			detected = append(detected, conns...)
		}
	}

	return detected, nil
}

func parseDockerComposeFile(filePath, projectPath, sourceFilename string) ([]config.ConnectionConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var detected []config.ConnectionConfig
	scanner := bufio.NewScanner(file)

	var currentImage, currentService string
	envVars := make(map[string]string)
	var hostPort int
	inServicesBlock := false

	flushService := func() {
		if currentService == "" {
			return
		}

		driver := inferDockerDriver(currentImage, currentService, envVars)
		if driver != "" {
			host := "localhost"
			port := hostPort
			if port == 0 {
				if driver == config.DriverPostgres {
					port = 5432
				} else if driver == config.DriverMySQL {
					port = 3306
				}
			}

			user := getEnvVal(envVars, "POSTGRES_USER", "MYSQL_USER", "DB_USER")
			pass := getEnvVal(envVars, "POSTGRES_PASSWORD", "MYSQL_ROOT_PASSWORD", "MYSQL_PASSWORD", "DB_PASS")
			dbname := getEnvVal(envVars, "POSTGRES_DB", "MYSQL_DATABASE", "DB_NAME")

			if user == "" && driver == config.DriverPostgres {
				user = "postgres"
			}
			if user == "" && driver == config.DriverMySQL {
				user = "root"
			}

			conn := config.ConnectionConfig{
				Name:        sourceFilename + " (" + currentService + ")",
				Driver:      driver,
				Host:        host,
				Port:        port,
				User:        user,
				Password:    pass,
				Database:    dbname,
				ProjectPath: projectPath,
			}
			detected = append(detected, conn)
		}

		currentImage = ""
		currentService = ""
		envVars = make(map[string]string)
		hostPort = 0
	}

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " \t"))

		if trimmed == "services:" {
			inServicesBlock = true
			continue
		}

		if inServicesBlock {
			// Top-level service under services: (indent level 2 or 4)
			if strings.HasSuffix(trimmed, ":") && indent <= 4 && !strings.Contains(trimmed, " ") {
				svcName := strings.TrimSuffix(trimmed, ":")
				if svcName != "environment" && svcName != "ports" && svcName != "volumes" && svcName != "networks" {
					flushService()
					currentService = svcName
					continue
				}
			}

			if currentService != "" {
				// Detect image
				if strings.HasPrefix(trimmed, "image:") {
					parts := strings.SplitN(trimmed, ":", 2)
					if len(parts) == 2 {
						currentImage = strings.TrimSpace(parts[1])
					}
				}

				// Detect ports (e.g. "- 5432:5432" or "- '5432:5432'")
				if strings.HasPrefix(trimmed, "-") && strings.Contains(trimmed, ":") {
					cleanPort := strings.TrimPrefix(trimmed, "-")
					cleanPort = strings.Trim(strings.TrimSpace(cleanPort), `"'`)
					portParts := strings.Split(cleanPort, ":")
					if len(portParts) >= 2 {
						p, _ := strconv.Atoi(strings.TrimSpace(portParts[0]))
						if p > 0 {
							hostPort = p
						}
					}
				}

				// Detect environment variables (e.g. "POSTGRES_USER: postgres" or "POSTGRES_PASSWORD=secret")
				if strings.Contains(trimmed, ":") || strings.Contains(trimmed, "=") {
					var k, v string
					if strings.Contains(trimmed, "=") {
						parts := strings.SplitN(trimmed, "=", 2)
						k, v = strings.TrimPrefix(parts[0], "-"), parts[1]
					} else {
						parts := strings.SplitN(trimmed, ":", 2)
						k, v = strings.TrimPrefix(parts[0], "-"), parts[1]
					}
					k = strings.TrimSpace(k)
					v = strings.Trim(strings.TrimSpace(v), `"'`)
					if k != "" && v != "" && k != "image" && k != "ports" && k != "command" && k != "volumes" {
						envVars[k] = v
					}
				}
			}
		}
	}

	flushService()
	return detected, nil
}

func inferDockerDriver(image, service string, env map[string]string) config.DriverType {
	combined := strings.ToLower(image + " " + service)
	if strings.Contains(combined, "postgres") || strings.Contains(combined, "pgsql") {
		return config.DriverPostgres
	}
	if strings.Contains(combined, "mysql") || strings.Contains(combined, "mariadb") {
		return config.DriverMySQL
	}
	if getEnvVal(env, "POSTGRES_USER", "POSTGRES_DB") != "" {
		return config.DriverPostgres
	}
	if getEnvVal(env, "MYSQL_ROOT_PASSWORD", "MYSQL_DATABASE") != "" {
		return config.DriverMySQL
	}
	return ""
}
