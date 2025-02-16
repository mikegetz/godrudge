package godrudge

import (
	"bytes"
	"net/http"
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
	Href  string
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

// Print Drudge
// Prints drudge page to stdout
//
// textOnly - prints to stdout without ansi links
func (c *Client) PrintDrudge(textOnly bool) {
	terminalWidth, _ := getTerminalWidth()
	printDrudgeHeader(c, terminalWidth, textOnly)
	printDrudgeBody(c, terminalWidth, textOnly)
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

func extractNodeAttr(n *html.Node, attrKey string) string {
	for _, attr := range n.Attr {
		if attr.Key == attrKey {
			return attr.Val
		}
	}
	return ""
}

// Extracts all headlines as h from startNode to stopNodeText
func extractHeadlines(startNode *html.Node, stopNodeText string) (h []Headline) {
	var buf bytes.Buffer
	const redTextPlaceholder = "<color:{red}>"
	const hrefPlaceholder = "<href>"
	const noColorTextPlaceholder = "<color:{none}>"

	var traverse func(*html.Node, *html.Node) bool
	traverse = func(node *html.Node, lastAnchor *html.Node) bool {
		if node.Type == html.CommentNode {
			//if we are on a comment node that matches a stopNodeText exit recursion
			if strings.Contains(strings.TrimSpace(strings.ReplaceAll(node.Data, " ", "")), stopNodeText) {
				return false
			}
		} else if node.Type == html.TextNode && strings.TrimSpace(node.Data) != "" {
			color := ""
			href := ""

			//set color
			if node.Parent.Type == html.ElementNode && node.Parent.Data == "font" {
				colorVal := extractNodeAttr(node.Parent, "color")
				if strings.ToLower(colorVal) == "red" {
					color = redTextPlaceholder
				}
			} else {
				color = noColorTextPlaceholder
			}

			//set url
			href = hrefPlaceholder + extractNodeAttr(lastAnchor, "href") + hrefPlaceholder

			buf.WriteString(href + color + strings.TrimSpace(node.Data) + color)
		} else if node.Type == html.ElementNode && node.Data == "a" {
			lastAnchor = node
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if !traverse(c, lastAnchor) {
				return false
			}
		}
		return true
	}
	traverse(startNode, nil)

	columnHeadlineString := strings.TrimSpace(buf.String())

	colorHeadlineStrings := sliceEveryOther(strings.Split(columnHeadlineString, "<color:"), 1)

	hrefHeadlineURLs := sliceEveryOther(strings.Split(columnHeadlineString, hrefPlaceholder), 1)

	coloredHeadlines := []Headline{}
	for index, colorHeadlineString := range colorHeadlineStrings {
		blueHeadlineString := strings.Split(colorHeadlineString, "{none}>")
		if len(blueHeadlineString) > 1 {
			coloredHeadlines = append(coloredHeadlines, Headline{Title: blueHeadlineString[1], Color: Blue, Href: hrefHeadlineURLs[index]})
		}

		redHeadlineString := strings.Split(colorHeadlineString, "{red}>")
		if len(redHeadlineString) > 1 {
			coloredHeadlines = append(coloredHeadlines, Headline{Title: redHeadlineString[1], Color: Red, Href: hrefHeadlineURLs[index]})
		}

	}

	return coloredHeadlines
}
