package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
)

// Box represents a single item on our monitoring screen.
type Box struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Size        string `json:"size"`
	Color       string `json:"color"`
	MaxTBU      string `json:"maxTBU"`
	LastUpdate  string `json:"lastUpdate"`
	LastMessage string `json:"lastMessage"`
}

var boxes []Box

type by func(p1, p2 *Box) bool

func (by by) Sort(boxes []Box) {
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
func getBoxes(jsonFile string) {
	byteValue, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(byteValue, &boxes)

	sortBoxes()

}

func sortBoxes() {
	Size := func(p1, p2 *Box) bool {
		if p1.Size == p2.Size {
			return false
		}
		size1 := sizeToNumber(p1.Size)
		size2 := sizeToNumber(p2.Size)
		return size1 > size2
	}
	by(Size).Sort(boxes)
}

func testBoxID(id string) bool {
	for _, box := range boxes {
		if box.ID == id {
			return true
		}
	}
	return false
}
