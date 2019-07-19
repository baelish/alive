package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Event struct is used to stream events to dashboard.
type Event struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	ExpireAfter string `json:"expireAfter"`
	Message     string `json:"lastMessage"`
	MaxTBU      string `json:"maxTBU"`
	Type        string `json:"type"`
}

func apiGetBoxes(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(boxes)
}

func apiGetBox(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	i, err := findBoxByID(params["id"])
	if err != nil {
		json.NewEncoder(w).Encode(json.RawMessage(`{"error": "id not found"}`))
		return
	}
	json.NewEncoder(w).Encode(boxes[i])
}

func apiCreateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var event Event
	_ = json.NewDecoder(r.Body).Decode(&event)
	event.ID = params["id"]
	event.Type = "updateBox"
	update(event)
	json.NewEncoder(w).Encode(event)
}

func apiCreateBox(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	ft := fmt.Sprintf("%s", t.Format(time.RFC3339))
	var newBox Box
	_ = json.NewDecoder(r.Body).Decode(&newBox)
	if newBox.ID != "" {
		if testBoxID(newBox.ID) {
			json.NewEncoder(w).Encode("Cannot create box, the ID requested already exists.")
			return
		}
	} else {
		for newBox.ID == "" || testBoxID(newBox.ID) {
			newBox.ID = randStringBytes(10)
		}

	}
	newBox.LastUpdate = ft
	boxes = append(boxes, newBox)
	sortBoxes()
	newBoxPrint, _ := json.Marshal(newBox)
	log.Printf(string(newBoxPrint))
	json.NewEncoder(w).Encode(newBox)
	var event Event
	event.Type = "reloadPage"
	stringData, _ := json.Marshal(event)
	events.messages <- fmt.Sprintf(string(stringData))
}

func apiDeleteBox(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if deleteBox(params["id"]) {
		json.NewEncoder(w).Encode("deleted " + params["id"])
	} else {
		json.NewEncoder(w).Encode("not found")
	}
	var event Event
	event.Type = "deleteBox"
	event.ID = params["id"]
	stringData, _ := json.Marshal(event)
	events.messages <- fmt.Sprintf(string(stringData))
}

func runAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", apiGetBoxes).Methods("GET")
	router.HandleFunc("/api/v1/new", apiCreateBox).Methods("POST")
	router.HandleFunc("/api/v1/{id}", apiGetBox).Methods("GET")
	router.HandleFunc("/api/v1/{id}", apiDeleteBox).Methods("DELETE")
	router.HandleFunc("/api/v1/events/{id}", apiCreateEvent).Methods("POST")
	listenOn := fmt.Sprintf(":%s", config.apiPort)
	log.Fatal(http.ListenAndServe(listenOn, router))
}
