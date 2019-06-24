package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

// runUpdater is just for testing some events. This is only run to generate some
// random events.
func runUpdater() {

	// Generate a constant stream of events that get pushed
	// into the Broker's messages channel and are then broadcast
	// out to any clients that are attached. This will be replaced
	// with something else Kafka?
	go func() {
		var event Event
		for x := 0; ; x++ {
			if x >= len(boxes) {
				x = 0
			}

			// Create a little message to send to clients,
			// including the current time.
			t := time.Now()
			ft := fmt.Sprintf("%s", t.Format(time.RFC3339))
			if boxes[x].ID != statusBarID {
				event.ID = boxes[x].ID
				event.Color = "green"
				event.Message = fmt.Sprintf("the time is %s", ft)
				update(event)
			}

			if rand.Intn(30) == 1 {
				y := rand.Intn(len(boxes) - 1)
				if boxes[x].ID != statusBarID {
					switch rand.Intn(3) {
					case 0:
						event.Color = "red"
						event.Message = "PANIC! Red Alert"
					case 1:
						event.Color = "amber"
						event.Message = "OH NOES! Something's not quite right"
					case 2:
						event.Color = "grey"
						event.Message = "Meh not sure what to do now...."
					}

					event.ID = boxes[y].ID
					update(event)
				}
			}

			// Print a nice log message and sleep for 5s.
			time.Sleep(1 * time.Second)

		}
	}()
}

func update(event Event) {
	t := time.Now()
	ft := fmt.Sprintf("%s", t.Format(time.RFC3339))

	i, err := findBoxByID(event.ID)
	if err != nil {
		log.Print(err)

		return
	}

	boxes[i].LastMessage = event.Message
	boxes[i].LastUpdate = ft
	boxes[i].Color = event.Color

	if event.MaxTBU != "" {
		boxes[i].MaxTBU = event.MaxTBU
	}

	if event.ExpireAfter != "" {
		boxes[i].ExpireAfter = event.ExpireAfter
	}

	// Write json
	byteValue, err := json.Marshal(&boxes)
	if err != nil {
		log.Fatal(err)
	}
	err2 := ioutil.WriteFile(config.dataFile, byteValue, 0644)
	if err2 != nil {
		log.Fatal(err2)
	}

	event.Type = "updateBox"
	dataString, _ := json.Marshal(event)
	events.messages <- fmt.Sprint(string(dataString))
}
