package wizard

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"db-studio-go/internal/config"
)

// RunCLIWizard presents an interactive CLI prompt to input database connection details
func RunCLIWizard(projectPath string) (*config.ConnectionConfig, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n==================================================")
	fmt.Println(" 🛠️  DBStudio Connection Wizard")
	fmt.Println(" Tidak ada konfigurasi database yang terdeteksi.")
	fmt.Println(" Silakan masukkan rincian koneksi database Anda.")
	fmt.Println("==================================================")

	fmt.Println("\nPilih Tipe Database:")
	fmt.Println(" [1] PostgreSQL (Default)")
	fmt.Println(" [2] MySQL / MariaDB")
	fmt.Println(" [3] SQLite")
	fmt.Print("Pilihan [1-3]: ")

	driverChoice, _ := reader.ReadString('\n')
	driverChoice = strings.TrimSpace(driverChoice)

	driver := config.DriverPostgres
	defaultPort := 5432

	switch driverChoice {
	case "2":
		driver = config.DriverMySQL
		defaultPort = 3306
	case "3":
		driver = config.DriverSQLite
	}

	conn := config.ConnectionConfig{
		Driver:      driver,
		ProjectPath: projectPath,
		Name:        filepath.Base(projectPath) + " (" + string(driver) + ")",
	}

	if driver == config.DriverSQLite {
		fmt.Print("\nMasukkan Path File SQLite (e.g. ./dev.db atau database.sqlite): ")
		filePath, _ := reader.ReadString('\n')
		conn.FilePath = strings.TrimSpace(filePath)
		if conn.FilePath == "" {
			conn.FilePath = "./local.db"
		}
		return &conn, nil
	}

	// Host
	fmt.Printf("\nHost Database [default: localhost]: ")
	host, _ := reader.ReadString('\n')
	host = strings.TrimSpace(host)
	if host == "" {
		host = "localhost"
	}
	conn.Host = host

	// Port
	fmt.Printf("Port Database [default: %d]: ", defaultPort)
	portStr, _ := reader.ReadString('\n')
	portStr = strings.TrimSpace(portStr)
	port, _ := strconv.Atoi(portStr)
	if port == 0 {
		port = defaultPort
	}
	conn.Port = port

	// User
	defaultUser := "postgres"
	if driver == config.DriverMySQL {
		defaultUser = "root"
	}
	fmt.Printf("Username Database [default: %s]: ", defaultUser)
	user, _ := reader.ReadString('\n')
	user = strings.TrimSpace(user)
	if user == "" {
		user = defaultUser
	}
	conn.User = user

	// Password
	fmt.Print("Password Database: ")
	pass, _ := reader.ReadString('\n')
	conn.Password = strings.TrimSpace(pass)

	// Database Name
	fmt.Print("Nama Database: ")
	dbname, _ := reader.ReadString('\n')
	conn.Database = strings.TrimSpace(dbname)

	fmt.Println("\n✅ Konfigurasi database berhasil disimpan!")
	return &conn, nil
}
