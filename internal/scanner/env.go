package scanner

import (
	"bufio"
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"db-studio-go/internal/config"
)

// EnvScanner scans .env files in the project path
type EnvScanner struct{}

func NewEnvScanner() *EnvScanner {
	return &EnvScanner{}
}

func (s *EnvScanner) Name() string {
	return "DotEnv Scanner"
}

func (s *EnvScanner) Scan(ctx context.Context, projectPath string) ([]config.ConnectionConfig, error) {
	envFiles := []string{".env", ".env.local", ".env.development"}
	var detected []config.ConnectionConfig

	for _, filename := range envFiles {
		envPath := filepath.Join(projectPath, filename)
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			continue
		}

		kvMap, err := parseEnvFile(envPath)
		if err != nil {
			continue
		}

		// 1. Try DATABASE_URL / DB_URL connection string parsing
		for _, urlKey := range []string{"DATABASE_URL", "DB_URL", "POSTGRES_URL", "MYSQL_URL"} {
			if rawURL, ok := kvMap[urlKey]; ok && rawURL != "" {
				if conn, ok := parseConnectionURL(rawURL, projectPath); ok {
					conn.Name = filename + " (" + string(conn.Driver) + ")"
					detected = append(detected, conn)
				}
			}
		}

		// 2. Try individual key-value pairs
		driver := inferDriver(kvMap)
		if driver != "" {
			host := getEnvVal(kvMap, "DB_HOST", "POSTGRES_HOST", "MYSQL_HOST", "DATABASE_HOST")
			if host == "" {
				host = "localhost"
			}

			portStr := getEnvVal(kvMap, "DB_PORT", "POSTGRES_PORT", "MYSQL_PORT")
			port, _ := strconv.Atoi(portStr)
			if port == 0 {
				if driver == config.DriverPostgres {
					port = 5432
				} else if driver == config.DriverMySQL {
					port = 3306
				}
			}

			user := getEnvVal(kvMap, "DB_USER", "DB_USERNAME", "POSTGRES_USER", "MYSQL_USER")
			pass := getEnvVal(kvMap, "DB_PASS", "DB_PASSWORD", "POSTGRES_PASSWORD", "MYSQL_PASSWORD")
			dbname := getEnvVal(kvMap, "DB_NAME", "DB_DATABASE", "POSTGRES_DB", "MYSQL_DATABASE")

			if dbname != "" || driver == config.DriverSQLite {
				conn := config.ConnectionConfig{
					Name:        filename + " (" + string(driver) + ")",
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
		}
	}

	return detected, nil
}

func parseEnvFile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	kv := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			kv[key] = val
		}
	}

	return kv, nil
}

func parseConnectionURL(rawURL, projectPath string) (config.ConnectionConfig, bool) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return config.ConnectionConfig{}, false
	}

	var driver config.DriverType
	switch strings.ToLower(u.Scheme) {
	case "postgres", "postgresql":
		driver = config.DriverPostgres
	case "mysql":
		driver = config.DriverMySQL
	case "file", "sqlite", "sqlite3":
		driver = config.DriverSQLite
	default:
		return config.ConnectionConfig{}, false
	}

	host := u.Hostname()
	if host == "" {
		host = "localhost"
	}

	port, _ := strconv.Atoi(u.Port())
	if port == 0 {
		if driver == config.DriverPostgres {
			port = 5432
		} else if driver == config.DriverMySQL {
			port = 3306
		}
	}

	user := u.User.Username()
	pass, _ := u.User.Password()
	dbname := strings.TrimPrefix(u.Path, "/")

	return config.ConnectionConfig{
		Driver:      driver,
		Host:        host,
		Port:        port,
		User:        user,
		Password:    pass,
		Database:    dbname,
		ProjectPath: projectPath,
	}, true
}

func inferDriver(kv map[string]string) config.DriverType {
	driverVal := strings.ToLower(getEnvVal(kv, "DB_CONNECTION", "DB_DRIVER", "DATABASE_DRIVER"))
	if strings.Contains(driverVal, "postgres") || strings.Contains(driverVal, "pgsql") {
		return config.DriverPostgres
	}
	if strings.Contains(driverVal, "mysql") {
		return config.DriverMySQL
	}
	if strings.Contains(driverVal, "sqlite") {
		return config.DriverSQLite
	}

	// Fallback inference by key names
	if getEnvVal(kv, "POSTGRES_USER", "POSTGRES_DB") != "" {
		return config.DriverPostgres
	}
	if getEnvVal(kv, "MYSQL_USER", "MYSQL_DATABASE") != "" {
		return config.DriverMySQL
	}

	return ""
}

func getEnvVal(kv map[string]string, keys ...string) string {
	for _, k := range keys {
		if v, ok := kv[k]; ok && v != "" {
			return v
		}
	}
	return ""
}
