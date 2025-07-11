package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/baelish/alive/api"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func apiGetBoxes(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(boxes)
	if err != nil {
		err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "could not get boxes"}`))
		if err != nil {
			logger.Error(err.Error())
		}
	}
}

func apiGetBox(w http.ResponseWriter, r *http.Request) {
	i, err := findBoxByID(chi.URLParam(r, "id"))

	if err != nil {
		err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "id not found"}`))
		if err != nil {
			logger.Error(err.Error())
		}

		return
	}

	err = json.NewEncoder(w).Encode(boxes[i])
	if err != nil {
		logger.Error(err.Error())
	}
}

func apiCreateEvent(w http.ResponseWriter, r *http.Request) {
	var event api.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(json.RawMessage(fmt.Sprintf(`{"message": "could not decode data received","error": "%s"}`, err.Error())))
		if jsonErr != nil {
			logger.Error(jsonErr.Error())
		}

		return
	}

	event.ID = chi.URLParam(r, "id")
	event.Type = "updateBox"
	update(event)
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		logger.Error(err.Error())
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func apiCreateBox(w http.ResponseWriter, r *http.Request) {
	var newBox api.Box
	err := json.NewDecoder(r.Body).Decode(&newBox)
	if err != nil {
		logger.Error(err.Error())
		err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "could not decode data received"}`))
		if err != nil {
			logger.Error(err.Error())
		}

		return
	}

	id, err := addBox(newBox)
	if err != nil {
		json.NewEncoder(w).Encode(json.RawMessage(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		logger.Error(err.Error())

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

func apiStatus(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(json.RawMessage(`{"status": "ok"}`))
	if err != nil {
		logger.Error(err.Error())
	}
}

func apiReplaceBox(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	var newBox api.Box
	err := json.NewDecoder(r.Body).Decode(&newBox)
	if err != nil {
		err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "could not decode data received"}`))
		if err != nil {
			logger.Error(err.Error())
		}

		return
	}

	if newBox.ID == "" {
		err := json.NewEncoder(w).Encode(json.RawMessage(`{"error": "Cannot update box without an ID."}`))
		if err != nil {
			logger.Error(err.Error())
		}

		return
	}

	deleteBox(newBox.ID, false)
	newBox.LastUpdate = t
	boxes = append(boxes, newBox)
	sortBoxes()
	logger.Info("updating box", zap.String("id", newBox.ID))
	logger.Debug("update details", logStructDetails(newBox)...)
	err = json.NewEncoder(w).Encode(newBox)
	if err != nil {
		logger.Error(err.Error())
	}
	var event api.Event
	event.Type = "reloadPage"
	stringData, _ := json.Marshal(event)
	events.messages <- string(stringData)
}

func apiDeleteBox(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if deleteBox(id, true) {
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
	router.Get("/api/v1/boxes", apiGetBoxes)                // Get all boxes
	router.Post("/api/v1/boxes", apiCreateBox)              // Create a new box
	router.Put("/api/v1/boxes/{id}", apiReplaceBox)         // Replace an existing box
	router.Delete("/api/v1/boxes/{id}", apiDeleteBox)       // Delete a box
	router.Get("/api/v1/boxes/{id}", apiGetBox)             // Get a specific box
	router.Post("/api/v1/boxes/{id}/event", apiCreateEvent) // Create a box event

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
