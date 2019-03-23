package main

import (
	"log"
	"net/http"
)

func main() {
	config := getConfiguration()
	log.Printf("%+v\n", config)
	createStaticContent(config.staticFilePath)
	getBoxes("/home/drosth/go/src/github.com/baelish/alive/test.json")
	runFrontPage(config.staticFilePath)
	events := runSse()
	runUpdater(events)
	go runAPI()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
