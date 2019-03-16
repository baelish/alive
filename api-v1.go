package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

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

func apiCreateBox(w http.ResponseWriter, r *http.Request) {}
func apiDeleteBox(w http.ResponseWriter, r *http.Request) {}

func runApi() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/", apiGetBoxes).Methods("GET")
	router.HandleFunc("/api/v1/{id}", apiGetBox).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", router))
}
