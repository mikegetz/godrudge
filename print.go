package godrudge

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"
)

func printDrudgeHeader(c *Client, terminalWidth int, textOnly bool) {
	fmt.Println(horizontalRule(terminalWidth, 1))
	fmt.Print(alignText(c.Page.Title, terminalWidth, "center"))
	fmt.Print(strings.Repeat("\n", 2))
	for _, headline := range c.Page.TopHeadlines {
		coloredHeadline := colorString(headline.Color, alignText(headline.Title, terminalWidth, "center"))
		if textOnly {
			fmt.Print(coloredHeadline)
		} else {
			fmt.Print(ansiLink(headline.Href, coloredHeadline))
		}
	}
	fmt.Print(strings.Repeat("\n", 2))
	fmt.Println(horizontalRule(terminalWidth, 1))
}

func printDrudgeBody(c *Client, terminalWidth int, textOnly bool) {
	numColumns := len(c.Page.HeadlineColumns)

	colWidth := terminalWidth / numColumns

	truncateWidth := colWidth - 3

	// Determine the maximum column size
	maxColumnSize := determineMaximumColumnSize(c.Page.HeadlineColumns)

	for row := 0; row < maxColumnSize; row++ {
		var line strings.Builder
		for column := 0; column < numColumns; column++ {
			var headline string
			if row < len(c.Page.HeadlineColumns[column]) {
				headline = truncateLine(c.Page.HeadlineColumns[column][row].Title, truncateWidth)

				alignment := "left"
				if column == 1 {
					alignment = "center"
				} else if column == 2 {
					alignment = "right"
				}

				headline = alignText(headline, colWidth, alignment)
				coloredHeadline := colorString(c.Page.HeadlineColumns[column][row].Color, headline)
				if textOnly {
					headline = coloredHeadline
				} else {
					headline = ansiLink(c.Page.HeadlineColumns[column][row].Href, coloredHeadline)
				}

			} else {
				headline = rowGap(terminalWidth, 3)
			}

			line.WriteString(headline)
		}
		fmt.Fprintln(os.Stdout, line.String())
	}
}

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
