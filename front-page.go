package main

import (
	"fmt"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<head><link rel='stylesheet' type='text/css' href='static/standard.css'/><script src='static/scripts.js'></script></head>")
	fmt.Fprintf(w, "<body onresize='rightSizeBigBox()' onload='rightSizeBigBox(); keepalive()'>")
	fmt.Fprintf(w, "<div id='big-box' class='big-box'>")

	for i := 0; i < len(boxes); i++ {
		fmt.Fprintf(w, "<div onclick='boxClick(this.id)' id='%s' class='%s %s box'>", boxes[i].ID, boxes[i].Status, boxes[i].Size)
		fmt.Fprintf(w, "<p class='title'>%s</p>", boxes[i].Name)
		fmt.Fprintf(w, "<p class='message'>%s</p>", boxes[i].LastMessage)
		fmt.Fprintf(w, "<p class='lastUpdated'>%s</p>", boxes[i].LastUpdate)
		fmt.Fprintf(w, "<p class='maxTBU'>%s</p>", boxes[i].MaxTBU)
		fmt.Fprintf(w, "</div>")
	}

	fmt.Fprintf(w, "</div>")
	fmt.Fprintf(w, "</body>")
}

func handleStatus(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, `{"status":"ok"}`)
}

func runFrontPage() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc ("/health", handleStatus)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.staticFilePath))))
}
