package godrudge

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mikegetz/godrudge/color"
	"golang.org/x/net/html"
)

const (
	topHeadline          = "top headline"
	mainHeadline         = "main headline"
	firstColumnHeadline  = "first column"
	secondColumnHeadline = "second column"
	thirdColumnHeadline  = "third column"
)

func (c *Client) parseRSS() error {
	feed := c.rssFeed
	c.Page.HeadlineColumns = make([][]Headline, 3)
	for _, item := range feed.Items {
		if item.PublishedParsed != nil {
			var headlineColor color.Color
			if isRed, _ := isFeedHeadlineRed(item.Description); isRed {
				headlineColor = color.Red
			} else {
				headlineColor = color.Blue
			}
			headline := Headline{Title: item.Title, Href: item.Link, Color: headlineColor, ColorTitle: string(headlineColor) + item.Title + string(color.Reset)}
			headlineType, err := getFeedHeadlineType(item.Description)
			if err != nil {
				return err
			}
			switch headlineType {
			case topHeadline:
				c.Page.TopHeadlines = append(c.Page.TopHeadlines, headline)
			case mainHeadline:
				c.Page.MainHeadlines = append(c.Page.MainHeadlines, headline)
			case firstColumnHeadline:
				c.Page.HeadlineColumns[0] = append(c.Page.HeadlineColumns[0], headline)
			case secondColumnHeadline:
				c.Page.HeadlineColumns[1] = append(c.Page.HeadlineColumns[1], headline)
			case thirdColumnHeadline:
				c.Page.HeadlineColumns[2] = append(c.Page.HeadlineColumns[2], headline)
			default:
				return fmt.Errorf("failure parsing RSS headlines")
			}
		}
	}
	return nil
}

func isFeedHeadlineRed(description string) (bool, error) {
	doc, err := html.Parse(strings.NewReader(description))
	if err != nil {
		return false, err
	}
	var traverseFirstTree func(*html.Node) bool
	traverseFirstTree = func(node *html.Node) bool {
		if node.Type == html.ElementNode && node.Data == "font" {
			colorVal := extractNodeAttr(node, "color")
			if strings.ToLower(colorVal) == "red" {
				return true
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if traverseFirstTree(c) {
				return true
			}
		}
		return false
	}

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if traverseFirstTree(c) {
			return true, nil
		}
	}

	return false, nil
}

func getFeedHeadlineType(description string) (string, error) {
	doc, err := html.Parse(strings.NewReader(description))
	if err != nil {
		return "", err
	}

	var traverseFirstTree func(*html.Node) string
	traverseFirstTree = func(node *html.Node) string {
		if node.Type == html.TextNode && strings.Contains(strings.ToLower(node.Data), "headline") {
			re := regexp.MustCompile(`(\w+)\sheadline`)
			headlineType := re.FindStringSubmatch(strings.ToLower(node.Data))[1]
			switch headlineType {
			case "top":
				return topHeadline
			case "main":
				return mainHeadline
			default:
				return ""
			}
		} else if node.Type == html.TextNode && strings.Contains(strings.ToLower(node.Data), "column") {
			re := regexp.MustCompile(`(\w+)\scolumn`)
			headlineType := re.FindStringSubmatch(strings.ToLower(node.Data))[1]
			switch headlineType {
			case "first":
				return firstColumnHeadline
			case "second":
				return secondColumnHeadline
			case "third":
				return thirdColumnHeadline
			default:
				return ""
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if result := traverseFirstTree(c); result != "" {
				return result
			}
		}
		return ""
	}

	result := ""
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if result = traverseFirstTree(c); result != "" {
			break
		}
	}

	return result, nil

}
