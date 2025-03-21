package main

import (
	"context"
	"encoding/json"
	"fmt"
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

type Message struct {
	Message   string `json:"message"`
	Status    string `json:"status"`
	TimeStamp string `json:"timeStamp"`
}

type Status int

const (
	Grey Status = iota
	Red
	Amber
	Green
	NoUpdate
)

func (s Status) String() string {
	return [...]string{
		"grey",
		"red",
		"amber",
		"green",
		"noUpdate",
	}[s]
}

func (s *Status) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	switch str {
	case "green":
		*s = Green
	case "grey":
		*s = Grey
	case "gray":
		*s = Grey
	case "noUpdate":
		*s = NoUpdate
	case "red":
		*s = Red
	case "amber":
		*s = Amber
	default:
		return fmt.Errorf("invalid status")
	}

	return nil
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

type BoxSize int

const (
	Dot BoxSize = iota
	Micro
	Dmicro
	Small
	Dsmall
	Medium
	Dmedium
	Large
	Dlarge
	Xlarge
)

func (bs BoxSize) String() string {
	return [...]string{
		"dot",
		"micro",
		"dmicro",
		"small",
		"dsmall",
		"medium",
		"dmedium",
		"large",
		"dlarge",
		"xlarge",
	}[bs]
}

func (bs *BoxSize) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	switch str {
	case "dot":
		*bs = Dot
	case "micro":
		*bs = Micro
	case "dmicro":
		*bs = Dmicro
	case "small":
		*bs = Small
	case "dsmall":
		*bs = Dsmall
	case "medium":
		*bs = Medium
	case "dmedium":
		*bs = Dmedium
	case "large":
		*bs = Large
	case "dlarge":
		*bs = Dlarge
	case "xlarge":
		*bs = Xlarge
	default:
		return fmt.Errorf("invalid box size")
	}

	return nil
}

func (bs BoxSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(bs.String())
}

// Box represents a single item on our monitoring screen.
type Box struct {
	ID          string    `json:"id"`
	Description string    `json:"description,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	Name        string    `json:"name"`
	Parent      string    `json:"parent,omitempty"`
	Size        BoxSize   `json:"size"`
	Status      Status    `json:"status"`
	ExpireAfter string    `json:"expireAfter,omitempty"`
	MaxTBU      string    `json:"maxTBU,omitempty"`
	Messages    []Message `json:"messages"`
	LastUpdate  string    `json:"lastUpdate"`
	LastMessage string    `json:"lastMessage"`
	Links       []Links   `json:"links"`
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

func addBox(box Box) (id string, err error) {
	t := time.Now()
	ft := t.Format(timeFormat)

	if box.ID != "" {
		if testBoxID(box.ID) {
			err = fmt.Errorf("a box already exists with that ID: %s", box.ID)
			return "", err
		}
	} else {
		for box.ID == "" || testBoxID(box.ID) {
			box.ID = randStringBytes(10)
		}

	}
	box.LastUpdate = ft
	boxes = append(boxes, box)

	sortBoxes()

	newBoxPrint, err := json.Marshal(box)
	if err != nil {
		return "", (err)
	}
	log.Printf("creating new box with these details:'%s'", string(newBoxPrint))
	//	time.Sleep(100 * time.Millisecond)

	var event Event
	event.Type = "createBox"
	event.Box = &box

	i, err := findBoxByID(box.ID)
	if err != nil {
		log.Printf("couldn't find box (%s)", err.Error())
	}
	if i == 0 {
		event.After = "status-bar"
	} else {
		event.After = boxes[i-1].ID
	}

	stringData, err := json.Marshal(event)
	if err != nil {
		return "", (err)
	}
	events.messages <- string(stringData)

	return box.ID, nil
}

func deleteBox(id string, event bool) bool {
	var newBoxes []Box
	var found bool

	for _, box := range boxes {
		if box.ID != id {
			newBoxes = append(newBoxes, box)
		} else {
			log.Printf("Deleting box %s (%s)", id, box.Name)
			found = true
		}
	}

	boxes = newBoxes

	if event {
		var event Event
		event.Type = "deleteBox"
		event.ID = id
		stringData, err := json.Marshal(event)
		if err != nil {
			log.Print(err)
		}
		events.messages <- string(stringData)
	}

	return found
}

// Find any boxes that have expired and delete them, find any boxes which have
// not had timely updates and update their status. Also saves box file
// periodically or on exit.
func maintainBoxes(ctx context.Context) {
	if options.Debug {
		log.Print("Starting box maintenance routine")
	}
	var err error
	var lastSave time.Time
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
					_ = deleteBox(box.ID, true)

					continue
				}

			}

			if box.MaxTBU != "0" && box.MaxTBU != "" {
				alertAfter, err := strconv.Atoi(box.MaxTBU)

				if err != nil {
					log.Println(err)
				} else if lastUpdate.Add(time.Second*time.Duration(alertAfter)).Before(time.Now()) && box.Status != NoUpdate {
					log.Printf("no events for box %s", box.ID)
					var event Event
					event.ID = box.ID
					event.Status = NoUpdate
					event.Message = fmt.Sprintf("No new updates for %ss.", box.MaxTBU)
					event.Type = NoUpdate.String()
					update(event)

					continue
				}

			}
		}
		// Write json
		if time.Since(lastSave) > time.Duration(1*time.Minute) {
			log.Print("Saving data file")
			err = saveBoxFile()
			if err != nil {
				log.Printf("Error saving data file (%s)", err.Error())
			} else {
				lastSave = time.Now()
			}
		}

		select {
		case <-ctx.Done():
			log.Printf("Saving data file before exit")
			for t := 0; t < 3; t++ {
				err = saveBoxFile()
				if err != nil {
					log.Printf("Error saving box file (%s)", err.Error())
				}
			}

			return

		case <-time.After(time.Duration(1 * time.Second)):
		}
	}
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

func sortBoxes() {
	Size := func(p1, p2 *Box) bool {
		if p1.Size == p2.Size {
			return p1.Name < p2.Name
		}

		return int(p1.Size) > int(p2.Size)
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

func update(event Event) {
	t := time.Now()
	ft := t.Format(timeFormat)
	i, err := findBoxByID(event.ID)

	if err != nil {
		log.Print(err)

		return
	}

	boxes[i].LastMessage = event.Message

	boxes[i].Messages = append(
		[]Message{
			{
				Message:   event.Message,
				Status:    event.Status.String(),
				TimeStamp: ft,
			},
		},
		boxes[i].Messages...,
	)

	m := 30
	if len(boxes[i].Messages) > m {
		boxes[i].Messages = boxes[i].Messages[:m]
	}

	if event.Type != NoUpdate.String() {
		boxes[i].LastUpdate = ft
	}

	boxes[i].Status = event.Status

	if event.MaxTBU != "" {
		boxes[i].MaxTBU = event.MaxTBU
	}

	if event.ExpireAfter != "" {
		boxes[i].ExpireAfter = event.ExpireAfter
	}

	event.Type = "updateBox"
	dataString, _ := json.Marshal(event)
	events.messages <- fmt.Sprint(string(dataString))
}
