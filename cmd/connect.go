package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"db-studio-go/internal/config"
	"db-studio-go/internal/wizard"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Launch interactive CLI wizard to configure database connection",
	Long:  `Directly prompts the interactive CLI wizard to set up or edit database connection parameters for the current project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		configMgr, err := config.NewManager()
		if err != nil {
			return fmt.Errorf("config manager error: %w", err)
		}

		conn, err := wizard.RunCLIWizard(cwd)
		if err != nil {
			return fmt.Errorf("wizard error: %w", err)
		}

		if err := configMgr.SaveConnection(*conn); err != nil {
			return fmt.Errorf("failed to save connection: %w", err)
		}

		fmt.Printf("✅ Connection for project '%s' saved to %s\n", conn.Name, configMgr.GetConfigFilePath())
		return nil
	},
}
