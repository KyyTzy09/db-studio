package wizard

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"

	"db-studio-go/internal/config"
)

// RunCLIWizard presents an interactive CLI wizard powered by Charm Huh
func RunCLIWizard(projectPath string) (*config.ConnectionConfig, error) {
	var driverStr string
	var host string = "localhost"
	var portStr string
	var user string
	var password string
	var dbName string
	var filePath string = "./local.db"

	// Step 1: Select Database Driver
	driverSelectForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Database Type").
				Options(
					huh.NewOption("PostgreSQL", "postgres"),
					huh.NewOption("MySQL / MariaDB", "mysql"),
					huh.NewOption("SQLite", "sqlite"),
				).
				Value(&driverStr),
		),
	)

	if err := driverSelectForm.Run(); err != nil {
		return nil, err
	}

	driver := config.DriverType(driverStr)
	conn := config.ConnectionConfig{
		Driver:      driver,
		ProjectPath: projectPath,
		Name:        filepath.Base(projectPath) + " (" + string(driver) + ")",
	}

	// SQLite Form
	if driver == config.DriverSQLite {
		sqliteForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("SQLite File Path").
					Placeholder("./dev.db or database.sqlite").
					Value(&filePath),
			),
		)
		if err := sqliteForm.Run(); err != nil {
			return nil, err
		}
		if strings.TrimSpace(filePath) == "" {
			filePath = "./local.db"
		}
		conn.FilePath = strings.TrimSpace(filePath)
		return &conn, nil
	}

	// Server Database Form (PostgreSQL / MySQL)
	defaultPort := "5432"
	defaultUser := "postgres"
	if driver == config.DriverMySQL {
		defaultPort = "3306"
		defaultUser = "root"
	}
	portStr = defaultPort
	user = defaultUser

	serverForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Database Host").
				Placeholder("localhost").
				Value(&host),
			huh.NewInput().
				Title("Database Port").
				Placeholder(defaultPort).
				Value(&portStr),
			huh.NewInput().
				Title("Database User").
				Placeholder(defaultUser).
				Value(&user),
			huh.NewInput().
				Title("Database Password").
				EchoMode(huh.EchoModePassword).
				Value(&password),
			huh.NewInput().
				Title("Database Name").
				Placeholder("myapp_db").
				Value(&dbName),
		),
	)

	if err := serverForm.Run(); err != nil {
		return nil, err
	}

	portInt, _ := strconv.Atoi(strings.TrimSpace(portStr))
	if portInt == 0 {
		if driver == config.DriverMySQL {
			portInt = 3306
		} else {
			portInt = 5432
		}
	}

	conn.Host = strings.TrimSpace(host)
	conn.Port = portInt
	conn.User = strings.TrimSpace(user)
	conn.Password = password
	conn.Database = strings.TrimSpace(dbName)

	return &conn, nil
}
