package main

import (
    "net/http"
    "log"
)


func main() {
    config := getConfiguration()
    log.Printf("%+v\n", config)
    createStaticContent(config.staticFilePath)
    runFrontPage(config.staticFilePath)
    events := runSse()
    runUpdater(events)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
