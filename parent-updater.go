package main

import (
	"context"
	"time"
)

func parentUpdater(ctx context.Context) {
	if options.Debug {
		logger.Info("starting parent update routine")
	}

	for {
		// Update parent

		select {
		case <-ctx.Done():
			if options.Debug {
				logger.Info("stopping parent update routine")
			}
			return

		case <-time.After(time.Duration(3 * time.Second)):
		}
	}
}
