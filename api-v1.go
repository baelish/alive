package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Event struct is used to stream events to dashboard.
type Event struct {
	ID          string `json:"id"`
	Color       string `json:"color"`
	Message 		string `json:"lastMessage"`
}


func apiGetBoxes(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(boxes)
}


func apiGetBox(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, box := range boxes {
		if box.ID == params["id"] {
			json.NewEncoder(w).Encode(box)
			return
		}
	}
	json.NewEncoder(w).Encode(&Box{})
}


func apiCreateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var event Event
	_ = json.NewDecoder(r.Body).Decode(&event)
	event.ID = params["id"]
	update(event.ID, event.Color, event.Message)
	json.NewEncoder(w).Encode(event)
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
		json.NewEncoder(w).Encode("deleted" + params["id"])
	} else {
		json.NewEncoder(w).Encode("not found")
	}
	events.messages <- fmt.Sprintf("reloadPage")
}


func runAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", apiGetBoxes).Methods("GET")
	router.HandleFunc("/api/v1/{id}", apiGetBox).Methods("GET")
	router.HandleFunc("/api/v1/{id}", apiDeleteBox).Methods("DELETE")
	router.HandleFunc("/api/v1/events/{id}", apiCreateEvent).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}
