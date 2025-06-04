package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/baelish/alive/api"
)

// CreateEvent sends a POST request to create a new event.
func (c *Client) CreateEvent(event api.Event) (*api.Event, error) {
	url := fmt.Sprintf("%s/api/events", c.baseURL)

	payload, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdEvent api.Event
	if err := json.NewDecoder(resp.Body).Decode(&createdEvent); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdEvent, nil
}
