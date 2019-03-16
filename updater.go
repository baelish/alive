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
		var t int
		for i := 0; ; i++ {
			if i > 55 {
				i = 0
			}

			// Create a little message to send to clients,
			// including the current time.
			events.messages <- fmt.Sprintf("%d,green, the time is %v", i, time.Now())
			if rand.Intn(30) == 1 {
				t = rand.Intn(55)
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
				events.messages <- fmt.Sprintf("%d,%s,%s", t, c, m)
			}

			// Print a nice log message and sleep for 5s.
			time.Sleep(1 * time.Second)

		}
	}()
}
