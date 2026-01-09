package api

import (
	"encoding/json"
	"testing"
	"time"
)

// TestStatusMarshaling tests Status JSON marshaling and unmarshaling
func TestStatusMarshaling(t *testing.T) {
	tests := []struct {
		name       string
		status     Status
		jsonStr    string
		shouldFail bool
	}{
		{"green", Green, `"green"`, false},
		{"red", Red, `"red"`, false},
		{"amber", Amber, `"amber"`, false},
		{"grey", Grey, `"grey"`, false},
		{"gray (alternative)", Grey, `"gray"`, false},
		{"noUpdate", NoUpdate, `"noUpdate"`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name+" marshal", func(t *testing.T) {
			data, err := json.Marshal(tt.status)
			if err != nil {
				t.Fatalf("failed to marshal status: %v", err)
			}

			if string(data) != tt.jsonStr && tt.name != "gray (alternative)" {
				t.Errorf("expected %s, got %s", tt.jsonStr, string(data))
			}
		})

		t.Run(tt.name+" unmarshal", func(t *testing.T) {
			var s Status
			err := json.Unmarshal([]byte(tt.jsonStr), &s)
			if err != nil {
				t.Fatalf("failed to unmarshal status: %v", err)
			}

			if s != tt.status {
				t.Errorf("expected status %v, got %v", tt.status, s)
			}
		})
	}
}

func TestStatusUnmarshalInvalid(t *testing.T) {
	var s Status
	err := json.Unmarshal([]byte(`"invalid"`), &s)
	if err == nil {
		t.Error("expected error for invalid status, got nil")
	}
}

func TestStatusString(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{Grey, "grey"},
		{Red, "red"},
		{Amber, "amber"},
		{Green, "green"},
		{NoUpdate, "noUpdate"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.status.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.status.String())
			}
		})
	}
}

// TestBoxSizeMarshaling tests BoxSize JSON marshaling and unmarshaling
func TestBoxSizeMarshaling(t *testing.T) {
	tests := []struct {
		name    string
		size    BoxSize
		jsonStr string
	}{
		{"dot", Dot, `"dot"`},
		{"micro", Micro, `"micro"`},
		{"dmicro", Dmicro, `"dmicro"`},
		{"small", Small, `"small"`},
		{"dsmall", Dsmall, `"dsmall"`},
		{"medium", Medium, `"medium"`},
		{"dmedium", Dmedium, `"dmedium"`},
		{"large", Large, `"large"`},
		{"dlarge", Dlarge, `"dlarge"`},
		{"xlarge", Xlarge, `"xlarge"`},
	}

	for _, tt := range tests {
		t.Run(tt.name+" marshal", func(t *testing.T) {
			data, err := json.Marshal(tt.size)
			if err != nil {
				t.Fatalf("failed to marshal box size: %v", err)
			}

			if string(data) != tt.jsonStr {
				t.Errorf("expected %s, got %s", tt.jsonStr, string(data))
			}
		})

		t.Run(tt.name+" unmarshal", func(t *testing.T) {
			var bs BoxSize
			err := json.Unmarshal([]byte(tt.jsonStr), &bs)
			if err != nil {
				t.Fatalf("failed to unmarshal box size: %v", err)
			}

			if bs != tt.size {
				t.Errorf("expected box size %v, got %v", tt.size, bs)
			}
		})
	}
}

func TestBoxSizeUnmarshalInvalid(t *testing.T) {
	var bs BoxSize
	err := json.Unmarshal([]byte(`"invalid"`), &bs)
	if err == nil {
		t.Error("expected error for invalid box size, got nil")
	}
}

func TestBoxSizeString(t *testing.T) {
	tests := []struct {
		size     BoxSize
		expected string
	}{
		{Dot, "dot"},
		{Micro, "micro"},
		{Small, "small"},
		{Medium, "medium"},
		{Large, "large"},
		{Xlarge, "xlarge"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.size.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.size.String())
			}
		})
	}
}

// TestBoxMarshaling tests complete Box marshaling
func TestBoxMarshaling(t *testing.T) {
	info := map[string]string{"key": "value"}
	expireTime := 15 * time.Minute
	maxTBU := 5 * time.Minute
	box := Box{
		ID:          "test-123",
		Name:        "Test Box",
		Description: "A test box",
		DisplayName: "Display",
		Status:      Green,
		Size:        Medium,
		Info:        &info,
		Parent:      "parent-id",
		Links: []Links{
			{Name: "GitHub", URL: "https://github.com"},
		},
		Messages: []Message{
			{Message: "Hello", Status: "ok", TimeStamp: time.Now()},
		},
		LastMessage: "Hello",
		ExpireAfter: &expireTime,
		MaxTBU:      &maxTBU,
	}

	data, err := json.Marshal(box)
	if err != nil {
		t.Fatalf("failed to marshal box: %v", err)
	}

	var unmarshaled Box
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal box: %v", err)
	}

	if unmarshaled.ID != box.ID {
		t.Errorf("expected ID %s, got %s", box.ID, unmarshaled.ID)
	}

	if unmarshaled.Name != box.Name {
		t.Errorf("expected Name %s, got %s", box.Name, unmarshaled.Name)
	}

	if unmarshaled.Status != box.Status {
		t.Errorf("expected Status %v, got %v", box.Status, unmarshaled.Status)
	}

	if unmarshaled.Size != box.Size {
		t.Errorf("expected Size %v, got %v", box.Size, unmarshaled.Size)
	}
}

// TestEventMarshaling tests Event marshaling
func TestEventMarshaling(t *testing.T) {
	box := Box{ID: "test", Name: "Test"}
	expireTime := 15 * time.Minute
	maxTBU := 5 * time.Minute
	event := Event{
		ID:          "event-123",
		After:       "after-id",
		Box:         &box,
		Status:      Green,
		Message:     "Test message",
		Type:        "updateBox",
		ExpireAfter: &expireTime,
		MaxTBU:      &maxTBU,
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}

	var unmarshaled Event
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal event: %v", err)
	}

	if unmarshaled.ID != event.ID {
		t.Errorf("expected ID %s, got %s", event.ID, unmarshaled.ID)
	}

	if unmarshaled.Type != event.Type {
		t.Errorf("expected Type %s, got %s", event.Type, unmarshaled.Type)
	}

	if unmarshaled.Status != event.Status {
		t.Errorf("expected Status %v, got %v", event.Status, unmarshaled.Status)
	}
}

// TestErrorResponse tests ErrorResponse struct
func TestErrorResponse(t *testing.T) {
	errResp := ErrorResponse{
		Message: "Something went wrong",
		Error:   "detailed error",
	}

	data, err := json.Marshal(errResp)
	if err != nil {
		t.Fatalf("failed to marshal error response: %v", err)
	}

	var unmarshaled ErrorResponse
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal error response: %v", err)
	}

	if unmarshaled.Message != errResp.Message {
		t.Errorf("expected Message %s, got %s", errResp.Message, unmarshaled.Message)
	}

	if unmarshaled.Error != errResp.Error {
		t.Errorf("expected Error %s, got %s", errResp.Error, unmarshaled.Error)
	}
}
