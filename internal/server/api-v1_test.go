// api-v1_test.go
package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baelish/alive/api"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Initialize a no-op logger for tests
	logger = zap.NewNop()

	// Initialize boxes slice if needed
	if boxes == nil {
		boxes = []api.Box{}
	}

	// Initialize the events broker for tests
	events = &Broker{
		clients:        make(map[chan string]bool),
		newClients:     make(chan (chan string)),
		defunctClients: make(chan (chan string)),
		messages:       make(chan string, 100),
	}

	// Start a goroutine to drain the messages channel so it doesn't block
	go func() {
		for range events.messages {
			// Discard events in tests
		}
	}()

	m.Run()
}

// TestApiStatus tests the /health endpoint
func TestApiStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	apiStatus(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json, got %s", contentType)
	}

	var response map[string]string
	err := json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("expected status 'ok', got '%s'", response["status"])
	}
}

// TestApiGetBoxes tests getting all boxes
func TestApiGetBoxes(t *testing.T) {
	// Save original boxes and restore after test
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	// Set up test data
	boxes = []api.Box{
		{ID: "1", Name: "Box 1", Status: api.Green, Size: api.Small},
		{ID: "2", Name: "Box 2", Status: api.Red, Size: api.Medium},
	}

	req := httptest.NewRequest("GET", "/api/v1/boxes", nil)
	rr := httptest.NewRecorder()

	apiGetBoxes(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type application/json, got %s", contentType)
	}

	var responseBoxes []api.Box
	err := json.NewDecoder(rr.Body).Decode(&responseBoxes)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(responseBoxes) != 2 {
		t.Errorf("expected 2 boxes, got %d", len(responseBoxes))
	}

	if responseBoxes[0].Name != "Box 1" {
		t.Errorf("expected first box name 'Box 1', got '%s'", responseBoxes[0].Name)
	}
}

// TestApiGetBox tests getting a specific box by ID
func TestApiGetBox(t *testing.T) {
	// Save original boxes and restore after test
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	tests := []struct {
		name           string
		boxID          string
		setupBoxes     []api.Box
		expectedStatus int
		expectError    bool
	}{
		{
			name:  "existing box",
			boxID: "test-123",
			setupBoxes: []api.Box{
				{ID: "test-123", Name: "Test Box", Status: api.Green, Size: api.Medium},
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "non-existing box",
			boxID:          "nonexistent",
			setupBoxes:     []api.Box{},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boxes = tt.setupBoxes

			r := chi.NewRouter()
			r.Get("/api/v1/boxes/{id}", apiGetBox)

			req := httptest.NewRequest("GET", "/api/v1/boxes/"+tt.boxID, nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if !tt.expectError {
				var box api.Box
				err := json.NewDecoder(rr.Body).Decode(&box)
				if err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if box.ID != tt.boxID {
					t.Errorf("expected box ID %s, got %s", tt.boxID, box.ID)
				}
			}
		})
	}
}

// TestStatusMarshaling tests the Status type JSON marshaling
func TestStatusMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		status   api.Status
		expected string
	}{
		{"green", api.Green, `"green"`},
		{"red", api.Red, `"red"`},
		{"amber", api.Amber, `"amber"`},
		{"grey", api.Grey, `"grey"`},
		{"noUpdate", api.NoUpdate, `"noUpdate"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.status)
			if err != nil {
				t.Fatalf("failed to marshal status: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(data))
			}

			// Test unmarshaling
			var s api.Status
			err = json.Unmarshal(data, &s)
			if err != nil {
				t.Fatalf("failed to unmarshal status: %v", err)
			}

			if s != tt.status {
				t.Errorf("expected status %v, got %v", tt.status, s)
			}
		})
	}
}

// TestBoxSizeMarshaling tests the BoxSize type JSON marshaling
func TestBoxSizeMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		size     api.BoxSize
		expected string
	}{
		{"dot", api.Dot, `"dot"`},
		{"micro", api.Micro, `"micro"`},
		{"small", api.Small, `"small"`},
		{"medium", api.Medium, `"medium"`},
		{"large", api.Large, `"large"`},
		{"xlarge", api.Xlarge, `"xlarge"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.size)
			if err != nil {
				t.Fatalf("failed to marshal box size: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(data))
			}

			// Test unmarshaling
			var bs api.BoxSize
			err = json.Unmarshal(data, &bs)
			if err != nil {
				t.Fatalf("failed to unmarshal box size: %v", err)
			}

			if bs != tt.size {
				t.Errorf("expected box size %v, got %v", tt.size, bs)
			}
		})
	}
}

// TestBoxWithAllFields tests creating a box with all fields populated
func TestBoxWithAllFields(t *testing.T) {
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	boxes = make([]api.Box, 0)

	info := map[string]string{"key": "value"}
	newBox := api.Box{
		Name:        "Full Box",
		Description: "A complete box",
		DisplayName: "Display Name",
		Status:      api.Green,
		Size:        api.Medium,
		Info:        &info,
		Parent:      "parent-id",
		Links: []api.Links{
			{Name: "GitHub", URL: "https://github.com"},
		},
		Messages: []api.Message{
			{Message: "Test", Status: "ok"},
		},
	}

	body, _ := json.Marshal(newBox)
	req := httptest.NewRequest("POST", "/api/v1/boxes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	apiCreateBox(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var createdBox api.Box
	err := json.NewDecoder(rr.Body).Decode(&createdBox)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdBox.Name != newBox.Name {
		t.Errorf("expected name %s, got %s", newBox.Name, createdBox.Name)
	}

	if createdBox.Description != newBox.Description {
		t.Errorf("expected description %s, got %s", newBox.Description, createdBox.Description)
	}

	if len(createdBox.Links) != 1 {
		t.Errorf("expected 1 link, got %d", len(createdBox.Links))
	}
}

// TestApiCreateBox tests creating a new box
func TestApiCreateBox(t *testing.T) {
	// Save original boxes and restore after test
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	// Initialize empty boxes slice
	boxes = make([]api.Box, 0)

	newBox := api.Box{
		Name:   "New Test Box",
		Status: api.Green,
		Size:   api.Medium,
	}

	body, _ := json.Marshal(newBox)
	req := httptest.NewRequest("POST", "/api/v1/boxes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	apiCreateBox(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	location := rr.Header().Get("Location")
	if location == "" {
		t.Error("expected Location header to be set")
	}

	var createdBox api.Box
	err := json.NewDecoder(rr.Body).Decode(&createdBox)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdBox.ID == "" {
		t.Error("expected created box to have an ID")
	}

	if createdBox.Name != newBox.Name {
		t.Errorf("expected box name %s, got %s", newBox.Name, createdBox.Name)
	}

	if createdBox.Status != newBox.Status {
		t.Errorf("expected box status %v, got %v", newBox.Status, createdBox.Status)
	}
}

// TestApiCreateBox_InvalidJSON tests creating a box with invalid JSON
func TestApiCreateBox_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/v1/boxes", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	apiCreateBox(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var errResp api.ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	if err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if errResp.Message == "" {
		t.Error("expected error message to be present")
	}
}

// TestApiReplaceBox tests replacing an existing box
func TestApiReplaceBox(t *testing.T) {
	// Save original boxes and restore after test
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	tests := []struct {
		name           string
		urlID          string
		boxToReplace   api.Box
		setupBoxes     []api.Box
		expectedStatus int
	}{
		{
			name:  "replace existing box",
			urlID: "existing-id",
			boxToReplace: api.Box{
				ID:     "existing-id",
				Name:   "Updated Name",
				Status: api.Amber,
				Size:   api.Large,
			},
			setupBoxes: []api.Box{
				{ID: "existing-id", Name: "Original Name", Status: api.Green, Size: api.Small},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "create new box when not found",
			urlID: "new-id",
			boxToReplace: api.Box{
				ID:     "new-id",
				Name:   "New Box",
				Status: api.Green,
				Size:   api.Medium,
			},
			setupBoxes:     []api.Box{},
			expectedStatus: http.StatusCreated,
		},
		{
			name:  "missing ID in body",
			urlID: "some-id",
			boxToReplace: api.Box{
				Name:   "No ID Box",
				Status: api.Green,
				Size:   api.Medium,
			},
			setupBoxes:     []api.Box{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boxes = tt.setupBoxes

			body, _ := json.Marshal(tt.boxToReplace)
			req := httptest.NewRequest("PUT", "/api/v1/boxes/"+tt.urlID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			r := chi.NewRouter()
			r.Put("/api/v1/boxes/{id}", apiReplaceBox)
			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

// TestApiDeleteBox tests deleting a box
func TestApiDeleteBox(t *testing.T) {
	// Save original boxes and restore after test
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	tests := []struct {
		name           string
		boxID          string
		setupBoxes     []api.Box
		expectedStatus int
	}{
		{
			name:  "delete existing box",
			boxID: "delete-me",
			setupBoxes: []api.Box{
				{ID: "delete-me", Name: "To Delete", Status: api.Red, Size: api.Small},
				{ID: "keep-me", Name: "Keep This", Status: api.Green, Size: api.Medium},
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "delete non-existing box",
			boxID:          "nonexistent",
			setupBoxes:     []api.Box{},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boxes = tt.setupBoxes

			r := chi.NewRouter()
			r.Delete("/api/v1/boxes/{id}", apiDeleteBox)

			req := httptest.NewRequest("DELETE", "/api/v1/boxes/"+tt.boxID, nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

// TestApiCreateEvent tests creating an event for a box
func TestApiCreateEvent(t *testing.T) {
	event := api.Event{
		Status:  api.Green,
		Message: "test message",
	}

	body, _ := json.Marshal(event)
	req := httptest.NewRequest("POST", "/api/v1/boxes/box-123/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Post("/api/v1/boxes/{id}/events", apiCreateEvent)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	var createdEvent api.Event
	err := json.NewDecoder(rr.Body).Decode(&createdEvent)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdEvent.ID != "box-123" {
		t.Errorf("expected event ID 'box-123', got '%s'", createdEvent.ID)
	}

	if createdEvent.Type != "updateBox" {
		t.Errorf("expected event type 'updateBox', got '%s'", createdEvent.Type)
	}

	if createdEvent.Status != api.Green {
		t.Errorf("expected status Green, got %v", createdEvent.Status)
	}
}

// TestDeprecatedRoute tests the deprecated route middleware
func TestDeprecatedRoute(t *testing.T) {
	handler := DeprecatedRoute("test deprecated message")(apiStatus)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	warning := rr.Header().Get("Warning")
	if warning == "" {
		t.Error("expected Warning header to be set")
	}

	expectedWarning := `299 alive "test deprecated message"`
	if warning != expectedWarning {
		t.Errorf("expected warning '%s', got '%s'", expectedWarning, warning)
	}
}

// TestHandleApiErrorResponse tests the error response handler
func TestHandleApiErrorResponse(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		err          error
		message      string
		includeError bool
	}{
		{
			name:         "with error included",
			status:       http.StatusBadRequest,
			err:          errors.New("test error"),
			message:      "test error message",
			includeError: true,
		},
		{
			name:         "without error included",
			status:       http.StatusInternalServerError,
			err:          errors.New("internal error occurred"),
			message:      "internal error",
			includeError: false,
		},
		{
			name:         "no error object",
			status:       http.StatusNotFound,
			err:          nil,
			message:      "not found",
			includeError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			handleApiErrorResponse(rr, tt.status, tt.err, tt.message, tt.includeError, true)

			if rr.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, rr.Code)
			}

			var errResp api.ErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&errResp)
			if err != nil {
				t.Fatalf("failed to decode error response: %v", err)
			}

			if errResp.Message != tt.message {
				t.Errorf("expected message '%s', got '%s'", tt.message, errResp.Message)
			}

			if tt.includeError && tt.err != nil {
				if errResp.Error == "" {
					t.Error("expected error field to be populated")
				}
			}
		})
	}
}
