package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// Links describes a URL with a name.
type Links struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Message struct {
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	TimeStamp time.Time `json:"timeStamp"`
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

type Duration struct {
	time.Duration
	bool
}

// This is a custom UnmarshalJSON() for a time.Duration that can have a null
// value ("") as well as being backward compatible with supplying a string with
// the number of seconds.
func (d *Duration) UnmarshalJSON(b []byte) (err error) {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		d.bool = true
		return nil
	case string:
		var err error
		if value == "" {
			d.Duration = 0
			d.bool = false
			return nil
		}
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			i, err2 := strconv.Atoi(value)
			if err2 != nil {
				return errors.Join(err, err2)
			}
			d.Duration = time.Duration(i) * time.Second
			d.bool = true
			return nil
		} else {
			d.bool = true
			return nil
		}
	default:
		return fmt.Errorf("invalid type for duration (%T)", v)
	}
}

// This is a custom MarshalJSON() for a time.Duration that can have a null value
// ("")
func (d Duration) MarshalJSON() (b []byte, err error) {
	if d.bool {
		return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
	} else {
		return []byte(`""`), nil
	}
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
	ExpireAfter Duration  `json:"expireAfter,omitempty"`
	MaxTBU      Duration  `json:"maxTBU,omitempty"`
	Messages    []Message `json:"messages"`
	LastUpdate  time.Time `json:"lastUpdate"`
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
	box.LastUpdate = t
	boxes = append(boxes, box)

	sortBoxes()

	logger.Info("creating a new box", zap.String("id", box.ID))
	logger.Debug("box detail", logStructDetails(box)...)

	var event Event
	event.Type = "createBox"
	event.Box = &box

	i, err := findBoxByID(box.ID)
	if err != nil {
		logger.Error(err.Error())
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
			logger.Info("deleting box", zap.String("id", box.ID), zap.String("name", box.Name))
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
			logger.Error(err.Error())
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
		logger.Info("Starting box maintenance routine")
	}
	var err error
	var lastSave time.Time
	for {
		for _, box := range boxes {
			if box.LastUpdate.IsZero() {
				continue
			}

			lastUpdate := box.LastUpdate

			if err != nil {
				logger.Error(err.Error())

				continue
			}

			if box.ExpireAfter.Duration != 0 {
				if time.Since(lastUpdate) > box.ExpireAfter.Duration {
					logger.Info("deleting expired box", zap.String("id", box.ID))
					_ = deleteBox(box.ID, true)

					continue
				}

			}

			if box.MaxTBU.Duration != 0 {
				if time.Since(lastUpdate) > box.MaxTBU.Duration && box.Status != NoUpdate {
					logger.Warn("no events for box", zap.String("id", box.ID))
					var event Event
					event.ID = box.ID
					event.Status = NoUpdate
					event.Message = fmt.Sprintf("No new updates for %s.", box.MaxTBU)
					event.Type = NoUpdate.String()
					update(event)

					continue
				}
			}
		}
		// Write json
		if time.Since(lastSave) > time.Duration(1*time.Minute) {
			logger.Info("Saving data file")
			err = saveBoxFile()
			if err != nil {
				logger.Error(err.Error())
			} else {
				lastSave = time.Now()
			}
		}

		select {
		case <-ctx.Done():
			logger.Info("Saving data file before exit")
			for t := 0; t < 3; t++ {
				err = saveBoxFile()
				if err != nil {
					logger.Error(err.Error())
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
	i, err := findBoxByID(event.ID)

	if err != nil {
		logger.Error(err.Error())

		return
	}

	boxes[i].LastMessage = event.Message

	boxes[i].Messages = append(
		[]Message{
			{
				Message:   event.Message,
				Status:    event.Status.String(),
				TimeStamp: t,
			},
		},
		boxes[i].Messages...,
	)

	m := 30
	if len(boxes[i].Messages) > m {
		boxes[i].Messages = boxes[i].Messages[:m]
	}

	if event.Type != NoUpdate.String() {
		boxes[i].LastUpdate = t
	}

	boxes[i].Status = event.Status
	if event.MaxTBU.bool {
		boxes[i].MaxTBU = event.MaxTBU
	}

	if event.ExpireAfter.bool {
		boxes[i].ExpireAfter = event.ExpireAfter
	}

	event.Type = "updateBox"
	dataString, _ := json.Marshal(event)
	events.messages <- fmt.Sprint(string(dataString))
}
