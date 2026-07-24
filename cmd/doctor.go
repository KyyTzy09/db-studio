package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"db-studio-go/internal/config"
	"db-studio-go/internal/db"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose database connection health and configuration",
	Long:  `Checks saved database connections for the current project, attempts a ping test, and prints diagnostic health metrics.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		configMgr, err := config.NewManager()
		if err != nil {
			return fmt.Errorf("config manager error: %w", err)
		}

		fmt.Println("🏥 DBStudio Doctor - Diagnostic Health Check")
		fmt.Printf("📂 Project Directory: %s\n", cwd)
		fmt.Printf("📄 Config File Path: %s\n\n", configMgr.GetConfigFilePath())

		conn, found, err := configMgr.FindByProjectPath(cwd)
		if err != nil {
			return fmt.Errorf("error reading config: %w", err)
		}

		if !found || conn == nil {
			fmt.Println("❌ Status: No saved connection found for this project.")
			fmt.Println("👉 Run 'dbstudio connect' or 'dbstudio' to configure connection.")
			return nil
		}

		fmt.Println("==================================================")
		fmt.Printf(" Connection Name : %s\n", conn.Name)
		fmt.Printf(" Driver Type     : %s\n", conn.Driver)
		if conn.Driver == config.DriverSQLite {
			fmt.Printf(" File Path       : %s\n", conn.FilePath)
		} else {
			fmt.Printf(" Host & Port     : %s:%d\n", conn.Host, conn.Port)
			fmt.Printf(" Username        : %s\n", conn.User)
			fmt.Printf(" Database Name   : %s\n", conn.Database)
		}
		fmt.Println("==================================================")

		fmt.Print("⏳ Testing Connection Ping... ")

		driverInst, err := db.NewDriver(*conn)
		if err != nil {
			fmt.Printf("❌ FAILED (%v)\n", err)
			return nil
		}
		defer driverInst.Disconnect()

		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := driverInst.Ping(ctx); err != nil {
			fmt.Printf("❌ FAILED\n")
			fmt.Printf("⚠️  Error Details: %v\n", err)
		} else {
			latency := time.Since(start).Milliseconds()
			fmt.Printf("✅ ONLINE (%d ms latency)\n", latency)

			tables, tErr := driverInst.GetTables(ctx)
			if tErr == nil {
				fmt.Printf("📊 Detected Tables: %d table(s) in schema\n", len(tables))
			}
		}

		return nil
	},
}
