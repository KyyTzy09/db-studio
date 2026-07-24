package main

import (
	"embed"
	"fmt"
	"os"

	"db-studio-go/cmd"
)

//go:embed all:web/build
var webFS embed.FS

func main() {
	rootCmd := cmd.NewRootCmd(webFS)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
