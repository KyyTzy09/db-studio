package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"db-studio-go/internal/config"
	"db-studio-go/internal/db"
	"db-studio-go/internal/ui"
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

		ui.PrintBanner(AppVersion)
		ui.PrintInfo(fmt.Sprintf("Diagnosing project: %s", ui.Gray(cwd)))

		conn, found, err := configMgr.FindByProjectPath(cwd)
		if err != nil {
			return fmt.Errorf("error reading config: %w", err)
		}

		if !found || conn == nil {
			ui.PrintError("No saved database connection found for this project.")
			fmt.Printf("   %s Run 'dbstudio connect' or 'dbstudio' to setup connection.\n\n", ui.Gray("└─"))
			return nil
		}

		fmt.Println()
		if conn.Driver == config.DriverSQLite {
			fmt.Printf("   %s %s: %s (%s)\n", ui.Blue("🗄"), ui.BoldText("Target DB"), conn.Name, conn.FilePath)
		} else {
			fmt.Printf("   %s %s: %s (%s@%s:%d/%s)\n",
				ui.Blue("🗄"),
				ui.BoldText("Target DB"),
				conn.Name,
				conn.Driver,
				conn.Host,
				conn.Port,
				conn.Database,
			)
		}

		driverInst, err := db.NewDriver(*conn)
		if err != nil {
			ui.PrintError(fmt.Sprintf("Driver initialization failed: %v", err))
			return nil
		}
		defer driverInst.Disconnect()

		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := driverInst.Ping(ctx); err != nil {
			ui.PrintError(fmt.Sprintf("Connection failed (%v)", err))
			fmt.Printf("   %s Check host availability, credentials, or network status.\n\n", ui.Gray("└─"))
		} else {
			latency := time.Since(start).Milliseconds()
			ui.PrintSuccess(fmt.Sprintf("Connected to %s %s", conn.Driver, ui.Gray(fmt.Sprintf("(%d ms)", latency))))

			tables, tErr := driverInst.GetTables(ctx)
			if tErr == nil {
				ui.PrintSuccess(fmt.Sprintf("%d table(s) detected in schema", len(tables)))
			}
			ui.PrintReady()
		}

		return nil
	},
}
