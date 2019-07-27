package main

import (
	"encoding/json"
	"fmt"
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
				event.Status = "green"
				event.Message = fmt.Sprintf("the time is %s", ft)
				update(event)
			}

			if rand.Intn(2) == 1 {
				max := len(boxes) - 1
				if max > 0 {
					y := rand.Intn(max)
					if boxes[y].ID != statusBarID {
						switch rand.Intn(3) {
						case 0:
							event.Status = "red"
							event.Message = "PANIC! Red Alert"
						case 1:
							event.Status = "amber"
							event.Message = "OH NOES! Something's not quite right"
						case 2:
							event.Status = "grey"
							event.Message = "Meh not sure what to do now...."
						}

						event.ID = boxes[y].ID
						update(event)
					}
				}
			}
			time.Sleep(10 * time.Millisecond)

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

	if event.Type != missedStatusUpdate {
		boxes[i].LastUpdate = ft
	}

	boxes[i].Status = event.Status

	if event.MaxTBU != "" {
		boxes[i].MaxTBU = event.MaxTBU
	}

	if event.ExpireAfter != "" {
		boxes[i].ExpireAfter = event.ExpireAfter
	}

	event.Type = "updateBox"
	dataString, _ := json.Marshal(event)
	events.messages <- fmt.Sprint(string(dataString))
}
