package printer

import (
	"os"
	"strconv"
	"strings"
)

func CenterText(text string) string {
	// Get terminal width
	width, _ := getTerminalWidth()

	// Calculate padding
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0 // Avoid negative padding if text is wider than the terminal
	}

	// Return centered text
	return strings.Repeat(" ", padding) + text
}

// Get terminal width (cross-platform)
func getTerminalWidth() (int, error) {
	if w, ok := os.LookupEnv("COLUMNS"); ok {
		return strconv.Atoi(w)
	}
	return 80, nil // Default width if not available
}
