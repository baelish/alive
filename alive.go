package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const statusBarID = "status-bar"

var events *Broker
var config *Config

func init() {
	config = getConfiguration(os.Args)
}

func main() {
	log.Printf("%+v\n", config)
	createStaticContent()
	getBoxes()
	runPages()
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
