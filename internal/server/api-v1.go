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

func apiUpdateBox(w http.ResponseWriter, r *http.Request) {
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
	var message json.RawMessage
	id := chi.URLParam(r, "id")

	if deleteBox(id, true) {
		message = json.RawMessage(fmt.Sprintf(`{"info": "deleted box %s"}`, id))
	} else {
		message = json.RawMessage(`{"error": "box not found"}`)
	}

	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		logger.Error(err.Error())
	}

}

func runAPI(_ context.Context) {
	if options.Debug {
		logger.Info("Starting up API")
	}
	router := chi.NewRouter()
	router.Get("/health", apiStatus)
	// deprecate old box paths.
	router.Get("/api/v1/box", apiGetBoxes)                // move to boxes
	router.Post("/api/v1/box/new", apiCreateBox)          // move to boxes remove new
	router.Post("/api/v1/box/update", apiUpdateBox)       // move to boxes remove update, change to patch or put
	router.Delete("/api/v1/box/{id}", apiDeleteBox)       // move to boxes
	router.Get("/api/v1/box/{id}", apiGetBox)             // move to boxes
	router.Post("/api/v1/box/{id}/event", apiCreateEvent) // move to boxes/{id}/events
	listenOn := fmt.Sprintf(":%s", options.ApiPort)
	logger.Fatal(http.ListenAndServe(listenOn, router).Error())
}
