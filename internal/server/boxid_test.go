package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/baelish/alive/api"
	"github.com/go-chi/chi/v5"
)

func TestBoxIDValidation(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()
	resetBoxStore()

	router := chi.NewRouter()
	router.Post("/api/v1/boxes", apiCreateBox)
	router.Put("/api/v1/boxes/{id}", apiReplaceBox)

	do := func(method, path, body string) int {
		t.Helper()
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(method, path, strings.NewReader(body)))
		return response.Code
	}

	tests := []struct {
		name, method, path, body string
		want                     int
	}{
		{"create with hyphen and underscore", http.MethodPost, "/api/v1/boxes", `{"id":"mc-empires_1","name":"x"}`, http.StatusCreated},
		{"create without id", http.MethodPost, "/api/v1/boxes", `{"name":"x"}`, http.StatusCreated},
		{"create with space", http.MethodPost, "/api/v1/boxes", `{"id":"has space","name":"x"}`, http.StatusBadRequest},
		{"create with html", http.MethodPost, "/api/v1/boxes", `{"id":"<b>x</b>","name":"x"}`, http.StatusBadRequest},
		{"create too long", http.MethodPost, "/api/v1/boxes", `{"id":"` + strings.Repeat("a", 65) + `","name":"x"}`, http.StatusBadRequest},
		{"replace valid", http.MethodPut, "/api/v1/boxes/factorio-ava", `{"id":"factorio-ava","name":"x"}`, http.StatusCreated},
		{"replace with slash", http.MethodPut, "/api/v1/boxes/x", `{"id":"a/b","name":"x"}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		if got := do(tt.method, tt.path, tt.body); got != tt.want {
			t.Errorf("%s: got status %d, want %d", tt.name, got, tt.want)
		}
	}

	for _, b := range boxStore.GetAll() {
		if !api.ValidID(b.ID) {
			t.Errorf("stored box has invalid ID %q", b.ID)
		}
	}

	if _, err := addBox(api.Box{ID: "bad id"}); err == nil {
		t.Errorf("addBox accepted an invalid ID")
	}
}

func TestGeneratedIDsAreValid(t *testing.T) {
	for i := 0; i < 1000; i++ {
		if id := randStringBytes(10); !api.ValidID(id) {
			t.Fatalf("generated ID %q is not valid", id)
		}
	}
}

// A box stored before IDs were restricted can still be reached with an
// encoded URL.
func TestLegacyIDReachableByEncodedURL(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()
	resetBoxStore()
	initTestLogger()
	if templates == nil {
		if err := loadTemplates(); err != nil {
			t.Fatal(err)
		}
	}

	legacy := "old box/1"
	boxStore.mu.Lock()
	boxStore.boxes = []api.Box{{ID: legacy, Name: "old"}}
	boxStore.mu.Unlock()

	router := chi.NewRouter()
	router.Get("/api/v1/boxes/{id}", apiGetBox)
	router.Get("/box/{id}", handleBox)
	router.Delete("/api/v1/boxes/{id}", apiDeleteBox)

	for _, path := range []string{"/api/v1/boxes/old%20box%2F1", "/box/old%20box%2F1"} {
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != http.StatusOK {
			t.Errorf("GET %s: got status %d, want 200", path, response.Code)
		}
	}

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodDelete, "/api/v1/boxes/old%20box%2F1", nil))
	if boxStore.Exists(legacy) {
		t.Errorf("DELETE with encoded ID did not delete the box (status %d)", response.Code)
	}
}
