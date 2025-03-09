package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
)

const timeFormat = "2006-01-02T15:04:05.000Z07:00"

var events *Broker
var wg sync.WaitGroup

func main() {

	ctx, cancel := context.WithCancel(context.Background())
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
			log.Println("SIGINT signal received, exiting gracefully")
			cancel()
		}
		<-signalChan // second signal, hard exit
		log.Println("2nd SIGINT signal received, hard exit")
		os.Exit(2)
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

	runDashboard(ctx)

	events = runSSE(ctx)

	runKeepalives(ctx)

	maintainBoxes(ctx)

	if options.Demo {
		runDemo(ctx)
	}

	runAPI(ctx)

	wg.Wait()
}
