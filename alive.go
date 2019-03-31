package main

import (
	"log"
	"net/http"
)

var events *Broker
var config *Config
var sizes = []string{"micro", "dmicro", "small", "dsmall", "medium", "dmedium", "large", "dlarge", "xlarge", "dxlarge"}

func main() {
	config = getConfiguration()
	log.Printf("%+v\n", config)
	createStaticContent(config.staticFilePath)
	getBoxes(config.dataFile)
	runFrontPage(config.staticFilePath)
	events = runSse()
	if config.updater {
		runUpdater()
	}
	go runAPI()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
