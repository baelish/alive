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
	params := mux.Vars(r)
	i, err := findBoxByID(params["id"])

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
	params := mux.Vars(r)
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "could not decode data received"}`))
		if err != nil {
			log.Print(err)
		}

		return
	}

	event.ID = params["id"]
	event.Type = "updateBox"
	update(event)
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		log.Print(err)
	}
}

func apiCreateBox(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	ft := fmt.Sprintf("%s", t.Format(time.RFC3339))
	var newBox Box
	err := json.NewDecoder(r.Body).Decode(&newBox)
	if err != nil {
		err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "could not decode data received"}`))
		if err != nil {
			log.Print(err)
		}

		return
	}

	if !validateBoxSize(newBox.Size) {
		err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "Cannot create box, the ID requested already exists."}`))
		if err != nil {
			log.Print(err)
		}
		return
	}

	if newBox.ID != "" {
		if testBoxID(newBox.ID) {
			err = json.NewEncoder(w).Encode(json.RawMessage(`{"error": "Cannot create box, the ID requested already exists."}`))
			if err != nil {
				log.Print(err)
			}

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
	newBoxPrint, err := json.Marshal(newBox)
	if err != nil {
		log.Print(err)
	}
	log.Printf(string(newBoxPrint))
	err = json.NewEncoder(w).Encode(newBox)
	if err != nil {
		log.Print(err)
	}

	var event Event
	event.Type = "reloadPage"
	stringData, err := json.Marshal(event)
	if err != nil {
		log.Print(err)
	}
	events.messages <- fmt.Sprintf(string(stringData))
}

func apiStatus(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(json.RawMessage(`{"status": "ok"}`))
	if err != nil {
		log.Print(err)
	}
}

func apiUpdateBox(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	ft := fmt.Sprintf("%s", t.Format(time.RFC3339))
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

	deleteBox(newBox.ID)
	newBox.LastUpdate = ft
	boxes = append(boxes, newBox)
	sortBoxes()
	newBoxPrint, _ := json.Marshal(newBox)
	log.Printf(string(newBoxPrint))
	err = json.NewEncoder(w).Encode(newBox)
	if err != nil {
		log.Print(err)
	}
	var event Event
	event.Type = "reloadPage"
	stringData, _ := json.Marshal(event)
	events.messages <- fmt.Sprintf(string(stringData))
}

func apiDeleteBox(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var message json.RawMessage

	if deleteBox(params["id"]) {
		message = json.RawMessage(fmt.Sprintf(`{"info": "deleted box %s"}`, params["id"]))
		var event Event
		event.Type = "deleteBox"
		event.ID = params["id"]
		stringData, err := json.Marshal(event)
		if err != nil {
			log.Print(err)
		}
		events.messages <- fmt.Sprintf(string(stringData))
	} else {
		message = json.RawMessage(`{"error": "box not found"}`)
	}

	err := json.NewEncoder(w).Encode(message)
	if err != nil {
		log.Print(err)
	}

}

func runAPI() {
	router := mux.NewRouter()
	router.HandleFunc("/health", apiStatus).Methods("GET")
	router.HandleFunc("/api/v1", apiGetBoxes).Methods("GET")
	router.HandleFunc("/api/v1/", apiGetBoxes).Methods("GET")
	router.HandleFunc("/api/v1/new", apiCreateBox).Methods("POST")
	router.HandleFunc("/api/v1/update", apiUpdateBox).Methods("POST")
	router.HandleFunc("/api/v1/{id}", apiGetBox).Methods("GET")
	router.HandleFunc("/api/v1/{id}", apiDeleteBox).Methods("DELETE")
	router.HandleFunc("/api/v1/events/{id}", apiCreateEvent).Methods("POST")
	listenOn := fmt.Sprintf(":%s", config.apiPort)
	log.Fatal(http.ListenAndServe(listenOn, router))
}
