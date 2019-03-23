package main

import (
	"fmt"
	"math/rand"
	"time"
)

func runUpdater(events *Broker) {

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
			events.messages <- fmt.Sprintf("%d,green, the time is %s", x, ft)
			boxes[x].LastMessage = fmt.Sprintf( "the time is %s", ft)
			boxes[x].LastUpdate = ft
			boxes[x].Color = "green"

			if rand.Intn(30) == 1 {
				y := rand.Intn( len(boxes) - 1 )
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
				events.messages <- fmt.Sprintf("%d,%s,%s", y, c, m)
				boxes[y].LastMessage = m
				boxes[y].LastUpdate = ft
				boxes[y].Color = c
			}

			// Print a nice log message and sleep for 5s.
			time.Sleep(1 * time.Second)

		}
	}()
}
