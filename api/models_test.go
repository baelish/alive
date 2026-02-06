package api

import (
	"encoding/json"
	"testing"
	"time"
)

func ptr[T any](v T) *T { return &v }

func TestDurationMethods(t *testing.T) {
	tests := []struct {
		name    string
		input   *Duration
		wantStr string
		wantDur time.Duration
	}{
		{
			name:    "nil pointer",
			input:   nil,
			wantStr: "<nil>", // String() is on value, so nil pointer can't call it; skip
			wantDur: 0,
		},
		{
			name:    "5 minutes",
			input:   func() *Duration { d := Duration(5 * time.Minute); return &d }(),
			wantStr: "5m0s",
			wantDur: 5 * time.Minute,
		},
		{
			name:    "zero duration",
			input:   func() *Duration { d := Duration(0); return &d }(),
			wantStr: "0s",
			wantDur: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input != nil {
				// Test String() on value
				gotStr := (*tt.input).String()
				if gotStr != tt.wantStr {
					t.Errorf("String() = %q; want %q", gotStr, tt.wantStr)
				}
			}

			// Test Duration() on pointer
			gotDur := tt.input.Duration()
			if gotDur != tt.wantDur {
				t.Errorf("Duration() = %v; want %v", gotDur, tt.wantDur)
			}
		})
	}
}

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

// TestDurationMarshaling tests Duration custom marshaling
func TestDurationMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		duration *Duration
		expected string
	}{
		{
			name:     "set duration",
			duration: ptr(Duration(5 * time.Minute)),
			expected: `"5m0s"`,
		},
		{
			name:     "unset duration",
			duration: nil,
			expected: `null`,
		},
		{
			name:     "zero but set",
			duration: ptr(Duration(0)),
			expected: `"0s"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.duration)
			if err != nil {
				t.Fatalf("failed to marshal duration: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(data))
			}
		})
	}
}

// TestDurationUnmarshaling tests Duration custom unmarshaling
func TestDurationUnmarshaling(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectedDur time.Duration
		expectedSet bool
		shouldFail  bool
	}{
		{
			name:        "null",
			input:       []byte(`null`),
			expectedDur: 0,
			expectedSet: false,
			shouldFail:  false,
		},
		{
			name:        "empty string",
			input:       []byte(`""`),
			expectedDur: 0,
			expectedSet: true,
			shouldFail:  false,
		},
		{
			name:        "duration string",
			input:       []byte(`"5m"`),
			expectedDur: 5 * time.Minute,
			expectedSet: true,
			shouldFail:  false,
		},
		{
			name:        "seconds as number",
			input:       []byte(`60`),
			expectedDur: 60 * time.Second,
			expectedSet: true,
			shouldFail:  false,
		},
		{
			name:        "seconds as string",
			input:       []byte(`"120"`),
			expectedDur: 120 * time.Second,
			expectedSet: true,
			shouldFail:  false,
		},
		{
			name:        "nano seconds as number",
			input:       []byte(`60000000000`),
			expectedDur: 1 * time.Minute,
			expectedSet: true,
			shouldFail:  false,
		},
		{
			name:        "complex duration",
			input:       []byte(`"1h30m"`),
			expectedDur: 90 * time.Minute,
			expectedSet: true,
			shouldFail:  false,
		},
		{
			name:        "nonsense",
			input:       []byte(`"nonsense"`),
			expectedDur: 0,
			expectedSet: false,
			shouldFail:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d *Duration
			err := json.Unmarshal(tt.input, &d)

			if tt.shouldFail {
				if err == nil {
					t.Errorf("%s: expected error, got nil", tt.name)
				}
				return
			}

			if err != nil {
				t.Fatalf("%s: unexpected error: %v", tt.name, err)
			}

			if d != nil {
				if !tt.expectedSet {
					t.Fatalf("%s: expected nil, got %s", tt.name, *d)
				}
				if *d != Duration(tt.expectedDur) {
					t.Errorf("%s: expected duration %s, got %s", tt.name, tt.expectedDur, *d)
				}
			} else {
				if tt.expectedSet {
					t.Errorf("%s: expected duration %s, but var was not set", tt.name, tt.expectedDur)
				}
			}
		})
	}
}

// TestBoxMarshaling tests complete Box marshaling
func TestBoxMarshaling(t *testing.T) {
	info := map[string]string{"key": "value"}
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
		ExpireAfter: ptr(Duration(5 * time.Minute)),
		MaxTBU:      ptr(Duration(10 * time.Minute)),
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
		t.Errorf("expected Status %s, got %s", box.Status, unmarshaled.Status)
	}

	if unmarshaled.Size != box.Size {
		t.Errorf("expected Size %s, got %s", box.Size, unmarshaled.Size)
	}

	if *unmarshaled.MaxTBU != *box.MaxTBU {
		t.Errorf("expected MaxTBU %s, got %s", *box.MaxTBU, *unmarshaled.MaxTBU)
	}

	if *unmarshaled.ExpireAfter != *box.ExpireAfter {
		t.Errorf("expected ExpireAfter %s, got %s", *box.ExpireAfter, *unmarshaled.ExpireAfter)
	}

	box.MaxTBU = nil
	box.ExpireAfter = nil
	data, err = json.Marshal(box)
	if err != nil {
		t.Fatalf("failed to marshal box: %v", err)
	}

	unmarshaled = Box{}
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal box: %v", err)
	}

	if unmarshaled.MaxTBU != nil {
		t.Errorf("expected MaxTBU to be nil, got %s", *unmarshaled.MaxTBU)
	}

	if unmarshaled.ExpireAfter != nil {
		t.Errorf("expected ExpireAfter to be nil, got %s", *unmarshaled.ExpireAfter)
	}

}

// TestEventMarshaling tests Event marshaling
func TestEventMarshaling(t *testing.T) {
	box := Box{ID: "test", Name: "Test"}
	event := Event{
		ID:          "event-123",
		After:       "after-id",
		Box:         &box,
		Status:      Green,
		Message:     "Test message",
		Type:        "updateBox",
		ExpireAfter: ptr(Duration(5 * time.Minute)),
		MaxTBU:      ptr(Duration(10 * time.Minute)),
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

	if *unmarshaled.MaxTBU != *event.MaxTBU {
		t.Errorf("expected MaxTBU %s, got %s", *event.MaxTBU, *unmarshaled.MaxTBU)
	}

	if *unmarshaled.ExpireAfter != *event.ExpireAfter {
		t.Errorf("expected ExpireAfter %s, got %s", *event.ExpireAfter, *unmarshaled.ExpireAfter)
	}

	event.MaxTBU = nil
	event.ExpireAfter = nil
	data, err = json.Marshal(event)
	if err != nil {
		t.Fatalf("failed to marshal event: %v", err)
	}

	unmarshaled = Event{}
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal event: %v", err)
	}

	if unmarshaled.MaxTBU != nil {
		t.Errorf("expected MaxTBU to be nil, got %s", *unmarshaled.MaxTBU)
	}

	if unmarshaled.ExpireAfter != nil {
		t.Errorf("expected ExpireAfter to be nil, got %s", *unmarshaled.ExpireAfter)
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

// TestBoxSanitisation tests that a zero value is changed into nil and any other
// valid value is unchanged.
func TestBoxSanitisation(t *testing.T) {
	tests := []struct {
		name    string
		input   Box
		wantMax *Duration
		wantExp *Duration
	}{
		{
			name: "Convert zero durations to nil",
			input: Box{
				MaxTBU:      ptr(Duration(0)),
				ExpireAfter: ptr(Duration(0)),
			},
			wantMax: nil,
			wantExp: nil,
		},
		{
			name: "Keep positive durations",
			input: Box{
				MaxTBU:      ptr(Duration(10 * time.Second)),
				ExpireAfter: ptr(Duration(1 * time.Hour)),
			},
			wantMax: ptr(Duration(10 * time.Second)),
			wantExp: ptr(Duration(1 * time.Hour)),
		},
		{
			name: "Handle mixed nil and zero",
			input: Box{
				MaxTBU:      nil,
				ExpireAfter: ptr(Duration(0)),
			},
			wantMax: nil,
			wantExp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input.Sanitise()

			if !compareDurationPtr(tt.input.MaxTBU, tt.wantMax) {
				t.Errorf("%s: MaxTBU = %v, want %v", tt.name, tt.input.MaxTBU, tt.wantMax)
			}

			if !compareDurationPtr(tt.input.ExpireAfter, tt.wantExp) {
				t.Errorf("%s: ExpireAfter = %v, want %v", tt.name, tt.input.ExpireAfter, tt.wantExp)
			}
		})
	}
}

// Check Duration pointers are as expected
func compareDurationPtr(got, want *Duration) bool {
	if got == nil && want == nil {
		return true
	}
	if got != nil && want != nil {
		return *got == *want
	}
	return false
}
