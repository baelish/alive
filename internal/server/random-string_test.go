package server

import (
	"testing"
)

func TestRandStringBytes(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"length 5", 5},
		{"length 10", 10},
		{"length 20", 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := randStringBytes(tt.length)

			if len(result) != tt.length {
				t.Errorf("expected length %d, got %d", tt.length, len(result))
			}

			// Verify all characters are from the allowed set
			for _, char := range result {
				found := false
				for _, allowed := range randomBytes {
					if char == allowed {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("character '%c' not in allowed character set", char)
				}
			}
		})
	}

	// Test uniqueness (probabilistic)
	results := make(map[string]bool)
	for i := 0; i < 100; i++ {
		str := randStringBytes(10)
		if results[str] {
			t.Error("generated duplicate random string")
		}
		results[str] = true
	}
}
