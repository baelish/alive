package server

import (
	"testing"
	"time"

	"github.com/baelish/alive/api"
)

// Helper to reset boxStore state between tests
func resetBoxStore() {
	boxStore.mu.Lock()
	boxStore.boxes = make([]api.Box, 0)
	boxStore.mu.Unlock()
}

func TestBoxStore_Add(t *testing.T) {
	// Save original state
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()

	resetBoxStore()

	tests := []struct {
		name        string
		box         api.Box
		expectError bool
	}{
		{
			name: "add box without ID (allowed)",
			box: api.Box{
				Name:   "Test Box",
				Status: api.Green,
				Size:   api.Medium,
			},
			expectError: false, // ID can be empty
		},
		{
			name: "add box with ID",
			box: api.Box{
				ID:     "test-123",
				Name:   "Test Box",
				Status: api.Green,
				Size:   api.Medium,
			},
			expectError: false,
		},
		{
			name: "add duplicate ID",
			box: api.Box{
				ID:     "test-123",
				Name:   "Duplicate",
				Status: api.Red,
				Size:   api.Small,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := boxStore.Add(tt.box)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError {
				// Verify box was added
				box, err := boxStore.GetByID(tt.box.ID)
				if err != nil {
					t.Fatalf("box not found: %v", err)
				}
				if box.Name != tt.box.Name {
					t.Errorf("expected name %s, got %s", tt.box.Name, box.Name)
				}
			}
		})
	}
}

func TestBoxStore_GetByID(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()

	resetBoxStore()

	// Add test boxes
	boxStore.Add(api.Box{ID: "box-1", Name: "Box 1", Size: api.Medium})
	boxStore.Add(api.Box{ID: "box-2", Name: "Box 2", Size: api.Large})

	tests := []struct {
		name         string
		id           string
		expectError  bool
		expectedName string
	}{
		{
			name:         "existing box",
			id:           "box-1",
			expectError:  false,
			expectedName: "Box 1",
		},
		{
			name:        "non-existing box",
			id:          "nonexistent",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			box, err := boxStore.GetByID(tt.id)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if box.Name != tt.expectedName {
					t.Errorf("expected name %s, got %s", tt.expectedName, box.Name)
				}
			}
		})
	}
}

func TestBoxStore_Delete(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()

	resetBoxStore()

	boxStore.Add(api.Box{ID: "box-1", Name: "Box 1", Size: api.Medium})
	boxStore.Add(api.Box{ID: "box-2", Name: "Box 2", Size: api.Large})

	tests := []struct {
		name         string
		id           string
		expectFound  bool
		expectedName string
	}{
		{
			name:         "delete existing box",
			id:           "box-1",
			expectFound:  true,
			expectedName: "Box 1",
		},
		{
			name:        "delete non-existing box",
			id:          "nonexistent",
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initialCount := boxStore.Len()
			found, deletedBox := boxStore.Delete(tt.id)

			if found != tt.expectFound {
				t.Errorf("expected found=%v, got %v", tt.expectFound, found)
			}

			if tt.expectFound {
				if deletedBox.Name != tt.expectedName {
					t.Errorf("expected name %s, got %s", tt.expectedName, deletedBox.Name)
				}
				if boxStore.Len() != initialCount-1 {
					t.Errorf("expected %d boxes, got %d", initialCount-1, boxStore.Len())
				}

				// Verify box is gone
				_, err := boxStore.GetByID(tt.id)
				if err == nil {
					t.Error("deleted box still exists")
				}
			} else {
				if boxStore.Len() != initialCount {
					t.Error("box count changed when deleting non-existent box")
				}
			}
		})
	}
}

func TestBoxStore_Update(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()

	resetBoxStore()

	boxStore.Add(api.Box{
		ID:     "test-box",
		Name:   "Original Name",
		Status: api.Grey,
	})

	err := boxStore.Update("test-box", func(box *api.Box) {
		box.Name = "Updated Name"
		box.Status = api.Green
		box.LastUpdate = time.Now()
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify update
	box, err := boxStore.GetByID("test-box")
	if err != nil {
		t.Fatalf("box not found: %v", err)
	}

	if box.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", box.Name)
	}
	if box.Status != api.Green {
		t.Errorf("expected status Green, got %v", box.Status)
	}
	if box.LastUpdate.IsZero() {
		t.Error("expected LastUpdate to be set")
	}

	// Test updating non-existent box
	err = boxStore.Update("nonexistent", func(box *api.Box) {
		box.Name = "Should Fail"
	})
	if err == nil {
		t.Error("expected error when updating non-existent box")
	}
}

func TestBoxStore_Sorting(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()

	resetBoxStore()

	// Add boxes in random order
	boxStore.Add(api.Box{ID: "1", Name: "Small Box", Size: api.Small})
	boxStore.Add(api.Box{ID: "2", Name: "XLarge Box", Size: api.Xlarge})
	boxStore.Add(api.Box{ID: "3", Name: "Medium Box", Size: api.Medium})
	boxStore.Add(api.Box{ID: "4", Name: "Large Box", Size: api.Large})

	boxes := boxStore.GetAll()

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

func TestBoxStore_SortingByName(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()

	resetBoxStore()

	// Add boxes with same size but different names
	boxStore.Add(api.Box{ID: "1", Name: "Zebra", Size: api.Medium})
	boxStore.Add(api.Box{ID: "2", Name: "Apple", Size: api.Medium})
	boxStore.Add(api.Box{ID: "3", Name: "Banana", Size: api.Medium})

	boxes := boxStore.GetAll()

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

func TestBoxStore_GetAll(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()

	resetBoxStore()

	boxStore.Add(api.Box{ID: "1", Name: "Box 1", Size: api.Small})
	boxStore.Add(api.Box{ID: "2", Name: "Box 2", Size: api.Medium})

	boxes := boxStore.GetAll()

	if len(boxes) != 2 {
		t.Errorf("expected 2 boxes, got %d", len(boxes))
	}

	// Verify it returns a copy (modifications don't affect store)
	boxes[0].Name = "Modified"

	boxes2 := boxStore.GetAll()
	if boxes2[0].Name == "Modified" {
		t.Error("GetAll should return a copy, not the original slice")
	}
}

func TestBoxStore_Concurrency(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()

	resetBoxStore()

	// Test concurrent adds
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			box := api.Box{
				ID:   string(rune('a' + id)),
				Name: "Box",
				Size: api.Medium,
			}
			boxStore.Add(box)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all boxes were added
	if boxStore.Len() != 10 {
		t.Errorf("expected 10 boxes, got %d", boxStore.Len())
	}
}

func TestMaintainBoxes_NoDeadlock(t *testing.T) {
	// This test verifies that the maintenance routine doesn't deadlock
	// by calling delete/update while holding the ForEach lock

	// Save original state
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()

	resetBoxStore()

	// Add multiple boxes with different expiration states
	boxes := []api.Box{
		{
			ID:          "expire-1",
			Name:        "Expire 1",
			Status:      api.Green,
			Size:        api.Small,
			LastUpdate:  time.Now().Add(-2 * time.Second),
			ExpireAfter: ptr(api.Duration(1 * time.Second)),
		},
		{
			ID:         "maxtbu-1",
			Name:       "MaxTBU 1",
			Status:     api.Green,
			Size:       api.Small,
			LastUpdate: time.Now().Add(-2 * time.Second),
			MaxTBU:     ptr(api.Duration(1 * time.Second)),
		},
		{
			ID:         "normal-1",
			Name:       "Normal 1",
			Status:     api.Green,
			Size:       api.Small,
			LastUpdate: time.Now(),
		},
	}

	for _, box := range boxes {
		if err := boxStore.Add(box); err != nil {
			t.Fatalf("failed to add box: %v", err)
		}
	}

	// Run maintenance routine logic with timeout to detect deadlock
	done := make(chan bool, 1)
	go func() {
		// Use the actual maintenance check function
		boxesToDelete, boxesToUpdate := maintainBoxes()

		// Now perform actions outside the lock - this should NOT deadlock
		for _, id := range boxesToDelete {
			boxStore.Delete(id) // This needs a write lock
		}

		for _, event := range boxesToUpdate {
			boxStore.Update(event.ID, func(box *api.Box) {
				box.Status = event.Status
			}) // This needs a write lock
		}

		done <- true
	}()

	// Wait for completion with timeout
	select {
	case <-done:
		// Success - no deadlock
	case <-time.After(2 * time.Second):
		t.Fatal("maintenance routine deadlocked - timeout waiting for completion")
	}

	// Verify expected results
	// expire-1 should be deleted
	if _, err := boxStore.GetByID("expire-1"); err == nil {
		t.Error("expire-1 should have been deleted")
	}

	// maxtbu-1 should exist with NoUpdate status
	if box, err := boxStore.GetByID("maxtbu-1"); err != nil {
		t.Error("maxtbu-1 should still exist")
	} else if box.Status != api.NoUpdate {
		t.Errorf("maxtbu-1 should have NoUpdate status, got %v", box.Status)
	}

	// normal-1 should exist unchanged
	if box, err := boxStore.GetByID("normal-1"); err != nil {
		t.Error("normal-1 should still exist")
	} else if box.Status != api.Green {
		t.Errorf("normal-1 should have Green status, got %v", box.Status)
	}
}
