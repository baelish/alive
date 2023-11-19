package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
)

const statusBarID = "status-bar"
const timeFormat = "2006-01-02T15:04:05.000Z07:00"

var events *Broker

func main() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		log.Printf("Running cleanup")
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			log.Println("SIGINT signal received, exiting")
			cancel()
		case <-ctx.Done():
			return
		}
		<-signalChan // second signal, hard exit
		os.Exit(1)
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
	go runDashboard(ctx)
	events = runSSE(ctx)
	runKeepalives(ctx)
	maintainBoxes(ctx)

	if options.Demo {
		runDemo(ctx)
	}

	go runAPI(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}
