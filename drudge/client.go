package drudge

import (
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Page       Page
}

func NewClient() *Client {
	return &Client{
		BaseURL: "https://www.drudgereport.com",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
