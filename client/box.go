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

func (c *Client) GetAllBoxes() (*[]api.Box, error) {
	url := fmt.Sprintf("%s/api/v1/boxes", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var boxes []api.Box
	if err := json.NewDecoder(resp.Body).Decode(&boxes); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &boxes, nil
}

func (c *Client) GetBox(id string) (*api.Box, error) {
	url := fmt.Sprintf("%s/api/v1/boxes/%s", c.baseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var box api.Box
	if err := json.NewDecoder(resp.Body).Decode(&box); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &box, nil
}

// CreateBox sends a PUT request to replace an existing box.
func (c *Client) ReplaceBox(box api.Box) (*api.Box, error) {
	url := fmt.Sprintf("%s/api/v1/boxes/%s", c.baseURL, box.ID)

	payload, err := json.Marshal(box)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal box: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var replacementBox api.Box
	if err := json.NewDecoder(resp.Body).Decode(&replacementBox); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &replacementBox, nil
}
