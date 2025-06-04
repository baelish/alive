package client

import (
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient initializes and returns a new API client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}
