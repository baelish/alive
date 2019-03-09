package main

import (
    "net/http"
    "log"
)


func main() {
    createStaticContent("./static")
    runFrontPage()
    events := runSse()
    runUpdater(events)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
