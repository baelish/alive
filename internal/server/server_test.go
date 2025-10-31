package server

import (
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	// Test that the time format constant is correct
	now := time.Now()
	formatted := now.Format(timeFormat)

	// Parse it back to ensure it's a valid format
	parsed, err := time.Parse(timeFormat, formatted)
	if err != nil {
		t.Errorf("timeFormat is invalid: %v", err)
	}

	// Check that parsing preserves the time (within a millisecond)
	if parsed.Sub(now).Abs() > time.Millisecond {
		t.Errorf("timeFormat loses precision: original=%v, parsed=%v", now, parsed)
	}
}
