package godrudge

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
)

type Client struct {
	domURL     string
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
		domURL:     "https://www.drudgereport.com",
		rssFeedURL: "http://feeds.feedburner.com/DrudgeReportFeed",
	}

	return c
}

func (c *Client) fetchRSS() error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(c.rssFeedURL)
	if err != nil {
		return nil
	}
	c.rssFeed = feed
	return nil
}

func (c *Client) ParseRSS() error {
	if c.rssFeed == nil {
		err := c.fetchRSS()
		if err != nil {
			return err
		}
	}
	err := c.parseRSS()
	if err != nil {
		return err
	}
	return nil
}
