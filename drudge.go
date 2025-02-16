package godrudge

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

type Client struct {
	domURL     string
	rssFeedURL string
	HTTPClient *http.Client
	Page       Page
	dom        *goquery.Document
	rssFeed    *gofeed.Feed
}

type Page struct {
	Title           string
	TopHeadlines    []Headline
	MainHeadlines   []Headline
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
		domURL:     "https://www.drudgereport.com",
		rssFeedURL: "http://feeds.feedburner.com/DrudgeReportFeed",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	for _, clientOverride := range c {
		defaultClient.HTTPClient = clientOverride
	}

	return defaultClient
}

func (c *Client) fetchDOM() error {
	resp, err := c.HTTPClient.Get(c.domURL)
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

func (c *Client) fetchRSS() error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(c.rssFeedURL)
	if err != nil {
		return nil
	}
	c.rssFeed = feed
	return nil
}

// Use HTML DOM parser
func (c *Client) ParseDOM() error {
	if c.dom == nil {
		err := c.fetchDOM()
		if err != nil {
			return err
		}
	}
	err := c.parseDOMTopHeadlines()
	if err != nil {
		return err
	}
	err = c.parseDOMMainHeadlines()
	if err != nil {
		return err
	}
	err = c.parseDOMHeadlines()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ParseRSS() error {
	if c.rssFeed == nil {
		err := c.fetchRSS()
		if err != nil {
			return err
		}
	}
	//TODO: parseRSS use c.rssFeed
	return nil
}

// Print Drudge
// Prints drudge page to stdout
//
// textOnly - prints to stdout without ansi links
func (c *Client) PrintDrudge(textOnly bool) {
	terminalWidth, _ := getTerminalWidth()
	printDrudgeTopHeadlines(c, terminalWidth, textOnly)
	printDrudgeMainHeadlines(c, terminalWidth, textOnly)
	printDrudgeBody(c, terminalWidth, textOnly)
}
