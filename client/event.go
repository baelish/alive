package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/baelish/alive/api"
)

// CreateEvent sends a POST request to create a new event.
func (c *Client) CreateEvent(event api.Event) error {
	url := fmt.Sprintf("%s/api/v1/boxes/%s/events", c.baseURL, event.ID)

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
