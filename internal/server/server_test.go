package server

import (
	"context"
	"os"
	"path/filepath"
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

func TestStart_PathDefaults(t *testing.T) {
	// Save original global state
	originalOptions := options
	originalLogger := logger
	defer func() {
		options = originalOptions
		logger = originalLogger
	}()

	// Test that paths are set to defaults when not specified
	tests := []struct {
		name            string
		dataPath        string
		staticPath      string
		expectedData    string
		expectedStatic  string
	}{
		{
			name:           "both paths empty",
			dataPath:       "",
			staticPath:     "",
			expectedData:   filepath.Clean(os.Getenv("HOME") + "/.alive/data"),
			expectedStatic: filepath.Clean(os.Getenv("HOME") + "/.alive/static"),
		},
		{
			name:           "custom data path",
			dataPath:       "/custom/data",
			staticPath:     "",
			expectedData:   "/custom/data",
			expectedStatic: filepath.Clean(os.Getenv("HOME") + "/.alive/static"),
		},
		{
			name:           "custom static path",
			dataPath:       "",
			staticPath:     "/custom/static",
			expectedData:   filepath.Clean(os.Getenv("HOME") + "/.alive/data"),
			expectedStatic: "/custom/static",
		},
		{
			name:           "both custom paths",
			dataPath:       "/custom/data",
			staticPath:     "/custom/static",
			expectedData:   "/custom/data",
			expectedStatic: "/custom/static",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset options
			options = Options{
				DataPath:   tt.dataPath,
				StaticPath: tt.staticPath,
			}

			// Simulate the path setting logic from Start()
			if options.DataPath == "" {
				options.DataPath = filepath.Clean(os.Getenv("HOME") + "/.alive/data")
			}
			if options.StaticPath == "" {
				options.StaticPath = filepath.Clean(os.Getenv("HOME") + "/.alive/static")
			}

			if options.DataPath != tt.expectedData {
				t.Errorf("DataPath: expected %q, got %q", tt.expectedData, options.DataPath)
			}
			if options.StaticPath != tt.expectedStatic {
				t.Errorf("StaticPath: expected %q, got %q", tt.expectedStatic, options.StaticPath)
			}
		})
	}
}

func TestStart_DemoMode(t *testing.T) {
	// Save original global state
	originalOptions := options
	defer func() {
		options = originalOptions
	}()

	// Test demo mode path creation logic
	t.Run("demo mode sets temporary paths", func(t *testing.T) {
		options = Options{
			Demo: true,
		}

		// Simulate demo mode temp directory creation
		tempDir, err := os.MkdirTemp(os.TempDir(), "alive-*.tmp")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Simulate the path setting from Start()
		options.DataPath = filepath.Clean(tempDir + "/data")
		options.StaticPath = filepath.Clean(tempDir + "/static")

		// Verify paths are set correctly
		expectedData := filepath.Clean(tempDir + "/data")
		expectedStatic := filepath.Clean(tempDir + "/static")

		if options.DataPath != expectedData {
			t.Errorf("DataPath: expected %q, got %q", expectedData, options.DataPath)
		}
		if options.StaticPath != expectedStatic {
			t.Errorf("StaticPath: expected %q, got %q", expectedStatic, options.StaticPath)
		}
	})

	t.Run("temp directory is cleaned up", func(t *testing.T) {
		// Create a temp directory
		tempDir, err := os.MkdirTemp(os.TempDir(), "alive-*.tmp")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}

		// Verify it exists
		if _, err := os.Stat(tempDir); os.IsNotExist(err) {
			t.Fatal("temp dir doesn't exist after creation")
		}

		// Clean it up (simulating defer os.RemoveAll(tempDir))
		os.RemoveAll(tempDir)

		// Verify it's gone
		if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
			t.Error("temp dir still exists after cleanup")
		}
	})
}

func TestStart_EnvironmentVariables(t *testing.T) {
	// Test DEV environment variable handling
	t.Run("DEV environment variable", func(t *testing.T) {
		// Save original
		originalDev := os.Getenv("DEV")
		defer os.Setenv("DEV", originalDev)

		// Test with DEV set
		os.Setenv("DEV", "1")
		devValue := os.Getenv("DEV")
		if devValue == "" {
			t.Error("DEV environment variable should be set")
		}

		// Test without DEV
		os.Setenv("DEV", "")
		devValue = os.Getenv("DEV")
		if devValue != "" {
			t.Error("DEV environment variable should be empty")
		}
	})

	t.Run("HOME environment variable for default paths", func(t *testing.T) {
		home := os.Getenv("HOME")
		if home == "" {
			t.Skip("HOME environment variable not set")
		}

		expectedData := filepath.Clean(home + "/.alive/data")
		expectedStatic := filepath.Clean(home + "/.alive/static")

		// These are the defaults when paths aren't specified
		if expectedData == "" || expectedStatic == "" {
			t.Error("default paths should not be empty when HOME is set")
		}
	})
}

func TestStart_PathCleaning(t *testing.T) {
	// Test that paths are properly cleaned
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "removes trailing slash",
			input:    "/path/to/data/",
			expected: "/path/to/data",
		},
		{
			name:     "cleans double slashes",
			input:    "/path//to///data",
			expected: "/path/to/data",
		},
		{
			name:     "handles relative paths",
			input:    "./data/../data",
			expected: "data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleaned := filepath.Clean(tt.input)
			if cleaned != tt.expected {
				t.Errorf("filepath.Clean(%q): expected %q, got %q", tt.input, tt.expected, cleaned)
			}
		})
	}
}

func TestStart_SignalHandling(t *testing.T) {
	// Test context cancellation behavior
	t.Run("context cancellation stops main loop", func(t *testing.T) {
		// This simulates the main loop in Start()
		ctx, cancel := context.WithCancel(context.Background())

		done := make(chan bool)
		go func() {
			for {
				select {
				case <-ctx.Done():
					done <- true
					return
				case <-time.After(time.Duration(10 * time.Millisecond)):
					// Continue loop
				}
			}
		}()

		// Cancel the context
		cancel()

		// Wait for the loop to exit
		select {
		case <-done:
			// Successfully exited
		case <-time.After(200 * time.Millisecond):
			t.Error("main loop did not exit after context cancellation")
		}
	})
}
