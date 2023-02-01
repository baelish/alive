package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"time"
)

// Links describes a URL with a name.
type Links struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Box represents a single item on our monitoring screen.
type Box struct {
	ID          string  `json:"id"`
	Description string  `json:"description,omitempty"`
	DisplayName string  `json:"displayName,omitempty"`
	Name        string  `json:"name"`
	Size        string  `json:"size"`
	Status      string  `json:"status"`
	ExpireAfter string  `json:"expireAfter,omitempty"`
	MaxTBU      string  `json:"maxTBU,omitempty"`
	LastUpdate  string  `json:"lastUpdate"`
	LastMessage string  `json:"lastMessage"`
	Links       []Links `json:"links"`
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
	case "status":
		return 110
	default:
		return 0
	}
}

func deleteBox(id string) bool {
	var newBoxes []Box
	var found bool

	for _, box := range boxes {
		if box.ID != id {
			newBoxes = append(newBoxes, box)
		} else {
			log.Printf("Deleting box %s", id)
			found = true
		}
	}

	boxes = newBoxes

	return found
}

// Find any boxes that have expired and delete them. Also find any boxes which
// have not had timely updates and update their status.
func maintainBoxes() {
	go func() {
		for {
			for _, box := range boxes {
				if box.LastUpdate == "" {
					continue
				}

				lastUpdate, err := time.Parse(time.RFC3339, box.LastUpdate)

				if err != nil {
					log.Println(err)

					continue
				}

				if box.ExpireAfter != "0" && box.ExpireAfter != "" {
					expireAfter, err := strconv.Atoi(box.ExpireAfter)

					if err != nil {
						log.Println(err)
					} else if lastUpdate.Add(time.Second * time.Duration(expireAfter)).Before(time.Now()) {
						log.Printf("deleting expired box %s", box.ID)
						_ = deleteBox(box.ID)
						var event Event
						event.ID = box.ID
						event.Type = "deleteBox"
						stringData, _ := json.Marshal(event)
						events.messages <- fmt.Sprintf(string(stringData))

						continue
					}

				}

				if box.MaxTBU != "0" && box.MaxTBU != "" {
					alertAfter, err := strconv.Atoi(box.MaxTBU)

					if err != nil {
						log.Println(err)
					} else if lastUpdate.Add(time.Second*time.Duration(alertAfter)).Before(time.Now()) && box.Status != missedStatusUpdate {
						log.Printf("no events for box %s", box.ID)
						var event Event
						event.ID = box.ID
						event.Status = missedStatusUpdate
						event.Message = fmt.Sprintf("No new updates for %ss. <br /> Last message: %s on %s", box.MaxTBU, box.LastMessage, box.LastUpdate)
						event.Type = missedStatusUpdate
						update(event)

						continue
					}

				}
			}

			// Write json
			byteValue, err := json.Marshal(&boxes)
			if err != nil {
				log.Fatal(err)
			}

			err = ioutil.WriteFile(options.DataFile, byteValue, 0644)

			if err != nil {
				log.Fatal(err)
			}

			// Sleep for 1s.
			time.Sleep(1 * time.Second)
		}
	}()
}

// Find a box in the boxes array, supply the box ID, will return the array id
func findBoxByID(id string) (int, error) {
	for i, box := range boxes {
		if box.ID == id {
			return i, nil
		}
	}

	return -1, fmt.Errorf("could not find %s", id)
}

// Loads Json from a file and returns Boxes sorted by size (Largest first)
func getBoxes() {
	byteValue, err := ioutil.ReadFile(options.DataFile)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(byteValue, &boxes)
	if err != nil {
		log.Fatal(err)
	}

	sortBoxes()

}

func sortBoxes() {
	Size := func(p1, p2 *Box) bool {
		if p1.Size == p2.Size {
			return p1.Name < p2.Name
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
