package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.000Z07:00"

var events *Broker

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer func() {
		log.Printf("Running cleanup")
		cancel()
	}()

	processOptions()

	if options.Demo {
		tempDir, err := os.MkdirTemp(os.TempDir(), "alive-*.tmp")
		if err != nil {
			log.Panicf("Unable to create a temporary directory")
		}
		defer os.RemoveAll(tempDir)

		log.Printf("Running demo using temporary files in %s", tempDir)

		options.DataPath = filepath.Clean(fmt.Sprintf("%s/data", tempDir))
		options.StaticPath = filepath.Clean(fmt.Sprintf("%s/static", tempDir))
	}

	if options.DataPath == "" {
		options.DataPath = filepath.Clean(fmt.Sprintf("%s/.alive/data", os.Getenv("HOME")))
	}

	if options.StaticPath == "" {
		options.StaticPath = filepath.Clean(fmt.Sprintf("%s/.alive/static", os.Getenv("HOME")))
	}

	log.Printf("%+v\n", options)

	createStaticContent()
	createDataFiles()
	getBoxesFromDataFile()

	events = runSSE(ctx)

	go runDashboard(ctx)
	go runAPI(ctx)

	go runKeepalives(ctx)
	go maintainBoxes(ctx)

	if options.ParentUrl != "" {
		go parentUpdater(ctx)
	}

	if options.Demo {
		go runDemo(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(1 * time.Second)):
		}
	}
}
