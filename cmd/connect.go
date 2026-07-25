package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"db-studio-go/internal/config"
	"db-studio-go/internal/scanner"
	"db-studio-go/internal/ui"
	"db-studio-go/internal/wizard"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Configure or select database connection for current project",
	Long:  `Scan project for database configurations (.env, docker-compose) or launch CLI wizard to set up database parameters.`,
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
		ui.PrintScanning()

		compositeScanner := scanner.NewCompositeScanner(
			scanner.NewEnvScanner(),
			scanner.NewDockerComposeScanner(),
		)

		detectedConns, _ := compositeScanner.Scan(context.Background(), cwd)
		reader := bufio.NewReader(os.Stdin)

		var selectedConn *config.ConnectionConfig

		if len(detectedConns) > 0 {
			ui.PrintInfo("Auto-detected connections:")

			for i, conn := range detectedConns {
				if conn.Driver == config.DriverSQLite {
					fmt.Printf("   [%d] %s ➔ %s (%s)\n", i+1, ui.BoldText(conn.Name), conn.FilePath, conn.Driver)
				} else {
					fmt.Printf("   [%d] %s ➔ %s:%d/%s (%s)\n", i+1, ui.BoldText(conn.Name), conn.Host, conn.Port, conn.Database, conn.Driver)
				}
			}

			manualChoiceIdx := len(detectedConns) + 1
			fmt.Printf("   [%d] %s\n", manualChoiceIdx, ui.Gray("Manual Connection Setup (Wizard)"))
			fmt.Printf("\nSelect [1-%d]: ", manualChoiceIdx)

			choiceStr, _ := reader.ReadString('\n')
			choiceStr = strings.TrimSpace(choiceStr)
			choiceIdx, _ := strconv.Atoi(choiceStr)

			if choiceIdx >= 1 && choiceIdx <= len(detectedConns) {
				selectedConn = &detectedConns[choiceIdx-1]
			}
		}

		if selectedConn == nil {
			wizardConn, err := wizard.RunCLIWizard(cwd)
			if err != nil {
				return fmt.Errorf("wizard error: %w", err)
			}
			selectedConn = wizardConn
		}

		if err := configMgr.SaveConnection(*selectedConn); err != nil {
			return fmt.Errorf("failed to save connection: %w", err)
		}

		ui.PrintSuccess(fmt.Sprintf("Connection saved: %s (%s)", selectedConn.Name, selectedConn.Driver))
		ui.PrintReady()
		return nil
	},
}
