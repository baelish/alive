// internal/server/boxes_test.go
package server

import (
	"testing"
	"time"

	"github.com/baelish/alive/api"
)

func TestAddBox(t *testing.T) {
	// Save and restore original boxes
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	boxes = []api.Box{}

	tests := []struct {
		name        string
		box         api.Box
		expectError bool
	}{
		{
			name: "add box without ID",
			box: api.Box{
				Name:   "Test Box",
				Status: api.Green,
				Size:   api.Medium,
			},
			expectError: false,
		},
		{
			name: "add box with custom ID",
			box: api.Box{
				ID:     "custom-id-123",
				Name:   "Custom ID Box",
				Status: api.Red,
				Size:   api.Small,
			},
			expectError: false,
		},
		{
			name: "add box with duplicate ID",
			box: api.Box{
				ID:     "custom-id-123",
				Name:   "Duplicate",
				Status: api.Amber,
				Size:   api.Large,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := addBox(tt.box)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if id == "" {
				t.Error("expected non-empty ID")
			}

			// Verify box was added
			found := false
			for _, box := range boxes {
				if box.ID == id {
					found = true
					if box.Name != tt.box.Name {
						t.Errorf("expected name %s, got %s", tt.box.Name, box.Name)
					}
					if box.LastUpdate.IsZero() {
						t.Error("expected LastUpdate to be set")
					}
					break
				}
			}

			if !found {
				t.Errorf("box with ID %s not found in boxes slice", id)
			}
		})
	}
}

func TestDeleteBox(t *testing.T) {
	// Save and restore original boxes
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	tests := []struct {
		name        string
		setupBoxes  []api.Box
		deleteID    string
		sendEvent   bool
		expectFound bool
	}{
		{
			name: "delete existing box",
			setupBoxes: []api.Box{
				{ID: "box-1", Name: "Box 1"},
				{ID: "box-2", Name: "Box 2"},
			},
			deleteID:    "box-1",
			sendEvent:   false,
			expectFound: true,
		},
		{
			name: "delete non-existing box",
			setupBoxes: []api.Box{
				{ID: "box-1", Name: "Box 1"},
			},
			deleteID:    "nonexistent",
			sendEvent:   false,
			expectFound: false,
		},
		{
			name: "delete with event",
			setupBoxes: []api.Box{
				{ID: "box-1", Name: "Box 1"},
			},
			deleteID:    "box-1",
			sendEvent:   true,
			expectFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boxes = tt.setupBoxes
			initialCount := len(boxes)

			found, deletedBox := deleteBox(tt.deleteID, tt.sendEvent)

			if found != tt.expectFound {
				t.Errorf("expected found=%v, got %v", tt.expectFound, found)
			}

			if tt.expectFound {
				if deletedBox.ID != tt.deleteID {
					t.Errorf("expected deleted box ID %s, got %s", tt.deleteID, deletedBox.ID)
				}

				if len(boxes) != initialCount-1 {
					t.Errorf("expected %d boxes after deletion, got %d", initialCount-1, len(boxes))
				}

				// Verify box is actually gone
				for _, box := range boxes {
					if box.ID == tt.deleteID {
						t.Error("deleted box still exists in boxes slice")
					}
				}
			} else {
				if len(boxes) != initialCount {
					t.Error("box count changed when deleting non-existent box")
				}
			}
		})
	}
}

func TestFindBoxByID(t *testing.T) {
	// Save and restore original boxes
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	boxes = []api.Box{
		{ID: "box-1", Name: "Box 1"},
		{ID: "box-2", Name: "Box 2"},
		{ID: "box-3", Name: "Box 3"},
	}

	tests := []struct {
		name          string
		searchID      string
		expectedIndex int
		expectError   bool
	}{
		{
			name:          "find first box",
			searchID:      "box-1",
			expectedIndex: 0,
			expectError:   false,
		},
		{
			name:          "find middle box",
			searchID:      "box-2",
			expectedIndex: 1,
			expectError:   false,
		},
		{
			name:          "find last box",
			searchID:      "box-3",
			expectedIndex: 2,
			expectError:   false,
		},
		{
			name:        "box not found",
			searchID:    "nonexistent",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, err := findBoxByID(tt.searchID)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if index != -1 {
					t.Errorf("expected index -1, got %d", index)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if index != tt.expectedIndex {
				t.Errorf("expected index %d, got %d", tt.expectedIndex, index)
			}

			if boxes[index].ID != tt.searchID {
				t.Errorf("expected box ID %s at index %d, got %s", tt.searchID, index, boxes[index].ID)
			}
		})
	}
}

func TestTestBoxID(t *testing.T) {
	// Save and restore original boxes
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	boxes = []api.Box{
		{ID: "existing-1", Name: "Box 1"},
		{ID: "existing-2", Name: "Box 2"},
	}

	tests := []struct {
		name     string
		testID   string
		expected bool
	}{
		{
			name:     "existing ID",
			testID:   "existing-1",
			expected: true,
		},
		{
			name:     "another existing ID",
			testID:   "existing-2",
			expected: true,
		},
		{
			name:     "non-existing ID",
			testID:   "nonexistent",
			expected: false,
		},
		{
			name:     "empty ID",
			testID:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testBoxID(tt.testID)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSortBoxes(t *testing.T) {
	// Save and restore original boxes
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	boxes = []api.Box{
		{ID: "1", Name: "Small Box", Size: api.Small},
		{ID: "2", Name: "Large Box", Size: api.Large},
		{ID: "3", Name: "Medium Box", Size: api.Medium},
		{ID: "4", Name: "XLarge Box", Size: api.Xlarge},
	}

	sortBoxes()

	// Verify boxes are sorted by size (largest first)
	if boxes[0].Size != api.Xlarge {
		t.Errorf("expected first box to be Xlarge, got %v", boxes[0].Size)
	}

	if boxes[len(boxes)-1].Size != api.Small {
		t.Errorf("expected last box to be Small, got %v", boxes[len(boxes)-1].Size)
	}

	// Verify descending order
	for i := 0; i < len(boxes)-1; i++ {
		if int(boxes[i].Size) < int(boxes[i+1].Size) {
			t.Errorf("boxes not sorted correctly at index %d: %v should be >= %v",
				i, boxes[i].Size, boxes[i+1].Size)
		}
	}
}

func TestSortBoxesByName(t *testing.T) {
	// Save and restore original boxes
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	// Same size boxes should be sorted by name
	boxes = []api.Box{
		{ID: "1", Name: "Zebra", Size: api.Medium},
		{ID: "2", Name: "Apple", Size: api.Medium},
		{ID: "3", Name: "Banana", Size: api.Medium},
	}

	sortBoxes()

	// Verify alphabetical order for same-size boxes
	if boxes[0].Name != "Apple" {
		t.Errorf("expected first box name 'Apple', got '%s'", boxes[0].Name)
	}

	if boxes[1].Name != "Banana" {
		t.Errorf("expected second box name 'Banana', got '%s'", boxes[1].Name)
	}

	if boxes[2].Name != "Zebra" {
		t.Errorf("expected third box name 'Zebra', got '%s'", boxes[2].Name)
	}
}

func TestUpdate(t *testing.T) {
	// Save and restore original boxes
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	// Setup a box to update
	boxes = []api.Box{
		{
			ID:          "test-box",
			Name:        "Test",
			Status:      api.Grey,
			LastMessage: "old message",
			Messages:    []api.Message{},
		},
	}

	event := api.Event{
		ID:      "test-box",
		Status:  api.Green,
		Message: "new message",
		Type:    "test",
	}

	update(event)

	// Find the updated box
	idx, err := findBoxByID("test-box")
	if err != nil {
		t.Fatalf("box not found: %v", err)
	}

	updatedBox := boxes[idx]

	// Verify updates
	if updatedBox.Status != api.Green {
		t.Errorf("expected status Green, got %v", updatedBox.Status)
	}

	if updatedBox.LastMessage != "new message" {
		t.Errorf("expected LastMessage 'new message', got '%s'", updatedBox.LastMessage)
	}

	if len(updatedBox.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(updatedBox.Messages))
	}

	if updatedBox.Messages[0].Message != "new message" {
		t.Errorf("expected message 'new message', got '%s'", updatedBox.Messages[0].Message)
	}

	if updatedBox.LastUpdate.IsZero() {
		t.Error("expected LastUpdate to be set")
	}
}

func TestUpdateWithDurations(t *testing.T) {
	// Save and restore original boxes
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	boxes = []api.Box{
		{
			ID:     "test-box",
			Name:   "Test",
			Status: api.Grey,
		},
	}

	event := api.Event{
		ID:          "test-box",
		Status:      api.Green,
		Message:     "test",
		Type:        "test",
		ExpireAfter: api.Duration{Duration: 5 * time.Minute, Set: true},
		MaxTBU:      api.Duration{Duration: 10 * time.Minute, Set: true},
	}

	update(event)

	idx, _ := findBoxByID("test-box")
	updatedBox := boxes[idx]

	if updatedBox.ExpireAfter.Duration != 5*time.Minute {
		t.Errorf("expected ExpireAfter 5m, got %v", updatedBox.ExpireAfter.Duration)
	}

	if updatedBox.MaxTBU.Duration != 10*time.Minute {
		t.Errorf("expected MaxTBU 10m, got %v", updatedBox.MaxTBU.Duration)
	}
}

func TestUpdateMessageLimit(t *testing.T) {
	// Save and restore original boxes
	originalBoxes := boxes
	defer func() { boxes = originalBoxes }()

	// Create a box with 30 existing messages
	existingMessages := make([]api.Message, 30)
	for i := 0; i < 30; i++ {
		existingMessages[i] = api.Message{
			Message: "old message",
			Status:  "old",
		}
	}

	boxes = []api.Box{
		{
			ID:       "test-box",
			Name:     "Test",
			Status:   api.Grey,
			Messages: existingMessages,
		},
	}

	// Add a new message
	event := api.Event{
		ID:      "test-box",
		Status:  api.Green,
		Message: "new message",
		Type:    "test",
	}

	update(event)

	idx, _ := findBoxByID("test-box")
	updatedBox := boxes[idx]

	// Should still have maximum of 30 messages
	if len(updatedBox.Messages) != 30 {
		t.Errorf("expected 30 messages (limit), got %d", len(updatedBox.Messages))
	}

	// First message should be the new one
	if updatedBox.Messages[0].Message != "new message" {
		t.Errorf("expected first message to be 'new message', got '%s'", updatedBox.Messages[0].Message)
	}
}
