package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// Event struct is used to stream events to dashboard.
type Event struct {
	ID          string   `json:"id,omitempty"`
	After       string   `json:"after,omitempty"`
	Box         *Box     `json:"box,omitempty"`
	Status      Status   `json:"status,omitempty"`
	Message     string   `json:"lastMessage,omitempty"`
	ExpireAfter Duration `json:"expireAfter,omitempty"`
	MaxTBU      Duration `json:"maxTBU,omitempty"`
	Type        string   `json:"type"`
}

func apiGetBoxes(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(boxes)
	if err != nil {
		err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "could not get boxes"}`))
		if err != nil {
			log.Print(err)
		}
	}
}

func apiGetBox(w http.ResponseWriter, r *http.Request) {
	i, err := findBoxByID(chi.URLParam(r, "id"))

	if err != nil {
		err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "id not found"}`))
		if err != nil {
			log.Print(err)
		}

		return
	}

	err = json.NewEncoder(w).Encode(boxes[i])
	if err != nil {
		log.Print(err)
	}
}

func apiCreateEvent(w http.ResponseWriter, r *http.Request) {
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		jsonErr := json.NewEncoder(w).Encode(json.RawMessage(fmt.Sprintf(`{"message": "could not decode data received","error": "%s"}`, err.Error())))
		if jsonErr != nil {
			log.Print(jsonErr)
		}

		return
	}

	event.ID = chi.URLParam(r, "id")
	event.Type = "updateBox"
	update(event)
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		log.Print(err)
	}
}

func apiCreateBox(w http.ResponseWriter, r *http.Request) {
	var newBox Box
	err := json.NewDecoder(r.Body).Decode(&newBox)
	if err != nil {
		log.Printf("error decoding json: %s", err.Error())
		err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "could not decode data received"}`))
		if err != nil {
			log.Print(err)
		}

		return
	}

	id, err := addBox(newBox)
	if err != nil {
		json.NewEncoder(w).Encode(json.RawMessage(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
		log.Print(err)

		return
	}

	newBox.ID = id

	err = json.NewEncoder(w).Encode(newBox)
	if err != nil {
		log.Print(err)
	}
}

func apiStatus(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(json.RawMessage(`{"status": "ok"}`))
	if err != nil {
		log.Print(err)
	}
}

func apiUpdateBox(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	var newBox Box
	err := json.NewDecoder(r.Body).Decode(&newBox)
	if err != nil {
		err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "could not decode data received"}`))
		if err != nil {
			log.Print(err)
		}

		return
	}

	if newBox.ID == "" {
		err := json.NewEncoder(w).Encode(json.RawMessage(`{"error": "Cannot update box without an ID."}`))
		if err != nil {
			log.Print(err)
		}

		return
	}

	deleteBox(newBox.ID, false)
	newBox.LastUpdate = t
	boxes = append(boxes, newBox)
	sortBoxes()
	newBoxPrint, _ := json.Marshal(newBox)
	log.Print(string(newBoxPrint))
	err = json.NewEncoder(w).Encode(newBox)
	if err != nil {
		log.Print(err)
	}
	var event Event
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
		log.Print(err)
	}

}

func runAPI(_ context.Context) {
	if options.Debug {
		log.Print("Starting up API")
	}
	router := chi.NewRouter()
	router.Get("/health", apiStatus)
	router.Get("/api/v1", apiGetBoxes)                 // deprecate
	router.Get("/api/v1/", apiGetBoxes)                // deprecate
	router.Post("/api/v1/new", apiCreateBox)           // deprecate
	router.Post("/api/v1/update", apiUpdateBox)        // deprecate
	router.Delete("/api/v1/{id}", apiDeleteBox)        // deprecate
	router.Get("/api/v1/{id}", apiGetBox)              // deprecate
	router.Post("/api/v1/events/{id}", apiCreateEvent) // deprecate
	router.Get("/api/v1/box", apiGetBoxes)
	router.Post("/api/v1/box/new", apiCreateBox)
	router.Post("/api/v1/box/update", apiUpdateBox)
	router.Delete("/api/v1/box/{id}", apiDeleteBox)
	router.Get("/api/v1/box/{id}", apiGetBox)
	router.Post("/api/v1/box/{id}/event", apiCreateEvent)
	listenOn := fmt.Sprintf(":%s", options.ApiPort)
	log.Fatal(http.ListenAndServe(listenOn, router))
}
