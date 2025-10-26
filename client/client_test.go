package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baelish/alive/api"
)

func TestNewClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := NewClient(baseURL)

	if client == nil {
		t.Fatal("expected client to be non-nil")
	}

	if client.baseURL != baseURL {
		t.Errorf("expected baseURL %s, got %s", baseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("expected httpClient to be initialized")
	}

	if client.httpClient.Timeout == 0 {
		t.Error("expected httpClient timeout to be set")
	}
}

func TestCreateBox(t *testing.T) {
	tests := []struct {
		name           string
		box            api.Box
		responseStatus int
		responseBody   api.Box
		expectError    bool
	}{
		{
			name: "successful creation",
			box: api.Box{
				Name:   "Test Box",
				Status: api.Green,
				Size:   api.Medium,
			},
			responseStatus: http.StatusCreated,
			responseBody: api.Box{
				ID:     "created-123",
				Name:   "Test Box",
				Status: api.Green,
				Size:   api.Medium,
			},
			expectError: false,
		},
		{
			name: "server error",
			box: api.Box{
				Name: "Test Box",
			},
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST request, got %s", r.Method)
				}

				if r.URL.Path != "/api/v1/boxes" {
					t.Errorf("expected path /api/v1/boxes, got %s", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)

				if tt.responseStatus == http.StatusCreated {
					json.NewEncoder(w).Encode(tt.responseBody)
				}
			}))
			defer server.Close()

			client := NewClient(server.URL)
			createdBox, err := client.CreateBox(tt.box)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if createdBox.ID != tt.responseBody.ID {
				t.Errorf("expected ID %s, got %s", tt.responseBody.ID, createdBox.ID)
			}

			if createdBox.Name != tt.responseBody.Name {
				t.Errorf("expected Name %s, got %s", tt.responseBody.Name, createdBox.Name)
			}
		})
	}
}

func TestDeleteBox(t *testing.T) {
	tests := []struct {
		name           string
		boxID          string
		responseStatus int
		expectError    bool
	}{
		{
			name:           "successful deletion",
			boxID:          "test-123",
			responseStatus: http.StatusNoContent,
			expectError:    false,
		},
		{
			name:           "box not found",
			boxID:          "nonexistent",
			responseStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("expected DELETE request, got %s", r.Method)
				}

				expectedPath := "/api/v1/boxes/" + tt.boxID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.responseStatus)
			}))
			defer server.Close()

			client := NewClient(server.URL)
			err := client.DeleteBox(tt.boxID)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGetAllBoxes(t *testing.T) {
	tests := []struct {
		name           string
		responseStatus int
		responseBoxes  []api.Box
		expectError    bool
	}{
		{
			name:           "successful retrieval",
			responseStatus: http.StatusOK,
			responseBoxes: []api.Box{
				{ID: "1", Name: "Box 1", Status: api.Green, Size: api.Small},
				{ID: "2", Name: "Box 2", Status: api.Red, Size: api.Medium},
			},
			expectError: false,
		},
		{
			name:           "empty list",
			responseStatus: http.StatusOK,
			responseBoxes:  []api.Box{},
			expectError:    false,
		},
		{
			name:           "server error",
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				if r.URL.Path != "/api/v1/boxes" {
					t.Errorf("expected path /api/v1/boxes, got %s", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)

				if tt.responseStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.responseBoxes)
				}
			}))
			defer server.Close()

			client := NewClient(server.URL)
			boxes, err := client.GetAllBoxes()

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(*boxes) != len(tt.responseBoxes) {
				t.Errorf("expected %d boxes, got %d", len(tt.responseBoxes), len(*boxes))
			}
		})
	}
}

func TestGetBox(t *testing.T) {
	tests := []struct {
		name           string
		boxID          string
		responseStatus int
		responseBox    api.Box
		expectError    bool
	}{
		{
			name:           "successful retrieval",
			boxID:          "test-123",
			responseStatus: http.StatusOK,
			responseBox: api.Box{
				ID:     "test-123",
				Name:   "Test Box",
				Status: api.Green,
				Size:   api.Medium,
			},
			expectError: false,
		},
		{
			name:           "box not found",
			boxID:          "nonexistent",
			responseStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("expected GET request, got %s", r.Method)
				}

				expectedPath := "/api/v1/boxes/" + tt.boxID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)

				if tt.responseStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.responseBox)
				}
			}))
			defer server.Close()

			client := NewClient(server.URL)
			box, err := client.GetBox(tt.boxID)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if box.ID != tt.responseBox.ID {
				t.Errorf("expected ID %s, got %s", tt.responseBox.ID, box.ID)
			}

			if box.Name != tt.responseBox.Name {
				t.Errorf("expected Name %s, got %s", tt.responseBox.Name, box.Name)
			}
		})
	}
}

func TestReplaceBox(t *testing.T) {
	tests := []struct {
		name           string
		box            api.Box
		responseStatus int
		responseBox    api.Box
		expectError    bool
	}{
		{
			name: "successful replacement",
			box: api.Box{
				ID:     "test-123",
				Name:   "Updated Box",
				Status: api.Amber,
				Size:   api.Large,
			},
			responseStatus: http.StatusOK,
			responseBox: api.Box{
				ID:     "test-123",
				Name:   "Updated Box",
				Status: api.Amber,
				Size:   api.Large,
			},
			expectError: false,
		},
		{
			name: "server error",
			box: api.Box{
				ID:   "test-123",
				Name: "Updated Box",
			},
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PUT" {
					t.Errorf("expected PUT request, got %s", r.Method)
				}

				expectedPath := "/api/v1/boxes/" + tt.box.ID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.responseStatus)

				if tt.responseStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tt.responseBox)
				}
			}))
			defer server.Close()

			client := NewClient(server.URL)
			replacedBox, err := client.ReplaceBox(tt.box)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if replacedBox.ID != tt.responseBox.ID {
				t.Errorf("expected ID %s, got %s", tt.responseBox.ID, replacedBox.ID)
			}

			if replacedBox.Name != tt.responseBox.Name {
				t.Errorf("expected Name %s, got %s", tt.responseBox.Name, replacedBox.Name)
			}
		})
	}
}

func TestCreateEvent(t *testing.T) {
	tests := []struct {
		name           string
		event          api.Event
		responseStatus int
		expectError    bool
	}{
		{
			name: "successful event creation",
			event: api.Event{
				ID:      "box-123",
				Status:  api.Green,
				Message: "Test event",
				Type:    "updateBox",
			},
			responseStatus: http.StatusCreated,
			expectError:    false,
		},
		{
			name: "server error",
			event: api.Event{
				ID:      "box-123",
				Message: "Test event",
			},
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("expected POST request, got %s", r.Method)
				}

				expectedPath := "/api/v1/boxes/" + tt.event.ID + "/events"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.responseStatus)
			}))
			defer server.Close()

			client := NewClient(server.URL)
			err := client.CreateEvent(tt.event)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}
