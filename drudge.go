package godrudge

import (
	"bytes"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

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
	Color Color
}

// provide a client override
func NewClient(c ...*http.Client) *Client {
	defaultClient := &Client{
		BaseURL: "https://www.drudgereport.com",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	for _, clientOverride := range c {
		defaultClient.HTTPClient = clientOverride
	}

	return defaultClient
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
	subHeadlineStartNodeSelection := findSubHeadlineNodeStartSelection(c.dom)
	columnStopLines := []string{"LINKSFIRSTCOLUMN", "LINKSSECONDCOLUMN", "LINKSANDSEARCHES3RDCOLUMN"}
	for count := 0; count < 3; count++ {
		headlinesNode := subHeadlineStartNodeSelection.Get(count)
		headlines := extractHeadlines(headlinesNode, columnStopLines[count])
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

	mainHeadlineNode := findMainHeadlineNode(c.dom)

	headlines := extractHeadlines(mainHeadlineNode, "MAINHEADLINEENDHERE")

	c.Page.TopHeadlines = append(c.Page.TopHeadlines, headlines...)

	return nil
}

func (c *Client) PrintHeadlines() {
	fmt.Println(horizontalRule(1))
	fmt.Print(centerText(c.Page.Title, 1))
	fmt.Print(strings.Repeat("\n", 2))
	for _, headline := range c.Page.TopHeadlines {
		fmt.Print(colorString(headline.Color, centerText(headline.Title, 1)))
	}
	fmt.Print(strings.Repeat("\n", 2))
	fmt.Println(horizontalRule(1))

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
				fmt.Print(colorString(headline.Color, centerText(headline.Title, 3)))
			} else {
				fmt.Print(rowGap(3))
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
	return result.Parent
}

func findSubHeadlineNodeStartSelection(dom *goquery.Document) (s *goquery.Selection) {
	return dom.Find("body").Find("center").Find("table").Find("tt")
}

// Extracts all headlines as h from startNode to stopNodeText
func extractHeadlines(startNode *html.Node, stopNodeText string) (h []Headline) {
	var buf bytes.Buffer
	const redTextPlaceholder = "<color:red>"

	var traverse func(*html.Node) bool
	traverse = func(node *html.Node) bool {
		if node.Type == html.CommentNode {
			if strings.Contains(strings.TrimSpace(strings.ReplaceAll(node.Data, " ", "")), stopNodeText) {
				return false
			}
		} else if node.Type == html.TextNode {
			color := ""
			if node.Parent.Type == html.ElementNode && node.Parent.Data == "font" {
				for _, attr := range node.Parent.Attr {
					if attr.Key == "color" && strings.ToLower(attr.Val) == "red" {
						color = redTextPlaceholder
					}
				}
			}
			buf.WriteString(color + strings.TrimSpace(node.Data) + color)
		} else if node.Data == "br" {
			buf.WriteString("<br>")
		} else if node.Data == "hr" {
			buf.WriteString("<br>")
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if !traverse(c) {
				return false
			}
		}
		return true
	}
	traverse(startNode)

	columnHeadlineString := strings.TrimSpace(buf.String())

	headlineStrings := strings.Split(columnHeadlineString, "<br>")

	//remove indexes with empty strings
	cleanedHeadlineStrings := slices.DeleteFunc(headlineStrings, func(s string) bool {
		return s == ""
	})
	//split cleaned headlines to get just red headlines based on redTextPlaceholder
	redHeadlineStrings := strings.Split(strings.Join(cleanedHeadlineStrings, ""), redTextPlaceholder)

	if len(redHeadlineStrings) == 1 {
		redHeadlineStrings = make([]string, 0)
	} else if len(redHeadlineStrings) > 1 {
		//red headlines only exist in every odd index ie. 1, 3, 5 etc. trimming the array to remove all even indexes 0, 2, 4 etc.
		newIndexCount := 0
		for count := 1; count < len(redHeadlineStrings); count += 2 {
			redHeadlineStrings[newIndexCount] = redHeadlineStrings[count]
			newIndexCount++
		}
		redHeadlineStrings = redHeadlineStrings[:newIndexCount]
	}

	coloredHeadlines := []Headline{}

	//add all cleanedHeadlines to coloredHeadlines
	for _, blueHeadlineString := range cleanedHeadlineStrings {
		coloredHeadlines = append(coloredHeadlines, Headline{Title: blueHeadlineString, Color: Blue})
	}

	//replace blue headlines with red headlines to maintain headline order
	for _, redHeadline := range redHeadlineStrings {
		indexOfRedHeadline := slices.IndexFunc(coloredHeadlines, func(h Headline) bool {
			return strings.Contains(h.Title, redHeadline)
		})
		coloredHeadlines[indexOfRedHeadline] = Headline{Title: redHeadline, Color: Red}
	}

	return coloredHeadlines
}
