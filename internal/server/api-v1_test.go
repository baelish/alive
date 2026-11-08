package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestGetBoxIncludesConfiguredDurations(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()
	resetBoxStore()

	router := chi.NewRouter()
	router.Post("/api/v1/boxes", apiCreateBox)
	router.Get("/api/v1/boxes/{id}", apiGetBox)
	router.Post("/api/v1/boxes/{id}/events", apiCreateEvent)

	checkDurations := func(expireAfter, maxTBU string) {
		t.Helper()
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/boxes/test-box", nil))
		if response.Code != http.StatusOK {
			t.Fatalf("get status = %d, body = %s", response.Code, response.Body.String())
		}

		var body map[string]json.RawMessage
		if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		for field, want := range map[string]string{"expireAfter": expireAfter, "maxTBU": maxTBU} {
			if want == "" {
				if _, present := body[field]; present {
					t.Errorf("%s should be absent after clearing, got %s", field, body[field])
				}
				continue
			}
			if string(body[field]) != want {
				t.Errorf("%s = %s, want %s", field, body[field], want)
			}
		}
	}

	created := httptest.NewRecorder()
	router.ServeHTTP(created, httptest.NewRequest(http.MethodPost, "/api/v1/boxes", strings.NewReader(`{"id":"test-box","name":"Test","size":"small","status":"grey","expireAfter":"2h","maxTBU":"30m"}`)))
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d, body = %s", created.Code, created.Body.String())
	}
	checkDurations(`"2h0m0s"`, `"30m0s"`)

	updated := httptest.NewRecorder()
	router.ServeHTTP(updated, httptest.NewRequest(http.MethodPost, "/api/v1/boxes/test-box/events", strings.NewReader(`{"status":"green","expireAfter":"3h","maxTBU":"45m"}`)))
	if updated.Code != http.StatusCreated {
		t.Fatalf("update status = %d, body = %s", updated.Code, updated.Body.String())
	}
	checkDurations(`"3h0m0s"`, `"45m0s"`)

	cleared := httptest.NewRecorder()
	router.ServeHTTP(cleared, httptest.NewRequest(http.MethodPost, "/api/v1/boxes/test-box/events", strings.NewReader(`{"status":"green","expireAfter":0,"maxTBU":0}`)))
	if cleared.Code != http.StatusCreated {
		t.Fatalf("clear status = %d, body = %s", cleared.Code, cleared.Body.String())
	}
	checkDurations("", "")
}
