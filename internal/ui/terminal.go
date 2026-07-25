package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

var (
	// Lip Gloss Color Definitions (Semantic)
	GreenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	YellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	RedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	BlueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	GrayStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	BoldStyle   = lipgloss.NewStyle().Bold(true)

	// Box Banner Style
	BannerStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 5).
			Bold(true)

	// Charm Logger
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false,
		ReportCaller:    false,
		TimeFormat:      time.Kitchen,
	})
)

// PrintBanner prints a clean rounded box banner using Lip Gloss
func PrintBanner(version string) {
	content := fmt.Sprintf("DBStudio %s", version)
	fmt.Println(BannerStyle.Render(content))
	fmt.Println()
}

// Standard Status Messages formatted to match CMD_DESIGN.md exactly
func PrintScanning() {
	fmt.Println(BlueStyle.Render("🔍 Scanning project..."))
}

func PrintSuccess(msg string) {
	fmt.Println(GreenStyle.Render("✔ ") + msg)
}

func PrintWarning(msg string) {
	fmt.Println(YellowStyle.Render("⚠ ") + msg)
}

func PrintError(msg string) {
	fmt.Println(RedStyle.Render("✖ ") + msg)
}

func PrintInfo(msg string) {
	fmt.Println(BlueStyle.Render("ℹ ") + msg)
}

func PrintStarting() {
	fmt.Println(BlueStyle.Render("\n🚀 Starting Studio..."))
}

func PrintListening(url string) {
	fmt.Println(GreenStyle.Render("✔ Listening on ") + BoldStyle.Render(url))
}

func PrintOpeningBrowser() {
	fmt.Println(BlueStyle.Render("🌐 Opening browser..."))
}

func PrintReady() {
	fmt.Println(GreenStyle.Render("✨ Ready.\n"))
}

// Verbose Log Helpers
func LogVerbose(msg string, keyvals ...interface{}) {
	Logger.Debug(msg, keyvals...)
}

// Backward-compatible Lip Gloss String Formatters
func Green(s string) string    { return GreenStyle.Render(s) }
func Yellow(s string) string   { return YellowStyle.Render(s) }
func Red(s string) string      { return RedStyle.Render(s) }
func Blue(s string) string     { return BlueStyle.Render(s) }
func Gray(s string) string     { return GrayStyle.Render(s) }
func BoldText(s string) string { return BoldStyle.Render(s) }
