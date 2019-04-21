package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Event struct is used to stream events to dashboard.
type Event struct {
	ID      string `json:"id"`
	Color   string `json:"color"`
	Message string `json:"lastMessage"`
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
	update(event.ID, event.Color, event.Message)
	json.NewEncoder(w).Encode(event)
}

func apiCreateBox(w http.ResponseWriter, r *http.Request) {
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
	boxes = append(boxes, newBox)
	sortBoxes()
	newBoxPrint, _ := json.Marshal(newBox)
	log.Printf(string(newBoxPrint))
	json.NewEncoder(w).Encode(newBox)
	events.messages <- fmt.Sprintf("reloadPage")
}

func apiDeleteBox(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	found := false
	var newBoxes []Box
	for _, box := range boxes {
		if box.ID != params["id"] {
			newBoxes = append(newBoxes, box)
		} else {
			log.Printf("Deleting box %s as requested by %s", params["id"], r.RemoteAddr)
			found = true
		}
	}
	boxes = newBoxes
	if found == true {
		json.NewEncoder(w).Encode("deleted " + params["id"])
	} else {
		json.NewEncoder(w).Encode("not found")
	}
	events.messages <- fmt.Sprintf("reloadPage")
}

func runAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", apiGetBoxes).Methods("GET")
	router.HandleFunc("/api/v1/new", apiCreateBox).Methods("POST")
	router.HandleFunc("/api/v1/{id}", apiGetBox).Methods("GET")
	router.HandleFunc("/api/v1/{id}", apiDeleteBox).Methods("DELETE")
	router.HandleFunc("/api/v1/events/{id}", apiCreateEvent).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}
