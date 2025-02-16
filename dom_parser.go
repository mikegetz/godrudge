package godrudge

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func (c *Client) parseDOMTopHeadlines() error {

	topHeadlineNode := findDOMTopHeadlineNode(c.dom)

	headlines := extractDOMHeadlines(topHeadlineNode, "TOPLEFTHEADLINESENDHERE")

	c.Page.TopHeadlines = append(c.Page.TopHeadlines, headlines...)

	return nil
}

func (c *Client) parseDOMMainHeadlines() error {

	title := c.dom.Find("title").Text()
	c.Page.Title = title

	mainHeadlineNode := findDOMMainHeadlineNode(c.dom)

	headlines := extractDOMHeadlines(mainHeadlineNode, "MAINHEADLINEENDHERE")

	c.Page.MainHeadlines = append(c.Page.MainHeadlines, headlines...)

	return nil
}

func (c *Client) parseDOMHeadlines() error {
	headlineColumns := [][]Headline{}
	subHeadlineStartNodeSelection := findDOMSubHeadlineNodeStartSelection(c.dom)
	columnStopLines := []string{"LINKSFIRSTCOLUMN", "LINKSSECONDCOLUMN", "LINKSANDSEARCHES3RDCOLUMN"}
	for count := 0; count < 3; count++ {
		headlinesNode := subHeadlineStartNodeSelection.Get(count)
		headlines := extractDOMHeadlines(headlinesNode, columnStopLines[count])
		headlineColumns = append(headlineColumns, headlines)
	}

	c.Page.HeadlineColumns = headlineColumns
	return nil
}

func findDOMTopHeadlineNode(dom *goquery.Document) (n *html.Node) {
	bodyNode := dom.Find("body").Get(0)
	var traverseFirstTree func(*html.Node) *html.Node
	traverseFirstTree = func(node *html.Node) *html.Node {
		if node.Type == html.CommentNode && strings.Contains(strings.ReplaceAll(node.Data, " ", ""), "TOPLEFTSTARTSHERE") {
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

func findDOMMainHeadlineNode(dom *goquery.Document) (n *html.Node) {
	bodyNode := dom.Find("body").Get(0)

	var traverseFirstTree func(*html.Node) *html.Node
	traverseFirstTree = func(node *html.Node) *html.Node {
		if node.Type == html.CommentNode && strings.Contains(strings.ReplaceAll(node.Data, " ", ""), "MAINHEADLINE") {
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

func findDOMSubHeadlineNodeStartSelection(dom *goquery.Document) (s *goquery.Selection) {
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
func extractDOMHeadlines(startNode *html.Node, stopNodeText string) (h []Headline) {
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
