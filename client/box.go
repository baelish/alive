package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/baelish/alive/api"
)

// CreateBox sends a POST request to create a new box.
func (c *Client) CreateBox(box api.Box) (*api.Box, error) {
	url := fmt.Sprintf("%s/api/v1/boxes", c.baseURL)

	payload, err := json.Marshal(box)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal box: %w", err)
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

	var createdBox api.Box
	if err := json.NewDecoder(resp.Body).Decode(&createdBox); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdBox, nil
}

func (c *Client) DeleteBox(id string) error {
	url := fmt.Sprintf("%s/api/v1/boxes/%s", c.baseURL, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
