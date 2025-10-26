package server

import (
	"testing"

	"github.com/baelish/alive/api"
)

// TestLogStructDetails tests the logger helper function
func TestLogStructDetails(t *testing.T) {
	box := api.Box{
		ID:     "test-123",
		Name:   "Test Box",
		Status: api.Green,
		Size:   api.Medium,
	}

	fields := logStructDetails(box)

	// Should have multiple fields from the Box struct
	if len(fields) == 0 {
		t.Error("expected fields to be extracted")
	}

	// Verify we can find specific fields
	foundID := false
	foundName := false
	for _, field := range fields {
		if field.Key == "id" {
			foundID = true
		}
		if field.Key == "name" {
			foundName = true
		}
	}

	if !foundID {
		t.Error("expected to find 'id' field")
	}

	if !foundName {
		t.Error("expected to find 'name' field")
	}
}

func TestLogStructDetailsWithPointer(t *testing.T) {
	type TestStruct struct {
		Value string `json:"value"`
	}

	testObj := &TestStruct{Value: "test"}
	fields := logStructDetails(testObj)

	if len(fields) != 1 {
		t.Errorf("expected 1 field, got %d", len(fields))
	}

	if fields[0].Key != "value" {
		t.Errorf("expected field key 'value', got '%s'", fields[0].Key)
	}
}

func TestIndexComma(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"name,omitempty", 4},
		{"value", -1},
		{"first,second,third", 5},
		{"", -1},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := indexComma(tt.input)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
