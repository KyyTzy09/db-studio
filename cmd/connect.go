package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
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

		var selectedConn *config.ConnectionConfig

		if len(detectedConns) > 0 {
			var options []huh.Option[int]
			for i, conn := range detectedConns {
				var label string
				if conn.Driver == config.DriverSQLite {
					label = fmt.Sprintf("%s → %s (%s)", conn.Name, conn.FilePath, conn.Driver)
				} else {
					label = fmt.Sprintf("%s → %s:%d/%s (%s)", conn.Name, conn.Host, conn.Port, conn.Database, conn.Driver)
				}
				options = append(options, huh.NewOption(label, i))
			}
			options = append(options, huh.NewOption("Manual Connection Setup (Wizard)", -1))

			var selectedIdx int
			connSelectForm := huh.NewForm(
				huh.NewGroup(
					huh.NewSelect[int]().
						Title("Select Database Connection").
						Options(options...).
						Value(&selectedIdx),
				),
			)

			if err := connSelectForm.Run(); err == nil && selectedIdx >= 0 && selectedIdx < len(detectedConns) {
				selectedConn = &detectedConns[selectedIdx]
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
