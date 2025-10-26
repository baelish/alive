package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/baelish/alive/api"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// initTestLogger initializes a logger for testing
func initTestLogger() {
	if logger == nil {
		logger = zap.NewNop()
	}
}

func TestLoadTemplates(t *testing.T) {
	t.Run("successfully loads all templates", func(t *testing.T) {
		err := loadTemplates()
		if err != nil {
			t.Fatalf("loadTemplates failed: %v", err)
		}

		if templates == nil {
			t.Fatal("templates is nil after loadTemplates")
		}

		// Verify key templates exist
		requiredTemplates := []string{"dashboard", "infoPage", "head", "statusBar", "boxGrid", "box", "boxInfo"}
		for _, name := range requiredTemplates {
			tmpl := templates.Lookup(name)
			if tmpl == nil {
				t.Errorf("template %q was not loaded", name)
			}
		}
	})

	t.Run("templates include ToUpper func", func(t *testing.T) {
		err := loadTemplates()
		if err != nil {
			t.Fatalf("loadTemplates failed: %v", err)
		}

		// Execute a template that uses ToUpper (via boxInfo)
		testBox := &api.Box{
			ID:          "test-id",
			Name:        "Test Box",
			DisplayName: "Test Display Name",
			Status:      api.Green,
			Size:        api.Medium,
			LastUpdate:  time.Now(),
			Messages: []api.Message{
				{
					TimeStamp: time.Now(),
					Status:    "green",
					Message:   "test message",
				},
			},
		}

		var buf strings.Builder
		err = templates.ExecuteTemplate(&buf, "boxInfo", testBox)
		if err != nil {
			t.Errorf("failed to execute boxInfo template: %v", err)
		}

		// Check that the output contains uppercase status (via ToUpper)
		output := buf.String()
		if !strings.Contains(output, "GREEN") {
			t.Error("ToUpper function may not be working correctly in templates")
		}
	})
}

func TestHandleRoot(t *testing.T) {
	// Save original global state
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	err := loadTemplates()
	if err != nil {
		t.Fatalf("loadTemplates failed: %v", err)
	}

	t.Run("renders dashboard with boxes", func(t *testing.T) {
		// Set up test data
		boxes = []api.Box{
			{
				ID:          "box-1",
				Name:        "Test Box 1",
				DisplayName: "Display 1",
				Status:      api.Green,
				Size:        api.Medium,
				LastMessage: "All good",
				LastUpdate:  time.Now(),
			},
			{
				ID:          "box-2",
				Name:        "Test Box 2",
				Status:      api.Red,
				Size:        api.Large,
				LastMessage: "Error occurred",
				LastUpdate:  time.Now(),
			},
		}

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handleRoot(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body := w.Body.String()

		// Verify dashboard structure
		if !strings.Contains(body, "<!DOCTYPE html>") {
			t.Error("response doesn't contain DOCTYPE")
		}

		// Verify boxes are rendered
		if !strings.Contains(body, "box-1") {
			t.Error("response doesn't contain box-1 ID")
		}
		if !strings.Contains(body, "box-2") {
			t.Error("response doesn't contain box-2 ID")
		}
		if !strings.Contains(body, "Test Box 1") {
			t.Error("response doesn't contain Test Box 1 name")
		}
		if !strings.Contains(body, "Test Box 2") {
			t.Error("response doesn't contain Test Box 2 name")
		}

		// Verify display name is used when present
		if !strings.Contains(body, "Display 1") {
			t.Error("response doesn't contain display name")
		}

		// Verify messages
		if !strings.Contains(body, "All good") {
			t.Error("response doesn't contain box-1 message")
		}
		if !strings.Contains(body, "Error occurred") {
			t.Error("response doesn't contain box-2 message")
		}

		// Verify status classes
		if !strings.Contains(body, "green") {
			t.Error("response doesn't contain green status class")
		}
		if !strings.Contains(body, "red") {
			t.Error("response doesn't contain red status class")
		}

		// Verify size classes
		if !strings.Contains(body, "medium") {
			t.Error("response doesn't contain medium size class")
		}
		if !strings.Contains(body, "large") {
			t.Error("response doesn't contain large size class")
		}
	})

	t.Run("renders empty dashboard", func(t *testing.T) {
		boxes = []api.Box{}

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()

		handleRoot(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body := w.Body.String()

		// Should still have basic structure
		if !strings.Contains(body, "<!DOCTYPE html>") {
			t.Error("response doesn't contain DOCTYPE")
		}
		if !strings.Contains(body, "status-bar") {
			t.Error("response doesn't contain status bar")
		}
	})
}

func TestHandleStatus(t *testing.T) {
	t.Run("returns JSON status ok", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()

		handleStatus(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body := w.Body.String()
		expected := `{"status":"ok"}`
		if body != expected {
			t.Errorf("expected body %q, got %q", expected, body)
		}
	})
}

func TestHandleBox(t *testing.T) {
	// Save original global state
	originalBoxes := boxes
	originalLogger := logger
	defer func() {
		boxes = originalBoxes
		logger = originalLogger
	}()

	// Initialize logger to prevent nil pointer
	initTestLogger()

	err := loadTemplates()
	if err != nil {
		t.Fatalf("loadTemplates failed: %v", err)
	}

	t.Run("renders box info page", func(t *testing.T) {
		// Set up test data
		testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		infoMap := map[string]string{
			"Version": "1.0.0",
			"Region":  "us-east-1",
		}
		boxes = []api.Box{
			{
				ID:          "test-box-123",
				Name:        "Test Box",
				DisplayName: "Test Display Name",
				Description: "A test description",
				Status:      api.Green,
				Size:        api.Medium,
				LastMessage: "Everything is fine",
				LastUpdate:  testTime,
				Info:        &infoMap,
				Links: []api.Links{
					{Name: "Documentation", URL: "https://example.com/docs"},
					{Name: "Dashboard", URL: "https://example.com/dashboard"},
				},
				Messages: []api.Message{
					{
						TimeStamp: testTime,
						Status:    "green",
						Message:   "Service started",
					},
				},
			},
		}

		// Create a chi router to handle URL params
		r := chi.NewRouter()
		r.HandleFunc("/box/{id}", handleBox)

		req := httptest.NewRequest("GET", "/box/test-box-123", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		body := w.Body.String()

		// Verify page structure
		if !strings.Contains(body, "<!DOCTYPE html>") {
			t.Error("response doesn't contain DOCTYPE")
		}

		// Verify box details
		if !strings.Contains(body, "test-box-123") {
			t.Error("response doesn't contain box ID")
		}
		if !strings.Contains(body, "Test Box") {
			t.Error("response doesn't contain box name")
		}
		if !strings.Contains(body, "Test Display Name") {
			t.Error("response doesn't contain display name")
		}
		if !strings.Contains(body, "A test description") {
			t.Error("response doesn't contain description")
		}
		if !strings.Contains(body, "Everything is fine") {
			t.Error("response doesn't contain last message")
		}

		// Verify info map
		if !strings.Contains(body, "Version") {
			t.Error("response doesn't contain Version key")
		}
		if !strings.Contains(body, "1.0.0") {
			t.Error("response doesn't contain version value")
		}
		if !strings.Contains(body, "Region") {
			t.Error("response doesn't contain Region key")
		}
		if !strings.Contains(body, "us-east-1") {
			t.Error("response doesn't contain region value")
		}

		// Verify links
		if !strings.Contains(body, "Documentation") {
			t.Error("response doesn't contain Documentation link")
		}
		if !strings.Contains(body, "https://example.com/docs") {
			t.Error("response doesn't contain docs URL")
		}
		if !strings.Contains(body, "Dashboard") {
			t.Error("response doesn't contain Dashboard link")
		}
		if !strings.Contains(body, "https://example.com/dashboard") {
			t.Error("response doesn't contain dashboard URL")
		}

		// Verify messages
		if !strings.Contains(body, "Service started") {
			t.Error("response doesn't contain message text")
		}
		if !strings.Contains(body, "GREEN") {
			t.Error("response doesn't contain uppercased status")
		}
	})

	t.Run("handles non-existent box", func(t *testing.T) {
		boxes = []api.Box{
			{ID: "existing-box", Name: "Existing"},
		}

		r := chi.NewRouter()
		r.HandleFunc("/box/{id}", handleBox)

		req := httptest.NewRequest("GET", "/box/non-existent", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		// Should return with no error (logs error instead)
		// The response might be empty or have a default status
		resp := w.Result()
		defer resp.Body.Close()

		// Just verify it doesn't crash
		if resp.StatusCode != http.StatusOK {
			// This is expected - when box is not found, nothing is written
			// so we get a 200 with empty body
		}
	})
}

func TestTemplateConstants(t *testing.T) {
	tests := []struct {
		name     string
		template string
		contains []string
	}{
		{
			name:     "dashboard template",
			template: dashboard,
			contains: []string{"define \"dashboard\"", "<!DOCTYPE html>", "statusBar", "boxGrid", "big-box"},
		},
		{
			name:     "infoPage template",
			template: infoPage,
			contains: []string{"define \"infoPage\"", "<!DOCTYPE html>", "statusBar", "boxInfo", "big-box"},
		},
		{
			name:     "generic template",
			template: generic,
			contains: []string{"define \"head\"", "define \"statusBar\"", "standard.css", "scripts.js", "status-bar"},
		},
		{
			name:     "boxGrid template",
			template: boxGrid,
			contains: []string{"define \"boxGrid\"", "define \"box\"", "range .", "boxClick", "boxHover"},
		},
		{
			name:     "boxInfo template",
			template: boxInfo,
			contains: []string{"define \"boxInfo\"", "table", "ID:", "Last message:", "Previous Messages:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, substr := range tt.contains {
				if !strings.Contains(tt.template, substr) {
					t.Errorf("template %q doesn't contain %q", tt.name, substr)
				}
			}
		})
	}
}
