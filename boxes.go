package main

import (
	"alive/api"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"go.uber.org/zap"
)

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

	if box.ID != nil {
		if testBoxID(*box.ID) {
			err = fmt.Errorf("a box already exists with that ID: %s", box.ID)
			return "", err
		}
	} else {
		for box.ID == nil {
			newID := randStringBytes(10)
			if !testBoxID(newID) {
				box.ID = &newID
			}
		}
	}

	box.LastUpdate = &t
	boxes = append(boxes, box)

	sortBoxes()

	logger.Info("creating a new box", zap.String("id", *box.ID))
	logger.Debug("box detail", logStructDetails(box)...)

	var event Event
	event.Type = "createBox"
	event.Box = &box

	i, err := findBoxByID(*box.ID)
	if err != nil {
		logger.Error(err.Error())
	}
	if i == 0 {
		event.After = Ptr("status-bar")
	} else {
		event.After = boxes[i-1].ID
	}

	stringData, err := json.Marshal(event)
	if err != nil {
		return "", (err)
	}
	events.messages <- string(stringData)

	return *box.ID, nil
}

func deleteBox(id string, event bool) bool {
	var newBoxes []Box
	var found bool

	for _, box := range boxes {
		if *box.ID != id {
			newBoxes = append(newBoxes, box)
		} else {
			logger.Info("deleting box", zap.String("id", *box.ID), zap.String("name", *box.Name))
			found = true
		}
	}

	boxes = newBoxes

	if event {
		var event Event
		event.Type = "deleteBox"
		event.ID = Ptr(id)
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

			if box.ExpireAfter != nil {
				if time.Since(*lastUpdate) > *DurationFromString(box.ExpireAfter) {
					logger.Info("deleting expired box", zap.String("id", *box.ID))
					_ = deleteBox(*box.ID, true)

					continue
				}

			}

			if box.MaxTBU != nil {
				if time.Since(*lastUpdate) > *DurationFromString(box.MaxTBU) && *box.Status != NoUpdate {
					logger.Warn("no events for box", zap.String("id", *box.ID))
					var event Event
					event.ID = box.ID
					event.Status = Ptr(NoUpdate)
					event.Message = Ptr(fmt.Sprintf("No new updates for %s.", box.MaxTBU))
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
		if box.ID == &id {
			return i, nil
		}
	}

	return -1, fmt.Errorf("could not find %s", id)
}

func sortBoxes() {
	Size := func(p1, p2 *Box) bool {
		if *p1.Size == *p2.Size {
			return *p1.Name < *p2.Name
		}

		return int(*p1.Size) > int(*p2.Size)
	}

	by(Size).Sort(boxes)
}

// Returns true if the ID exists false if it doesn't.
func testBoxID(id string) bool {
	for _, box := range boxes {
		if *box.ID == id {
			return true
		}
	}
	return false
}

func update(event Event) {
	t := time.Now()
	i, err := findBoxByID(*event.ID)

	if err != nil {
		logger.Error(err.Error())

		return
	}

	boxes[i].LastMessage = event.Message
	var msgStatus string
	if event.Status != nil {
		msgStatus = event.Status.String()
	} else {
		logger.Warn("event status is missing", zap.String("id", *event.ID))
		msgStatus = api.Grey.String()
	}
	var msgString string
	if event.Message != nil {
		msgString = *event.Message
	} else {
		logger.Warn("event message is missing", zap.String("id", *event.ID))
		msgString = ""
	}

	newMessage := &Message{
		Message:   msgString,
		Status:    msgStatus,
		TimeStamp: t,
	}
	boxes[i].Messages = Ptr(append([]Message{*newMessage}, *boxes[i].Messages...))

	const maxMessages = 30
	if messages := boxes[i].Messages; messages != nil && len(*messages) > maxMessages {
		trimmed := (*messages)[:maxMessages]
		boxes[i].Messages = &trimmed
	}

	if event.Type != NoUpdate.String() {
		boxes[i].LastUpdate = &t
	}

	boxes[i].Status = event.Status
	if event.MaxTBU != nil {
		boxes[i].MaxTBU = event.MaxTBU
	}

	if event.ExpireAfter != nil {
		boxes[i].ExpireAfter = event.ExpireAfter
	}

	event.Type = "updateBox"
	dataString, _ := json.Marshal(event)
	events.messages <- fmt.Sprint(string(dataString))
}
