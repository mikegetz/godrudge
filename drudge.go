package godrudge

import (
	"charm.land/lipgloss/v2"
	"github.com/mmcdole/gofeed"
)

type Client struct {
	rssFeedURL string
	Page       Page
	rssFeed    *gofeed.Feed
}

type Page struct {
	TopHeadlines    []Headline
	MainHeadlines   []Headline
	HeadlineColumns [][]Headline
}

type Headline struct {
	Title string
	URL   string
	Style lipgloss.Style
}

// provide a client override
func NewClient() *Client {
	c := &Client{
		rssFeedURL: "http://feeds.feedburner.com/DrudgeReportFeed",
	}

	return c
}

func (c *Client) ParseRSS() error {
	err := c.fetchRSS()
	if err != nil {
		return err
	}

	err = c.parseRSS()
	if err != nil {
		return err
	}

	return nil
}
