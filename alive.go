package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type Box struct {
	Name  string `json:"name"`
	Size  string `json:"size"`
	Color string `json:"color"`
}

type By func(p1, p2 *Box) bool

func (by By) Sort(boxes []Box) {
	bs := &boxSorter{
		boxes: boxes,
		by:    by,
	}
	sort.Sort(bs)
}

type boxSorter struct {
	boxes []Box
	by    func(p1, p2 *Box) bool
}

func (s *boxSorter) Len() int {
	return len(s.boxes)
}

func (s *boxSorter) Swap(i, j int) {
	s.boxes[i], s.boxes[j] = s.boxes[j], s.boxes[i]
}

func (s *boxSorter) Less(i, j int) bool {
	return s.by(&s.boxes[i], &s.boxes[j])
}

func sizeToNumber(size string) int {
	switch size {
	case "micro":
		return 10
	case "dmicro":
		return 20
	case "small":
		return 30
	case "dsmall":
		return 40
	case "medium":
		return 50
	case "dmedium":
		return 60
	case "large":
		return 70
	case "dlarge":
		return 80
	case "xlarge":
		return 90
	case "dxlarge":
		return 100
	default:
		return 0
	}
}

// Loads Json from a file and returns Boxes sorted by size (Largest first)
func loadJsonFromFile(jsonFile string) []Box {
	byteValue, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	var boxes []Box
	json.Unmarshal(byteValue, &boxes)

	Size := func(p1, p2 *Box) bool {
		if p1.Size == p2.Size {
			return false
		}
		size1 := sizeToNumber(p1.Size)
		size2 := sizeToNumber(p2.Size)
		return size1 > size2
	}

	By(Size).Sort(boxes)

	return boxes
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	boxes := loadJsonFromFile("/home/drosth/go/src/github.com/baelish/alive/test.json")
	fmt.Fprintf(w, "<head><link rel='stylesheet' type='text/css' href='static/standard.css'/><script src='static/scripts.js'></script></head>")
	fmt.Fprintf(w, "<div class='big-box'>")
	for i := 0; i < len(boxes); i++ {
		fmt.Fprintf(w, "<div onclick='boxClick(this.id)' id='%d' class='%s %s box'>%s</div>", i, boxes[i].Color, boxes[i].Size, boxes[i].Name)
	}
	fmt.Fprintf(w, "</div>")
}

func main() {
	http.HandleFunc("/", handleRoot)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
