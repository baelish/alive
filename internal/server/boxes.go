package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/baelish/alive/api"

	"go.uber.org/zap"
)

var boxes []api.Box

type by func(p1, p2 *api.Box) bool

func (by by) Sort(boxes []api.Box) {
	bs := &boxSorter{
		boxes: boxes,
		by:    by,
	}

	sort.Sort(bs)
}

type boxSorter struct {
	boxes []api.Box
	by    func(p1, p2 *api.Box) bool
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

func addBox(box api.Box) (id string, err error) {
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

	var event api.Event
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

func deleteBox(id string, sendEvent bool) (found bool, deletedBox api.Box) {
	var newBoxes []api.Box

	for _, box := range boxes {
		if box.ID != id {
			newBoxes = append(newBoxes, box)
		} else {
			logger.Info("deleting box", zap.String("id", box.ID), zap.String("name", box.Name))
			deletedBox = box
			found = true
		}
	}

	boxes = newBoxes

	if sendEvent {
		event := api.Event{Type: "deleteBox", ID: id}
		if stringData, err := json.Marshal(event); err != nil {
			logger.Error(err.Error())
		} else {
			events.messages <- string(stringData)
		}
	}

	return found, deletedBox
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
					deleteBox(box.ID, true)

					continue
				}

			}

			if box.MaxTBU.Duration != 0 {
				if time.Since(lastUpdate) > box.MaxTBU.Duration && box.Status != api.NoUpdate {
					logger.Warn("no events for box", zap.String("id", box.ID))
					var event api.Event
					event.ID = box.ID
					event.Status = api.NoUpdate
					event.Message = fmt.Sprintf("No new updates for %s.", box.MaxTBU)
					event.Type = api.NoUpdate.String()
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
			for range 3 {
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
	Size := func(p1, p2 *api.Box) bool {
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

func update(event api.Event) {
	t := time.Now()
	i, err := findBoxByID(event.ID)

	if err != nil {
		logger.Error(err.Error())

		return
	}

	boxes[i].LastMessage = event.Message

	boxes[i].Messages = append(
		[]api.Message{
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

	if event.Type != api.NoUpdate.String() {
		boxes[i].LastUpdate = t
	}

	boxes[i].Status = event.Status
	if event.MaxTBU.Set {
		boxes[i].MaxTBU = event.MaxTBU
	}

	if event.ExpireAfter.Set {
		boxes[i].ExpireAfter = event.ExpireAfter
	}

	event.Type = "updateBox"
	dataString, _ := json.Marshal(event)
	events.messages <- fmt.Sprint(string(dataString))
}
