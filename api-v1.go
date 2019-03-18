package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Event struct {
	Id          string `json:"id"`
	Color       string `json:"color"`
	Message 		string `json:lastMessage`
}


func apiGetBoxes(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(boxes)
}
func apiGetBox(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, box := range boxes {
		if box.Id == params["id"] {
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
	event.Id = params["id"]
	json.NewEncoder(w).Encode(event)
}
func apiDeleteBox(w http.ResponseWriter, r *http.Request) {}

func runApi() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", apiGetBoxes).Methods("GET")
	router.HandleFunc("/api/v1/{id}", apiGetBox).Methods("GET")
	router.HandleFunc("/api/v1/events/{id}", apiCreateEvent).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}
