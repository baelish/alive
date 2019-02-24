package main

import (
    "encoding/json"
	"fmt"
    "io/ioutil"
	"log"
	"net/http"
)

type Box struct {
    Name string `json:"name"`
    Size string `json:"size"`
    Color string `json:"color"`
}

type Boxes struct {
    Boxes []Box `json:"boxes"`
}

func loadJsonFromFile(jsonFile string) Boxes {
    byteValue, err := ioutil.ReadFile(jsonFile)
    if err != nil {
        log.Fatal(err)
    }
    var boxes Boxes
    json.Unmarshal(byteValue, &boxes)
    return boxes
}


func handleRoot(w http.ResponseWriter, r *http.Request) {
    boxes := loadJsonFromFile("/home/drosth/go/src/github.com/baelish/alive/test.json")
	fmt.Fprintf(w, "<head><meta http-equiv='refresh' content='5'><link rel='stylesheet' type='text/css' href='css/standard.css'/></head>")
	fmt.Fprintf(w, "<div class='big-box'>")
    for i := 0; i < len(boxes.Boxes); i++ {
	    fmt.Fprintf(w, "<div class='%s %s box'>%s</div>", boxes.Boxes[i].Color, boxes.Boxes[i].Size, boxes.Boxes[i].Name)
    }
	fmt.Fprintf(w, "</div>")
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
