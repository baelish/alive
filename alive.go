package main

import (
    "net/http"
    "log"
)


func main() {
    createStaticContent("./static")
    runFrontPage()
    runSse()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
