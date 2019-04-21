package main

import (
	"log"
	"net/http"
)

const statusBarID = "status-bar"

var events *Broker
var config *Config
var sizes = []string{"micro", "dmicro", "small", "dsmall", "medium", "dmedium", "large", "dlarge", "xlarge", "dxlarge", "status"}

func main() {
	config = getConfiguration()
	log.Printf("%+v\n", config)
	createStaticContent(config.staticFilePath)
	getBoxes(config.dataFile)
	runFrontPage(config.staticFilePath)
	events = runSse()
	runKeepalives()

	if config.updater {
		runUpdater()
	}
	go runAPI()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
