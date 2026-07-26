package cmd

import (
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"db-studio-go/internal/config"
	"db-studio-go/internal/db"
	httpServer "db-studio-go/internal/http"
	"db-studio-go/internal/scanner"
	"db-studio-go/internal/ui"
	"db-studio-go/internal/wizard"
)

var (
	portFlag    int
	verboseFlag bool
	WebFS       embed.FS
)

func NewRootCmd(webFS embed.FS) *cobra.Command {
	WebFS = webFS

	rootCmd := &cobra.Command{
		Use:   "dbstudio",
		Short: "One command database studio for every developer",
		Long:  `DBStudio auto-detects database configurations in your project and launches a local web studio interface instantly.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verboseFlag {
				ui.Logger.SetLevel(log.DebugLevel)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current working directory: %w", err)
			}

			configMgr, err := config.NewManager()
			if err != nil {
				return fmt.Errorf("config manager error: %w", err)
			}

			// Print Banner
			ui.PrintBanner(AppVersion)
			ui.PrintScanning()

			// Step 1: Check saved config in global path (~/.config/dbstudio/connections.json)
			targetConn, found, err := configMgr.FindByProjectPath(cwd)
			if err != nil {
				ui.PrintWarning(fmt.Sprintf("Global config warning: %v", err))
			}

			// Step 2: Auto-Detection Scanner if not found in global config
			if !found || targetConn == nil {
				compositeScanner := scanner.NewCompositeScanner(
					scanner.NewEnvScanner(),
					scanner.NewDockerComposeScanner(),
				)

				detectedConns, _ := compositeScanner.Scan(context.Background(), cwd)

				if len(detectedConns) == 1 {
					ui.PrintSuccess(fmt.Sprintf("Auto-detected 1 connection: %s (%s)", detectedConns[0].Name, detectedConns[0].Driver))
					targetConn = &detectedConns[0]
					_ = configMgr.SaveConnection(*targetConn)
				} else if len(detectedConns) > 1 {
					ui.PrintSuccess(fmt.Sprintf("Auto-detected %d database connections. Selecting %s", len(detectedConns), detectedConns[0].Name))
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
				ui.PrintSuccess(fmt.Sprintf("Found saved config: %s (%s)", targetConn.Name, targetConn.Driver))
			}

			ui.PrintSuccess("Connected")

			// Step 4: Instantiate Database Driver (Lazy Connection)
			driver, err := db.NewDriver(*targetConn)
			if err != nil {
				return fmt.Errorf("failed to initialize driver: %w", err)
			}

			ui.PrintStarting()

			// Step 5: Start Chi HTTP Server (Automatic Port Fallback) & Open Browser
			srv := httpServer.NewServer(driver, configMgr, WebFS, portFlag)

			return srv.ListenAndServe(func(actualPort int) {
				url := fmt.Sprintf("http://localhost:%d", actualPort)
				ui.PrintListening(url)
				ui.PrintOpeningBrowser()
				ui.PrintReady()

				go openBrowser(url)
			})
		},
	}

	rootCmd.PersistentFlags().IntVarP(&portFlag, "port", "p", 8080, "Port for DBStudio Web HTTP Server")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable verbose output and logs")

	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(doctorCmd)
	rootCmd.AddCommand(versionCmd)

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
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		ui.PrintWarning(fmt.Sprintf("Could not open browser automatically: %v", err))
	}
}
