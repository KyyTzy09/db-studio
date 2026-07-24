package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

const AppVersion = "v0.1.0-mvp"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of DBStudio",
	Long:  `Display DBStudio build version, Go runtime version, and operating system architecture.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("⚡ DBStudio %s (%s/%s, %s)\n", AppVersion, runtime.GOOS, runtime.GOARCH, runtime.Version())
	},
}

func init() {
	// Registered in root.go
}
