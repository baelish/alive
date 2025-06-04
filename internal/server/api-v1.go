package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/baelish/alive/api"

	"github.com/go-chi/chi/v5"
)

func handleApiErrorResponse(w http.ResponseWriter, status int, e error, message string, includeError bool, skipServerLog bool) {
	if e != nil && !skipServerLog {
		logger.Error(e.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	resp := api.ErrorResponse{
		Message: message,
	}

	if e != nil && includeError {
		resp.Error = e.Error()
	}

	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logger.Error(err.Error())
	}
}

func apiGetBoxes(w http.ResponseWriter, _ *http.Request) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(boxes)
	if err != nil {
		handleApiErrorResponse(w, http.StatusInternalServerError, err, "could not get boxes", false, false)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func apiGetBox(w http.ResponseWriter, r *http.Request) {
	i, err := findBoxByID(chi.URLParam(r, "id"))

	if err != nil {
		handleApiErrorResponse(w, http.StatusNotFound, err, "id not found", false, false)

		return
	}

	// TODO Mutex boxes
	_ = json.NewEncoder(w).Encode(boxes[i])
}

func apiCreateEvent(w http.ResponseWriter, r *http.Request) {
	var event api.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		handleApiErrorResponse(w, http.StatusBadRequest, err, "could not decode data received", true, false)

		return
	}

	event.ID = chi.URLParam(r, "id")
	event.Type = "updateBox"
	logger.Debug("update event details", logStructDetails(event)...)
	update(event)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(event); err != nil {
		logger.Error("failed to encode response: " + err.Error())
	}
}

func apiCreateBox(w http.ResponseWriter, r *http.Request) {
	var newBox api.Box

	err := json.NewDecoder(r.Body).Decode(&newBox)
	if err != nil {
		handleApiErrorResponse(w, http.StatusBadRequest, err, "failed to decode data received", true, false)

		return
	}

	id, err := addBox(newBox)
	if err != nil {
		handleApiErrorResponse(w, http.StatusInternalServerError, err, "failed to create the box", false, false)

		return
	}

	newBox.ID = id

	w.Header().Set("Location", fmt.Sprintf("/api/boxes/%s", newBox.ID))
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newBox)
	if err != nil {
		logger.Error(err.Error())
	}
}

// Will replace the box if found, will create a new one if not found
func apiReplaceBox(w http.ResponseWriter, r *http.Request) {
	var newBox api.Box
	if err := json.NewDecoder(r.Body).Decode(&newBox); err != nil {
		handleApiErrorResponse(w, http.StatusBadRequest, err, "failed to decode data received", true, false)

		return
	}

	if newBox.ID == "" {
		// Create a custom error for missing ID
		missingIDErr := errors.New("missing ID when requesting a box replacement")
		handleApiErrorResponse(w, http.StatusBadRequest, missingIDErr, "cannot replace a box without an ID", true, false)
		return
	}

	// TODO Add mutex
	found, oldBox := deleteBox(newBox.ID, true)
	if _, err := addBox(newBox); err != nil {
		var msg string
		if found {
			if _, err = addBox(oldBox); err != nil {
				msg = "failed to replace the box and the old box was lost"
			} else {
				msg = "failed to replace the box and the old box was restored"
			}
		} else {
			msg = "failed to create the box"
		}
		handleApiErrorResponse(w, http.StatusInternalServerError, err, msg, false, false)

		return
	}

	// Send success
	if found {
		w.WriteHeader(http.StatusOK) // replaced existing

	} else {
		w.WriteHeader(http.StatusCreated) // created new
	}
	if err := json.NewEncoder(w).Encode(newBox); err != nil {
		logger.Error(err.Error())
	}
}

func apiStatus(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]string{"status": "ok"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error(err.Error())
	}
}

func apiDeleteBox(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if found, _ := deleteBox(id, true); found {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	err := json.NewEncoder(w).Encode(map[string]string{"error": "box not found"})
	if err != nil {
		logger.Error(err.Error())
	}

}

func DeprecatedRoute(msg string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Warning", `299 alive "`+msg+`"`)
			next.ServeHTTP(w, r)
		})
	}
}

func runAPI(_ context.Context) {
	if options.Debug {
		logger.Info("Starting up API")
	}
	router := chi.NewRouter()
	router.Get("/health", apiStatus)
	router.Get("/api/v1/boxes", apiGetBoxes)                 // Get all boxes
	router.Post("/api/v1/boxes", apiCreateBox)               // Create a new box
	router.Put("/api/v1/boxes/{id}", apiReplaceBox)          // Replace an existing box
	router.Delete("/api/v1/boxes/{id}", apiDeleteBox)        // Delete a box
	router.Get("/api/v1/boxes/{id}", apiGetBox)              // Get a specific box
	router.Post("/api/v1/boxes/{id}/events", apiCreateEvent) // Create a box event

	// Old paths, Deprecated.
	router.Get("/api/v1/box", DeprecatedRoute("use GET /api/v1/boxes instead")(apiGetBoxes))
	router.Post("/api/v1/box/new", DeprecatedRoute("use POST /api/v1/boxes instead")(apiCreateBox))
	router.Post("/api/v1/box/update", DeprecatedRoute("use PUT /api/v1/boxes/{id} instead")(apiReplaceBox))
	router.Delete("/api/v1/box/{id}", DeprecatedRoute("use DELETE /api/v1/boxes/{id} instead")(apiDeleteBox))
	router.Get("/api/v1/box/{id}", DeprecatedRoute("use GET /api/v1/boxes/{id} instead")(apiGetBox))
	router.Post("/api/v1/box/{id}/event", DeprecatedRoute("use POST /api/v1/boxes/{id}/event instead")(apiCreateEvent))

	listenOn := fmt.Sprintf(":%s", options.ApiPort)
	logger.Fatal(http.ListenAndServe(listenOn, router).Error())
}
