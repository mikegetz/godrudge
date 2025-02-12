package godrudge

import (
	"os"
	"strings"

	"golang.org/x/term"
)

func centerText(text string, columns int) string {
	if columns < 1 {
		columns = 1
	}

	width, _ := getTerminalWidth()

	// Calculate padding
	columnWidth := width / columns
	padding := (columnWidth - len(text)) / 2
	if padding < 0 {
		padding = 0 // Avoid negative padding if text is wider than the terminal
	}

	// Return centered text
	centered := strings.Repeat(" ", padding) + text
	if columnWidth < len(centered) {
		centered = centered[:columnWidth-3] + "..."
	} else if len(centered) < columnWidth {
		centered = centered + strings.Repeat(" ", columnWidth-len(centered))
	}
	return centered
}

func horizontalRule(columns int) string {
	if columns < 1 {
		columns = 1
	}
	width, _ := getTerminalWidth()

	// Return horizontal rule
	return strings.Repeat("-", (width / columns))
}

func rowGap(columns int) string {
	if columns < 1 {
		columns = 1
	}
	width, _ := getTerminalWidth()

	// Return row gap
	return strings.Repeat(" ", (width / columns))
}

func getTerminalWidth() (int, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 120, err
	}
	return width, nil
}
