// internal/server/df-functions_test.go
package server

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/baelish/alive/api"
)

func TestSaveAndLoadBoxFile(t *testing.T) {
	// Save and restore original state
	originalBoxes := boxes
	originalBoxFile := boxFile
	defer func() {
		boxes = originalBoxes
		boxFile = originalBoxFile
	}()

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	boxFile = filepath.Join(tempDir, "boxes.json")

	// Setup test boxes
	boxes = []api.Box{
		{ID: "test-1", Name: "Test Box 1", Status: api.Green, Size: api.Medium},
		{ID: "test-2", Name: "Test Box 2", Status: api.Red, Size: api.Small},
	}

	// Test saving
	err := saveBoxFile()
	if err != nil {
		t.Fatalf("saveBoxFile failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(boxFile); os.IsNotExist(err) {
		t.Fatal("box file was not created")
	}

	// Read and verify contents
	data, err := os.ReadFile(boxFile)
	if err != nil {
		t.Fatalf("failed to read box file: %v", err)
	}

	var loadedBoxes []api.Box
	err = json.Unmarshal(data, &loadedBoxes)
	if err != nil {
		t.Fatalf("failed to unmarshal box data: %v", err)
	}

	if len(loadedBoxes) != 2 {
		t.Errorf("expected 2 boxes, got %d", len(loadedBoxes))
	}

	if loadedBoxes[0].ID != "test-1" {
		t.Errorf("expected first box ID 'test-1', got '%s'", loadedBoxes[0].ID)
	}
}

func TestSaveBoxFileBackups(t *testing.T) {
	// Save and restore original state
	originalBoxes := boxes
	originalBoxFile := boxFile
	defer func() {
		boxes = originalBoxes
		boxFile = originalBoxFile
	}()

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	boxFile = filepath.Join(tempDir, "boxes.json")

	boxes = []api.Box{
		{ID: "test-1", Name: "Test Box 1"},
	}

	// Save first time
	err := saveBoxFile()
	if err != nil {
		t.Fatalf("first save failed: %v", err)
	}

	// Modify and save again
	boxes = []api.Box{
		{ID: "test-1", Name: "Modified Box 1"},
		{ID: "test-2", Name: "New Box 2"},
	}

	err = saveBoxFile()
	if err != nil {
		t.Fatalf("second save failed: %v", err)
	}

	// Verify backup was created
	backup := boxFile + ".bak1"
	if _, err := os.Stat(backup); os.IsNotExist(err) {
		t.Error("backup file .bak1 was not created")
	}

	// Read backup and verify it has the old data
	data, err := os.ReadFile(backup)
	if err != nil {
		t.Fatalf("failed to read backup file: %v", err)
	}

	var backupBoxes []api.Box
	err = json.Unmarshal(data, &backupBoxes)
	if err != nil {
		t.Fatalf("failed to unmarshal backup data: %v", err)
	}

	if len(backupBoxes) != 1 {
		t.Errorf("expected 1 box in backup, got %d", len(backupBoxes))
	}

	if backupBoxes[0].Name != "Test Box 1" {
		t.Errorf("expected backup to contain original name 'Test Box 1', got '%s'", backupBoxes[0].Name)
	}
}

func TestCreateDataFiles(t *testing.T) {
	// Save and restore original state
	originalOptions := options
	originalBoxFile := boxFile
	defer func() {
		options = originalOptions
		boxFile = originalBoxFile
	}()

	// Create a temporary directory
	tempDir := t.TempDir()
	options.DataPath = tempDir

	// Test creating data files
	createDataFiles()

	// Verify directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("data directory was not created")
	}

	// Verify box file exists
	expectedBoxFile := filepath.Join(tempDir, "boxes.json")
	if _, err := os.Stat(expectedBoxFile); os.IsNotExist(err) {
		t.Error("boxes.json file was not created")
	}

	// Verify it's valid JSON
	data, err := os.ReadFile(expectedBoxFile)
	if err != nil {
		t.Fatalf("failed to read created box file: %v", err)
	}

	var testBoxes []api.Box
	err = json.Unmarshal(data, &testBoxes)
	if err != nil {
		t.Fatalf("created box file is not valid JSON: %v", err)
	}
}

func TestGetBoxesFromDataFile(t *testing.T) {
	// Save and restore original state
	originalBoxes := boxes
	originalBoxFile := boxFile
	defer func() {
		boxes = originalBoxes
		boxFile = originalBoxFile
	}()

	// Create a temporary directory and file
	tempDir := t.TempDir()
	boxFile = filepath.Join(tempDir, "boxes.json")

	// Create test data
	testBoxes := []api.Box{
		{ID: "small-1", Name: "Small Box", Status: api.Green, Size: api.Small},
		{ID: "large-1", Name: "Large Box", Status: api.Red, Size: api.Large},
		{ID: "medium-1", Name: "Medium Box", Status: api.Amber, Size: api.Medium},
	}

	// Write test data to file
	data, err := json.Marshal(testBoxes)
	if err != nil {
		t.Fatalf("failed to marshal test data: %v", err)
	}

	err = os.WriteFile(boxFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Clear boxes and load from file
	boxes = []api.Box{}
	getBoxesFromDataFile()

	// Verify boxes were loaded
	if len(boxes) != 3 {
		t.Errorf("expected 3 boxes, got %d", len(boxes))
	}

	// Verify sorting (largest first)
	if boxes[0].Size != api.Large {
		t.Errorf("expected first box to be Large, got %v", boxes[0].Size)
	}

	if boxes[len(boxes)-1].Size != api.Small {
		t.Errorf("expected last box to be Small, got %v", boxes[len(boxes)-1].Size)
	}
}
