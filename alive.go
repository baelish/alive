package main

import (
	"fmt"
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<head><link rel='stylesheet' type='text/css' href='css/standard.css'/></head>")
	fmt.Fprintf(w, "<div class='green small box'>Hi there, I love %s!</div>", r.Proto)
	fmt.Fprintf(w, "<div class='green medium box'>Hi there, I love %s!</div>", r.URL.Path[1:])
	fmt.Fprintf(w, "<div class='green small box'>Hi there, I love %s!</div>", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
