package main

import (
	"fmt"
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
	maintainBoxes()

	if config.updater {
		runUpdater()
	}
	go runAPI()
	listenOn := fmt.Sprintf(":%s", config.sitePort)
	log.Fatal(http.ListenAndServe(listenOn, nil))
}
