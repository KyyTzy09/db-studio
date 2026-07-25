package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"db-studio-go/internal/ui"
)

const AppVersion = "v0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of DBStudio",
	Long:  `Display DBStudio build version, Go runtime version, and operating system architecture.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner(AppVersion)
		fmt.Printf("%s DBStudio %s %s\n",
			ui.Green("✔"),
			ui.BoldText(AppVersion),
			ui.Gray(fmt.Sprintf("(%s/%s, %s)", runtime.GOOS, runtime.GOARCH, runtime.Version())),
		)
	},
}

func init() {
	// Registered in root.go
}
