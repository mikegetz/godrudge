package drudge

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MGuitar24/go-drudge/color"
	"github.com/MGuitar24/go-drudge/printer"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Page       Page
	dom        *goquery.Document
}

type Page struct {
	Title           string
	TopHeadlines    []Headline
	HeadlineColumns [][]Headline
}

type Headline struct {
	Title string
	Color color.Color
}

func NewClient() *Client {
	return &Client{
		BaseURL: "https://www.drudgereport.com",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) fetch() error {
	resp, err := c.HTTPClient.Get(c.BaseURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dom, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	c.dom = dom
	return nil
}

func (c *Client) Parse() error {
	err := c.parseTopHeadlines()
	if err != nil {
		return err
	}
	err = c.parseHeadlines()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) parseHeadlines() error {
	if c.dom == nil {
		err := c.fetch()
		if err != nil {
			return err
		}
	}
	headlineColumns := [][]Headline{}
	ttNodes := c.dom.Find("body").Find("center").Find("table").Find("tt")
	for count := 0; count < 3; count++ {
		headlinesNode := ttNodes.Get(count)
		headlineText := extractTextWithNewlines(headlinesNode)

		headlines := []Headline{}
		headlineTextGroups := strings.Split(headlineText, "<hr>")
		if count == 2 {
			//drop last 8 headline text groups in 3rd column
			headlineTextGroups = headlineTextGroups[:len(headlineTextGroups)-8]
		} else {
			//drop last headline text group
			headlineTextGroups = headlineTextGroups[:len(headlineTextGroups)-1]
		}
		headlineText = strings.Join(headlineTextGroups, "<br>")
		for _, headline := range strings.Split(headlineText, "<br>") {
			if len(headline) > 0 {
				headlines = append(headlines, Headline{Title: headline, Color: color.Blue})
			}
		}

		headlineColumns = append(headlineColumns, headlines)
	}

	c.Page.HeadlineColumns = headlineColumns
	return nil
}

func (c *Client) parseTopHeadlines() error {
	if c.dom == nil {
		err := c.fetch()
		if err != nil {
			return err
		}
	}

	title := c.dom.Find("title").Text()
	c.Page.Title = title

	// Find the <font> block that contains the main headline
	mainHeadlineNode := findMainHeadlineNode(c.dom)

	headlineText := extractTextWithNewlines(mainHeadlineNode)

	headlines := strings.Split(headlineText, "<br>")

	for _, headline := range headlines {
		c.Page.TopHeadlines = append(c.Page.TopHeadlines, Headline{Title: headline, Color: color.Blue})
	}

	return nil
}

func (c *Client) PrintHeadlines() {
	fmt.Println(printer.HorizontalRule(1))
	fmt.Print(printer.CenterText(c.Page.Title, 1))
	fmt.Print(strings.Repeat("\n", 2))
	for _, headline := range c.Page.TopHeadlines {
		fmt.Print(color.ColorString(headline.Color, printer.CenterText(headline.Title, 1)))
	}
	fmt.Print(strings.Repeat("\n", 2))
	fmt.Println(printer.HorizontalRule(1))

	// Find max column count (length of longest inner slice)
	maxCols := 0
	for _, column := range c.Page.HeadlineColumns {
		if len(column) > maxCols {
			maxCols = len(column)
		}
	}

	// Iterate column by column
	for column := 0; column < maxCols; column++ {
		for row := 0; row < len(c.Page.HeadlineColumns); row++ {
			if column < len(c.Page.HeadlineColumns[row]) { // Avoid out-of-bounds access
				headline := c.Page.HeadlineColumns[row][column]
				fmt.Print(color.ColorString(headline.Color, printer.CenterText(headline.Title, 3)))
			} else {
				fmt.Print(printer.RowGap(3))
			}
		}
	}
}

func findMainHeadlineNode(dom *goquery.Document) (n *html.Node) {
	bodyNode := dom.Find("body").Get(0)

	var traverseFirstTree func(*html.Node) *html.Node
	traverseFirstTree = func(node *html.Node) *html.Node {
		if node.Type == html.CommentNode && strings.Contains(node.Data, "MAIN HEADLINE") {
			return node
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if result := traverseFirstTree(c); result != nil {
				return result
			}
		}
		return nil
	}

	var result *html.Node
	for c := bodyNode.FirstChild; c != nil; c = c.NextSibling {
		if result = traverseFirstTree(c); result != nil {
			break
		}
	}
	return result
}

func extractTextWithNewlines(n *html.Node) string {
	var buf bytes.Buffer

	var traverse func(*html.Node)
	traverse = func(node *html.Node) {
		if node.Type == html.TextNode {
			buf.WriteString(strings.TrimSpace(node.Data))
		} else if node.Data == "br" {
			buf.WriteString("<br>")
		} else if node.Data == "hr" {
			buf.WriteString("<hr>")
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(n)

	return strings.TrimSpace(buf.String())
}
