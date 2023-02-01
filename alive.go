package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	goflags "github.com/jessevdk/go-flags"
)

const statusBarID = "status-bar"

var events *Broker

type Options struct {
	ApiPort       string `long:"api-port" description:"The port to use for api calls" default:"8081"`
	SitePort      string `short:"p" long:"port" description:"The port to use for the dashboard" default:"8080"`
	Updater       bool   `long:"updater" description:"?"`
	DefaultStatic bool   `long:"default-static" description:"Use default static content"`
	DataFile      string `short:"f" long:"data-file" description:"Data file location (default: $HOME/.alive/data.json)"`
	StaticPath    string `long:"static-path" description:"Path to store static files (default: $HOME/.alive/static)"`
}

var options Options

func processOptions() {
	goflagParser := goflags.NewParser(&options, goflags.Default)

	if _, err := goflagParser.Parse(); err != nil {
		if flagsErr, ok := err.(*goflags.Error); ok && flagsErr.Type == goflags.ErrHelp {
			os.Exit(0)
		} else {
			log.Panicf("Parse failed: %v", err)
		}
	}

	if options.DataFile == "" {
		options.DataFile = filepath.Clean(fmt.Sprintf("%s/.alive/data.json", os.Getenv("HOME")))
	}
	if options.StaticPath == "" {
		options.StaticPath = filepath.Clean(fmt.Sprintf("%s/.alive/static", os.Getenv("HOME")))
	}

	log.Printf("%+v\n", options)
}

func main() {
	processOptions()
	createStaticContent()
	getBoxes()
	runPages()
	events = runSse()
	runKeepalives()
	maintainBoxes()

	if options.Updater {
		runUpdater()
	}

	go runAPI()

	listenOn := fmt.Sprintf(":%s", options.SitePort)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
