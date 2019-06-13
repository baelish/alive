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
		var c, m string
		for x := 0; ; x++ {
			if x >= len(boxes) {
				x = 0
			}

			// Create a little message to send to clients,
			// including the current time.
			t := time.Now()
			ft := fmt.Sprintf("%s", t.Format(time.RFC3339))
			if boxes[x].ID != statusBarID {
				update(boxes[x].ID, "green", fmt.Sprintf("the time is %s", ft))
			}

			if rand.Intn(30) == 1 {
				y := rand.Intn(len(boxes) - 1)
				if boxes[x].ID != statusBarID {
					switch rand.Intn(3) {
					case 0:
						c = "red"
						m = "PANIC! Red Alert"
					case 1:
						c = "amber"
						m = "OH NOES! Something's not quite right"
					case 2:
						c = "grey"
						m = "Meh not sure what to do now...."
					}

					update(boxes[y].ID, c, m)
				}
			}

			// Print a nice log message and sleep for 5s.
			time.Sleep(1 * time.Second)

		}
	}()
}

func update(params ...string) {
	t := time.Now()
	ft := fmt.Sprintf("%s", t.Format(time.RFC3339))

	i, err := findBoxByID(params[0])
	if err != nil {
		log.Print(err)

		return
	}

	boxes[i].LastMessage = params[2]
	boxes[i].LastUpdate = ft
	boxes[i].Color = params[1]
	log.Printf("params: %d", len(params))

	if len(params) > 3 {
		log.Printf("params[3]: %s", params[3])
		boxes[i].MaxTBU = params[3]
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

	events.messages <- fmt.Sprintf("updateBox,%s,%s,%s,%s", params[0], boxes[i].Color, boxes[i].MaxTBU, boxes[i].LastMessage)
}
