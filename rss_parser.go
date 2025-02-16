package godrudge

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

const (
	TopHeadline          = "top headline"
	MainHeadline         = "main headline"
	FirstColumnHeadline  = "first column"
	SecondColumnHeadline = "second column"
	ThirdColumnHeadline  = "third column"
)

func (c *Client) parseRSS() error {
	feed := c.rssFeed
	c.Page.HeadlineColumns = make([][]Headline, 3)
	for _, item := range feed.Items {
		if item.PublishedParsed != nil {
			headline := Headline{Title: item.Title, Href: item.Link, Color: Blue}
			headlineType, err := c.getHeadlineType(item.Description)
			if err != nil {
				return nil
			}
			switch headlineType {
			case TopHeadline:
				c.Page.TopHeadlines = append(c.Page.TopHeadlines, headline)
			case MainHeadline:
				c.Page.MainHeadlines = append(c.Page.MainHeadlines, headline)
			case FirstColumnHeadline:
				c.Page.HeadlineColumns[0] = append(c.Page.HeadlineColumns[0], headline)
			case SecondColumnHeadline:
				c.Page.HeadlineColumns[1] = append(c.Page.HeadlineColumns[1], headline)
			case ThirdColumnHeadline:
				c.Page.HeadlineColumns[2] = append(c.Page.HeadlineColumns[2], headline)
			default:
				return fmt.Errorf("failure parsing RSS headlines")
			}
		}
	}
	return nil
}

func (c *Client) getHeadlineType(description string) (string, error) {
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
				return TopHeadline
			case "main":
				return MainHeadline
			default:
				return ""
			}
		} else if node.Type == html.TextNode && strings.Contains(strings.ToLower(node.Data), "column") {
			re := regexp.MustCompile(`(\w+)\scolumn`)
			headlineType := re.FindStringSubmatch(strings.ToLower(node.Data))[1]
			switch headlineType {
			case "first":
				return FirstColumnHeadline
			case "second":
				return SecondColumnHeadline
			case "third":
				return ThirdColumnHeadline
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
