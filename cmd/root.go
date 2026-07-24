package cmd

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"db-studio-go/internal/config"
	"db-studio-go/internal/db"
	httpServer "db-studio-go/internal/http"
	"db-studio-go/internal/scanner"
	"db-studio-go/internal/wizard"
)

var (
	portFlag int
	WebFS    embed.FS
)

func NewRootCmd(webFS embed.FS) *cobra.Command {
	WebFS = webFS

	rootCmd := &cobra.Command{
		Use:   "dbstudio",
		Short: "One command database studio for every developer",
		Long:  `DBStudio auto-detects database configurations in your project and launches a local web studio interface instantly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %w", err)
			}

			configMgr, err := config.NewManager()
			if err != nil {
				return fmt.Errorf("config manager error: %w", err)
			}

			fmt.Println("🔍 DBStudio: Checking database configurations...")

			// Step 1: Check saved config in global path (~/.config/dbstudio/connections.json)
			targetConn, found, err := configMgr.FindByProjectPath(cwd)
			if err != nil {
				fmt.Printf("⚠️ Warning reading global config: %v\n", err)
			}

			// Step 2: Auto-Detection Scanner if not found in global config
			if !found || targetConn == nil {
				fmt.Println("🔍 Running auto-detection scanners (.env & docker-compose)...")
				compositeScanner := scanner.NewCompositeScanner(
					scanner.NewEnvScanner(),
					scanner.NewDockerComposeScanner(),
				)

				detectedConns, _ := compositeScanner.Scan(context.Background(), cwd)

				if len(detectedConns) == 1 {
					fmt.Printf("✨ Auto-detected 1 connection: %s (%s)\n", detectedConns[0].Name, detectedConns[0].Driver)
					targetConn = &detectedConns[0]
					_ = configMgr.SaveConnection(*targetConn)
				} else if len(detectedConns) > 1 {
					fmt.Printf("✨ Auto-detected %d database connections. Selecting first: %s\n", len(detectedConns), detectedConns[0].Name)
					targetConn = &detectedConns[0]
					_ = configMgr.SaveConnection(*targetConn)
				} else {
					// Step 3: Fallback to CLI Interactive Wizard
					wizardConn, err := wizard.RunCLIWizard(cwd)
					if err != nil {
						return fmt.Errorf("wizard failed: %w", err)
					}
					targetConn = wizardConn
					_ = configMgr.SaveConnection(*targetConn)
				}
			} else {
				fmt.Printf("⚡ Found saved connection for project: %s (%s)\n", targetConn.Name, targetConn.Driver)
			}

			// Step 4: Instantiate Database Driver (Lazy Connection)
			driver, err := db.NewDriver(*targetConn)
			if err != nil {
				return fmt.Errorf("failed to initialize driver: %w", err)
			}

			// Step 5: Start Chi HTTP Server (Automatic Port Fallback) & Open Browser
			srv := httpServer.NewServer(driver, WebFS, portFlag)

			return srv.ListenAndServe(func(actualPort int) {
				go func() {
					url := fmt.Sprintf("http://localhost:%d", actualPort)
					fmt.Printf("🌐 Opening Web Studio at %s...\n", url)
					openBrowser(url)
				}()
			})
		},
	}

	rootCmd.Flags().IntVarP(&portFlag, "port", "p", 8080, "Port to run the HTTP web studio server on")

	// Register Subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(doctorCmd)

	return rootCmd
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	}
	if err != nil {
		fmt.Printf("⚠️ Could not open browser automatically: %v\n", err)
	}
}
