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

		fmt.Println("🔍 DBStudio: Scanning project for database configurations...")

		compositeScanner := scanner.NewCompositeScanner(
			scanner.NewEnvScanner(),
			scanner.NewDockerComposeScanner(),
		)

		detectedConns, _ := compositeScanner.Scan(context.Background(), cwd)
		reader := bufio.NewReader(os.Stdin)

		var selectedConn *config.ConnectionConfig

		if len(detectedConns) > 0 {
			fmt.Println("\n==================================================")
			fmt.Println(" 🛠️  DBStudio Connection Manager")
			fmt.Println("==================================================")
			fmt.Println("Ditemukan koneksi otomatis di proyek ini:")

			for i, conn := range detectedConns {
				if conn.Driver == config.DriverSQLite {
					fmt.Printf(" [%d] %s ➔ %s\n", i+1, conn.Name, conn.FilePath)
				} else {
					fmt.Printf(" [%d] %s ➔ %s:%d (Database: %s)\n", i+1, conn.Name, conn.Host, conn.Port, conn.Database)
				}
			}

			manualChoiceIdx := len(detectedConns) + 1
			fmt.Printf(" [%d] Input Koneksi Manual (Wizard Baru)\n", manualChoiceIdx)
			fmt.Printf("\nPilihan [1-%d]: ", manualChoiceIdx)

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

		fmt.Printf("\n✅ Koneksi '%s' (%s) berhasil disimpan!\n", selectedConn.Name, selectedConn.Driver)
		fmt.Printf("📂 Global config updated at: %s\n", configMgr.GetConfigFilePath())
		return nil
	},
}
