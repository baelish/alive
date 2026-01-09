package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/baelish/alive/api"

	"go.uber.org/zap"
)

// BoxStore provides thread-safe access to boxes
type BoxStore struct {
	mu    sync.RWMutex
	boxes []api.Box
}

// Global box store instance
var boxStore = &BoxStore{
	boxes: make([]api.Box, 0),
}

// GetAll returns a copy of all boxes (thread-safe read)
func (bs *BoxStore) GetAll() []api.Box {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	// Return a copy to prevent external modifications
	result := make([]api.Box, len(bs.boxes))
	copy(result, bs.boxes)
	return result
}

// GetByID returns a copy of a box by ID (thread-safe read)
func (bs *BoxStore) GetByID(id string) (*api.Box, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	for i := range bs.boxes {
		if bs.boxes[i].ID == id {
			// Return a copy
			box := bs.boxes[i]
			return &box, nil
		}
	}
	return nil, fmt.Errorf("could not find %s", id)
}

// FindIndexByID returns the index of a box by ID (thread-safe read)
func (bs *BoxStore) FindIndexByID(id string) (int, error) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	for i := range bs.boxes {
		if bs.boxes[i].ID == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("could not find %s", id)
}

// Exists checks if a box with given ID exists (thread-safe read)
func (bs *BoxStore) Exists(id string) bool {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	for i := range bs.boxes {
		if bs.boxes[i].ID == id {
			return true
		}
	}
	return false
}

// Add adds a new box (thread-safe write)
func (bs *BoxStore) Add(box api.Box) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	// Check if ID already exists
	for i := range bs.boxes {
		if bs.boxes[i].ID == box.ID {
			return fmt.Errorf("a box already exists with that ID: %s", box.ID)
		}
	}

	bs.boxes = append(bs.boxes, box)
	bs.sortUnsafe()
	return nil
}

// Delete removes a box by ID (thread-safe write)
func (bs *BoxStore) Delete(id string) (found bool, deletedBox api.Box) {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	for i := range bs.boxes {
		if bs.boxes[i].ID == id {
			deletedBox = bs.boxes[i]
			// Remove from slice
			bs.boxes = append(bs.boxes[:i], bs.boxes[i+1:]...)
			return true, deletedBox
		}
	}
	return false, api.Box{}
}

// Update modifies an existing box (thread-safe write)
func (bs *BoxStore) Update(id string, updateFn func(*api.Box)) error {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	for i := range bs.boxes {
		if bs.boxes[i].ID == id {
			updateFn(&bs.boxes[i])
			return nil
		}
	}
	return fmt.Errorf("could not find box %s", id)
}

// ForEach iterates over all boxes with a read lock
func (bs *BoxStore) ForEach(fn func(api.Box) bool) {
	bs.mu.RLock()
	defer bs.mu.RUnlock()

	for _, box := range bs.boxes {
		if !fn(box) {
			break
		}
	}
}

// Len returns the number of boxes (thread-safe read)
func (bs *BoxStore) Len() int {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	return len(bs.boxes)
}

// sortUnsafe sorts boxes (must be called with lock held)
func (bs *BoxStore) sortUnsafe() {
	Size := func(p1, p2 *api.Box) bool {
		if p1.Size == p2.Size {
			return p1.Name < p2.Name
		}
		return int(p1.Size) > int(p2.Size)
	}
	by(Size).Sort(bs.boxes)
}

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
		if boxStore.Exists(box.ID) {
			err = fmt.Errorf("a box already exists with that ID: %s", box.ID)
			return "", err
		}
	} else {
		for box.ID == "" || boxStore.Exists(box.ID) {
			box.ID = randStringBytes(10)
		}
	}

	box.LastUpdate = t

	// Add to store (thread-safe)
	if err := boxStore.Add(box); err != nil {
		return "", err
	}

	logger.Info("creating a new box", zap.String("id", box.ID))
	logger.Debug("box detail", logStructDetails(box)...)

	var event api.Event
	event.Type = "createBox"
	event.Box = &box

	i, err := boxStore.FindIndexByID(box.ID)
	if err != nil {
		logger.Error(err.Error())
	}
	if i == 0 {
		event.After = "status-bar"
	} else {
		// Get the box before this one
		allBoxes := boxStore.GetAll()
		if i > 0 && i <= len(allBoxes) {
			event.After = allBoxes[i-1].ID
		}
	}

	stringData, err := json.Marshal(event)
	if err != nil {
		return "", err
	}
	events.messages <- string(stringData)

	return box.ID, nil
}

func deleteBox(id string, sendEvent bool) (found bool, deletedBox api.Box) {
	// Delete from store (thread-safe)
	found, deletedBox = boxStore.Delete(id)

	if found {
		logger.Info("deleting box", zap.String("id", deletedBox.ID), zap.String("name", deletedBox.Name))
	}

	if sendEvent && found {
		event := api.Event{Type: "deleteBox", ID: id}
		if stringData, err := json.Marshal(event); err != nil {
			logger.Error(err.Error())
		} else {
			events.messages <- string(stringData)
		}
	}

	return found, deletedBox
}

// maintainBoxes examines boxes and returns lists of boxes to delete and update.
// This function is extracted to be testable and avoid deadlocks by not modifying
// the store while iterating.
func maintainBoxes() (boxesToDelete []string, boxesToUpdate []api.Event) {
	boxStore.ForEach(func(box api.Box) bool {
		if box.LastUpdate.IsZero() {
			return true // continue
		}

		lastUpdate := box.LastUpdate

		if box.ExpireAfter != nil {
			if time.Since(lastUpdate) > box.ExpireAfter.Duration() {
				if logger != nil {
					logger.Info("marking expired box for deletion", zap.String("id", box.ID))
				}
				boxesToDelete = append(boxesToDelete, box.ID)
				return true // continue
			}
		}

		if box.MaxTBU != nil {
			if time.Since(lastUpdate) > box.MaxTBU.Duration() && box.Status != api.NoUpdate {
				if logger != nil {
					logger.Warn("marking box for no-update event", zap.String("id", box.ID))
				}
				var event api.Event
				event.ID = box.ID
				event.Status = api.NoUpdate
				event.Message = fmt.Sprintf("No new updates for %s.", *box.MaxTBU)
				event.Type = api.NoUpdate.String()
				boxesToUpdate = append(boxesToUpdate, event)
				return true // continue
			}
		}
		return true // continue
	})
	return boxesToDelete, boxesToUpdate
}

// Find any boxes that have expired and delete them, find any boxes which have
// not had timely updates and update their status. Also saves box file
// periodically or on exit.
func maintenanceRoutine(ctx context.Context) {
	if options.Debug {
		logger.Info("Starting box maintenance routine")
	}
	var err error
	var lastSave time.Time
	for {
		// Check which boxes need maintenance
		boxesToDelete, boxesToUpdate := maintainBoxes()

		// Now perform actions outside of the ForEach lock
		for _, id := range boxesToDelete {
			deleteBox(id, true)
		}
		for _, event := range boxesToUpdate {
			update(event)
		}
		// Write json
		if time.Since(lastSave) > 1*time.Minute {
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

		case <-time.After(1 * time.Second):
		}
	}
}

func update(event api.Event) error {
	t := time.Now()
	const maxMessages = 30

	// Update box in store (thread-safe)
	err := boxStore.Update(event.ID, func(box *api.Box) {
		box.LastMessage = event.Message

		// Prepend new message
		box.Messages = append(
			[]api.Message{
				{
					Message:   event.Message,
					Status:    event.Status.String(),
					TimeStamp: t,
				},
			},
			box.Messages...,
		)

		// Trim to max messages
		if len(box.Messages) > maxMessages {
			box.Messages = box.Messages[:maxMessages]
		}

		if event.Type != api.NoUpdate.String() {
			box.LastUpdate = t
		}

		box.Status = event.Status
		if event.MaxTBU != nil {
			if *event.MaxTBU == api.Duration(0) {
				box.ExpireAfter = nil
			} else {
				box.MaxTBU = event.MaxTBU
			}
		}

		if event.ExpireAfter != nil {
			if *event.ExpireAfter == api.Duration(0) {
				box.ExpireAfter = nil
			} else {
				box.ExpireAfter = event.ExpireAfter
			}
		}
	})

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	event.Type = "updateBox"
	dataString, err := json.Marshal(event)
	if err != nil {
		return err
	}
	events.messages <- string(dataString)
	return nil
}
