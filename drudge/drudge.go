package drudge

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/MGuitar24/go-drudge/color"
	"github.com/MGuitar24/go-drudge/printer"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Page struct {
	Title        string
	TopHeadlines []Headline
	Headlines    []Headline
}

type Headline struct {
	Title string
	Color color.Color
}

func (c *Client) FetchTopHeadlines() error {
	resp, err := c.HTTPClient.Get(c.BaseURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	title := doc.Find("title").Text()
	c.Page.Title = title

	// Find the <font> block that contains the main headline
	mainHeadlineTag := doc.Find(`font[size="+7"]`).First()

	// Find the first <a> tag within that block
	headlineSelector := mainHeadlineTag.Find("a").First()
	headlineNode := headlineSelector.Get(0)
	headlineColorText := headlineSelector.Find("font").First().AttrOr("color", "NOCOLOR")
	headlineText := extractTextWithNewlines(headlineNode)

	headlines := strings.Split(headlineText, "\n")

	for _, headline := range headlines {
		var headlineColor color.Color
		if strings.ToUpper(headlineColorText) == "RED" {
			headlineColor = color.Red
		} else {
			headlineColor = color.Blue
		}

		c.Page.TopHeadlines = append(c.Page.TopHeadlines, Headline{Title: headline, Color: headlineColor})
	}

	return err
}

func (c *Client) PrintHeadlines() {
	for _, headline := range c.Page.TopHeadlines {
		colorString := color.ColorString(headline.Color, headline.Title)
		fmt.Println(printer.CenterText(colorString))
	}
}

func extractTextWithNewlines(n *html.Node) string {
	var buf bytes.Buffer

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			buf.WriteString(node.Data)
		} else if node.Data == "br" {
			buf.WriteString("\n")
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(n)

	return strings.TrimSpace(buf.String())
}
