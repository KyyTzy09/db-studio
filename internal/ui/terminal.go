package ui

import (
	"fmt"
	"strings"
)

// ANSI Escape Codes for Minimal & Professional Styling
const (
	Reset       = "\033[0m"
	Bold        = "\033[1m"
	Dim         = "\033[2m"
	GreenColor  = "\033[32m"
	YellowColor = "\033[33m"
	RedColor    = "\033[31m"
	BlueColor   = "\033[34m"
	CyanColor   = "\033[36m"
	GrayColor   = "\033[90m"
)

// Semantic Color Helpers
func Green(s string) string {
	return GreenColor + s + Reset
}

func Yellow(s string) string {
	return YellowColor + s + Reset
}

func Red(s string) string {
	return RedColor + s + Reset
}

func Blue(s string) string {
	return BlueColor + s + Reset
}

func Gray(s string) string {
	return GrayColor + s + Reset
}

func BoldText(s string) string {
	return Bold + s + Reset
}

// PrintBanner prints a clean rounded box banner as specified in CMD_DESIGN.md
func PrintBanner(version string) {
	title := fmt.Sprintf("           DBStudio %s            ", version)
	borderWidth := len(title)
	border := strings.Repeat("─", borderWidth)

	fmt.Println(Gray("╭" + border + "╮"))
	fmt.Println(Gray("│") + BoldText(title) + Gray("│"))
	fmt.Println(Gray("╰" + border + "╯"))
	fmt.Println()
}

// Standard Status Messages
func PrintScanning() {
	fmt.Println(Blue("🔍 Scanning project..."))
}

func PrintSuccess(msg string) {
	fmt.Println(Green("✔ ") + msg)
}

func PrintWarning(msg string) {
	fmt.Println(Yellow("⚠ ") + msg)
}

func PrintError(msg string) {
	fmt.Println(Red("✖ ") + msg)
}

func PrintInfo(msg string) {
	fmt.Println(Blue("ℹ ") + msg)
}

func PrintStarting() {
	fmt.Println(Blue("\n🚀 Starting Studio..."))
}

func PrintListening(url string) {
	fmt.Println(Green("✔ Listening on ") + BoldText(url))
}

func PrintOpeningBrowser() {
	fmt.Println(Blue("🌐 Opening browser..."))
}

func PrintReady() {
	fmt.Println(Green("✨ Ready.\n"))
}
