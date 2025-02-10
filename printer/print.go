package printer

import (
	"os"
	"strings"

	"golang.org/x/term"
)

func CenterText(text string, columns int) string {
	if columns < 1 {
		columns = 1
	}
	// Get terminal width
	width, _ := getTerminalWidth()

	// Calculate padding
	// re := regexp.MustCompile(`\033\[[0-9;]*m`)
	// colorlessText := re.ReplaceAllString(text, "")
	columnWidth := width / columns
	padding := (columnWidth - len(text)) / 2
	if padding < 0 {
		padding = 0 // Avoid negative padding if text is wider than the terminal
	}

	// Return centered text
	centered := strings.Repeat(" ", padding) + text
	if columnWidth < len(centered) {

		truncatedText := centered[:columnWidth-3]
		trimmedTruncatedText := strings.TrimSpace(truncatedText)

		untrimmedLength := len(truncatedText)
		trimmedLength := len(trimmedTruncatedText)

		ellipsesGap := 0
		if trimmedLength < untrimmedLength {
			ellipsesGap = untrimmedLength - trimmedLength
		}

		centered = strings.TrimSpace(centered[:columnWidth-3]) + strings.Repeat(".", ellipsesGap) + "..."
	} else if len(centered) < columnWidth {
		centered = centered + strings.Repeat(" ", columnWidth-len(centered))
	}
	return centered
}

func HorizontalRule(columns int) string {
	if columns < 1 {
		columns = 1
	}
	width, _ := getTerminalWidth()

	// Return horizontal rule
	return strings.Repeat("-", (width / columns))
}

func RowGap(columns int) string {
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
