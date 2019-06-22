package main

import (
	"fmt"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	expireBoxes()
	fmt.Fprintf(w, "<head><link rel='stylesheet' type='text/css' href='static/standard.css'/><script src='static/scripts.js'></script></head>")
	fmt.Fprintf(w, "<body onresize='rightSizeBigBox()' onload='rightSizeBigBox(); alertNoUpdate(\"status-bar\",20)'>")
	fmt.Fprintf(w, "<div id='big-box' class='big-box'>")

	for i := 0; i < len(boxes); i++ {
		fmt.Fprintf(w, "<div onclick='boxClick(this.id)' id='%s' class='%s %s box'>", boxes[i].ID, boxes[i].Color, boxes[i].Size)
		fmt.Fprintf(w, "<p class='title'>%s</p>", boxes[i].Name)
		fmt.Fprintf(w, "<p class='message'>%s</p>", boxes[i].LastMessage)
		fmt.Fprintf(w, "<p class='lastUpdated'>%s</p>", boxes[i].LastUpdate)
		fmt.Fprintf(w, "<p class='maxTBU'>%s</p>", boxes[i].MaxTBU)
		fmt.Fprintf(w, "</div>")
	}

	fmt.Fprintf(w, "</div>")
	fmt.Fprintf(w, "</body>")
}

func runFrontPage(staticFilePath string) {
	http.HandleFunc("/", handleRoot)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticFilePath))))
}
