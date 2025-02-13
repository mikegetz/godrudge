package godrudge

import (
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

// Align text based on column style
func alignText(text string, width int, alignment string) string {
	textLength := utf8.RuneCountInString(text)
	if textLength >= width {
		return text // If text is too long, return as is
	}

	totalPadding := width - textLength

	switch alignment {
	case "left": // Padding at the end (right side)
		return text + strings.Repeat(" ", totalPadding)
	case "center": // Padding split evenly
		leftPadding := totalPadding / 2
		rightPadding := totalPadding - leftPadding
		return strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding)
	case "right": // Padding at the front (left side)
		return strings.Repeat(" ", totalPadding) + text
	default:
		return text // Fallback (should not happen)
	}
}

func truncateLine(text string, maxLength int) string {
	if utf8.RuneCountInString(text) > maxLength {
		return text[:maxLength] + "..."
	}
	return text
}

func horizontalRule(terminalWidth int, columns int) string {
	if columns < 1 {
		columns = 1
	}

	// Return horizontal rule
	return strings.Repeat("-", (terminalWidth / columns))
}

// Create a row gap based on the number of columns
func rowGap(terminalWidth int, columns int) string {
	if columns < 1 {
		columns = 1
	}

	// Return row gap
	return strings.Repeat(" ", (terminalWidth / columns))
}

func getTerminalWidth() (int, error) {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 120, err
	}
	return width, nil
}
