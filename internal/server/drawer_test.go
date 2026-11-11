package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/baelish/alive/api"
	"github.com/go-chi/chi/v5"
)

func TestBoxDrawerIsLoadedOnDemand(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()
	resetBoxStore()
	initTestLogger()
	if err := loadTemplates(); err != nil {
		t.Fatal(err)
	}

	boxStore.Add(api.Box{
		ID:          "drawer-box",
		Name:        "<b>Big</b>-horned sheep",
		DisplayName: "BHS",
		Size:        api.Micro,
		Stat:        "3/4",
		LastMessage: "hello",
		Links:       []api.Links{{Name: "bad", URL: "javascript:alert(1)"}},
		Messages:    []api.Message{{Message: "first", Status: "green"}},
	})

	router := chi.NewRouter()
	router.Get("/box/{id}/drawer", handleBoxDrawer)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/box/drawer-box/drawer", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", response.Code)
	}
	body := response.Body.String()
	for _, want := range []string{`class="drawer-title">BHS<`, "&lt;b&gt;Big&lt;/b&gt;-horned sheep", `drawer-stat`, ">3/4<", "first", "#ZgotmplZ"} {
		if !strings.Contains(body, want) {
			t.Errorf("drawer HTML does not contain %q", want)
		}
	}
	if strings.Contains(body, "<b>Big</b>") || strings.Contains(body, "javascript:") {
		t.Errorf("drawer HTML contains unescaped box data")
	}
	if got := response.Header().Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}

	response = httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/box/missing/drawer", nil))
	if response.Code != http.StatusNotFound {
		t.Errorf("unknown box: got status %d, want 404", response.Code)
	}

	// The dashboard page no longer embeds a drawer per box.
	response = httptest.NewRecorder()
	handleRoot(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if strings.Contains(response.Body.String(), "box-detail-") || strings.Contains(response.Body.String(), "drawer-tab") {
		t.Errorf("dashboard page still embeds box drawers")
	}
}
