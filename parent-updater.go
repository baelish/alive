package main

import (
	"context"
	"log"
	"time"
)

// Send keepalives to the status bar.
func parentUpdater(ctx context.Context) {
	if options.Debug {
		log.Print("starting parent update routine")
	}

	for {
		// Update parent

		select {
		case <-ctx.Done():
			if options.Debug {
				log.Print("stopping parent update routine")
			}
			return

		case <-time.After(time.Duration(3 * time.Second)):
		}
	}
}
