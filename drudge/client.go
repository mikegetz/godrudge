package drudge

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Page       Page
	dom        *goquery.Document
}

func NewClient() *Client {
	return &Client{
		BaseURL: "https://www.drudgereport.com",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
