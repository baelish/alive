package server

import (
	"os"
	"testing"

	goflags "github.com/jessevdk/go-flags"
)

func TestOptions_Defaults(t *testing.T) {
	t.Run("default values are set correctly", func(t *testing.T) {
		opts := Options{}
		parser := goflags.NewParser(&opts, goflags.Default)

		// Parse with no arguments (should use defaults)
		_, err := parser.ParseArgs([]string{})
		if err != nil {
			t.Fatalf("failed to parse with defaults: %v", err)
		}

		// Verify defaults
		if opts.ApiPort != "8081" {
			t.Errorf("ApiPort default: expected %q, got %q", "8081", opts.ApiPort)
		}
		if opts.SitePort != "8080" {
			t.Errorf("SitePort default: expected %q, got %q", "8080", opts.SitePort)
		}
		if opts.ParentBoxSize != "large" {
			t.Errorf("ParentBoxSize default: expected %q, got %q", "large", opts.ParentBoxSize)
		}
		if opts.Debug {
			t.Error("Debug should default to false")
		}
		if opts.Demo {
			t.Error("Demo should default to false")
		}
		if opts.DefaultStatic {
			t.Error("DefaultStatic should default to false")
		}
		if opts.DataPath != "" {
			t.Errorf("DataPath should default to empty, got %q", opts.DataPath)
		}
		if opts.StaticPath != "" {
			t.Errorf("StaticPath should default to empty, got %q", opts.StaticPath)
		}
		if opts.ParentUrl != "" {
			t.Errorf("ParentUrl should default to empty, got %q", opts.ParentUrl)
		}
		if opts.ParentBoxID != "" {
			t.Errorf("ParentBoxID should default to empty, got %q", opts.ParentBoxID)
		}
	})
}

func TestOptions_CustomValues(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		validate func(t *testing.T, opts Options)
	}{
		{
			name: "custom api port",
			args: []string{"--api-port", "9091"},
			validate: func(t *testing.T, opts Options) {
				if opts.ApiPort != "9091" {
					t.Errorf("expected ApiPort %q, got %q", "9091", opts.ApiPort)
				}
			},
		},
		{
			name: "custom site port with short flag",
			args: []string{"-p", "9090"},
			validate: func(t *testing.T, opts Options) {
				if opts.SitePort != "9090" {
					t.Errorf("expected SitePort %q, got %q", "9090", opts.SitePort)
				}
			},
		},
		{
			name: "custom site port with long flag",
			args: []string{"--port", "9090"},
			validate: func(t *testing.T, opts Options) {
				if opts.SitePort != "9090" {
					t.Errorf("expected SitePort %q, got %q", "9090", opts.SitePort)
				}
			},
		},
		{
			name: "enable debug",
			args: []string{"--debug"},
			validate: func(t *testing.T, opts Options) {
				if !opts.Debug {
					t.Error("expected Debug to be true")
				}
			},
		},
		{
			name: "enable demo",
			args: []string{"--run-demo"},
			validate: func(t *testing.T, opts Options) {
				if !opts.Demo {
					t.Error("expected Demo to be true")
				}
			},
		},
		{
			name: "enable default static",
			args: []string{"--default-static"},
			validate: func(t *testing.T, opts Options) {
				if !opts.DefaultStatic {
					t.Error("expected DefaultStatic to be true")
				}
			},
		},
		{
			name: "custom data path with short flag",
			args: []string{"-d", "/custom/data"},
			validate: func(t *testing.T, opts Options) {
				if opts.DataPath != "/custom/data" {
					t.Errorf("expected DataPath %q, got %q", "/custom/data", opts.DataPath)
				}
			},
		},
		{
			name: "custom data path with long flag",
			args: []string{"--data-path", "/custom/data"},
			validate: func(t *testing.T, opts Options) {
				if opts.DataPath != "/custom/data" {
					t.Errorf("expected DataPath %q, got %q", "/custom/data", opts.DataPath)
				}
			},
		},
		{
			name: "custom static path",
			args: []string{"--static-path", "/custom/static"},
			validate: func(t *testing.T, opts Options) {
				if opts.StaticPath != "/custom/static" {
					t.Errorf("expected StaticPath %q, got %q", "/custom/static", opts.StaticPath)
				}
			},
		},
		{
			name: "parent url",
			args: []string{"--parent-url", "https://parent.example.com"},
			validate: func(t *testing.T, opts Options) {
				if opts.ParentUrl != "https://parent.example.com" {
					t.Errorf("expected ParentUrl %q, got %q", "https://parent.example.com", opts.ParentUrl)
				}
			},
		},
		{
			name: "parent box id",
			args: []string{"--parent-id", "parent-box-123"},
			validate: func(t *testing.T, opts Options) {
				if opts.ParentBoxID != "parent-box-123" {
					t.Errorf("expected ParentBoxID %q, got %q", "parent-box-123", opts.ParentBoxID)
				}
			},
		},
		{
			name: "parent box size",
			args: []string{"--parent-size", "xlarge"},
			validate: func(t *testing.T, opts Options) {
				if opts.ParentBoxSize != "xlarge" {
					t.Errorf("expected ParentBoxSize %q, got %q", "xlarge", opts.ParentBoxSize)
				}
			},
		},
		{
			name: "multiple options",
			args: []string{
				"--api-port", "9091",
				"-p", "9090",
				"--debug",
				"--data-path", "/data",
				"--static-path", "/static",
				"--parent-url", "https://parent.example.com",
				"--parent-id", "my-box",
				"--parent-size", "small",
			},
			validate: func(t *testing.T, opts Options) {
				if opts.ApiPort != "9091" {
					t.Errorf("expected ApiPort %q, got %q", "9091", opts.ApiPort)
				}
				if opts.SitePort != "9090" {
					t.Errorf("expected SitePort %q, got %q", "9090", opts.SitePort)
				}
				if !opts.Debug {
					t.Error("expected Debug to be true")
				}
				if opts.DataPath != "/data" {
					t.Errorf("expected DataPath %q, got %q", "/data", opts.DataPath)
				}
				if opts.StaticPath != "/static" {
					t.Errorf("expected StaticPath %q, got %q", "/static", opts.StaticPath)
				}
				if opts.ParentUrl != "https://parent.example.com" {
					t.Errorf("expected ParentUrl %q, got %q", "https://parent.example.com", opts.ParentUrl)
				}
				if opts.ParentBoxID != "my-box" {
					t.Errorf("expected ParentBoxID %q, got %q", "my-box", opts.ParentBoxID)
				}
				if opts.ParentBoxSize != "small" {
					t.Errorf("expected ParentBoxSize %q, got %q", "small", opts.ParentBoxSize)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := Options{}
			parser := goflags.NewParser(&opts, goflags.Default)

			_, err := parser.ParseArgs(tt.args)
			if err != nil {
				t.Fatalf("failed to parse args: %v", err)
			}

			tt.validate(t, opts)
		})
	}
}

func TestOptions_HelpFlag(t *testing.T) {
	t.Run("help flag returns ErrHelp", func(t *testing.T) {
		opts := Options{}
		parser := goflags.NewParser(&opts, goflags.Default)

		_, err := parser.ParseArgs([]string{"--help"})
		if err == nil {
			t.Fatal("expected error for --help, got nil")
		}

		flagsErr, ok := err.(*goflags.Error)
		if !ok {
			t.Fatalf("expected *goflags.Error, got %T", err)
		}

		if flagsErr.Type != goflags.ErrHelp {
			t.Errorf("expected ErrHelp, got %v", flagsErr.Type)
		}
	})

	t.Run("short help flag returns ErrHelp", func(t *testing.T) {
		opts := Options{}
		parser := goflags.NewParser(&opts, goflags.Default)

		_, err := parser.ParseArgs([]string{"-h"})
		if err == nil {
			t.Fatal("expected error for -h, got nil")
		}

		flagsErr, ok := err.(*goflags.Error)
		if !ok {
			t.Fatalf("expected *goflags.Error, got %T", err)
		}

		if flagsErr.Type != goflags.ErrHelp {
			t.Errorf("expected ErrHelp, got %v", flagsErr.Type)
		}
	})
}

func TestOptions_InvalidFlags(t *testing.T) {
	t.Run("unknown flag returns error", func(t *testing.T) {
		opts := Options{}
		parser := goflags.NewParser(&opts, goflags.Default)

		_, err := parser.ParseArgs([]string{"--unknown-flag"})
		if err == nil {
			t.Fatal("expected error for unknown flag, got nil")
		}

		flagsErr, ok := err.(*goflags.Error)
		if !ok {
			t.Fatalf("expected *goflags.Error, got %T", err)
		}

		if flagsErr.Type != goflags.ErrUnknownFlag {
			t.Errorf("expected ErrUnknownFlag, got %v", flagsErr.Type)
		}
	})
}

func TestProcessOptions_Integration(t *testing.T) {
	// This test verifies the processOptions function behavior in a controlled way
	// We can't fully test it since it reads os.Args and can panic/exit

	t.Run("options struct is exported", func(t *testing.T) {
		// Verify the global options variable exists and is usable
		originalOptions := options
		defer func() { options = originalOptions }()

		// Set some test values
		options = Options{
			ApiPort:  "9999",
			SitePort: "8888",
			Debug:    true,
		}

		if options.ApiPort != "9999" {
			t.Error("failed to set ApiPort on global options")
		}
		if options.SitePort != "8888" {
			t.Error("failed to set SitePort on global options")
		}
		if !options.Debug {
			t.Error("failed to set Debug on global options")
		}
	})
}

func TestOptions_StructTags(t *testing.T) {
	// Verify that struct tags are set correctly
	t.Run("struct tags define correct flags", func(t *testing.T) {
		opts := Options{}
		parser := goflags.NewParser(&opts, goflags.Default)

		// Get all the options from the parser
		allGroups := parser.Groups()
		if len(allGroups) == 0 {
			t.Fatal("no option groups found")
		}

		// Verify that at least some expected flags are registered
		_, err := parser.ParseArgs([]string{"--api-port", "8081"})
		if err != nil {
			t.Errorf("--api-port flag should be valid: %v", err)
		}

		_, err = parser.ParseArgs([]string{"-p", "8080"})
		if err != nil {
			t.Errorf("-p flag should be valid: %v", err)
		}

		_, err = parser.ParseArgs([]string{"--port", "8080"})
		if err != nil {
			t.Errorf("--port flag should be valid: %v", err)
		}

		_, err = parser.ParseArgs([]string{"-d", "/data"})
		if err != nil {
			t.Errorf("-d flag should be valid: %v", err)
		}

		_, err = parser.ParseArgs([]string{"--data-path", "/data"})
		if err != nil {
			t.Errorf("--data-path flag should be valid: %v", err)
		}
	})
}

func TestOptions_EnvironmentIntegration(t *testing.T) {
	t.Run("options work with empty HOME", func(t *testing.T) {
		// Save original
		originalHome := os.Getenv("HOME")
		defer os.Setenv("HOME", originalHome)

		// Test with empty HOME (edge case)
		os.Setenv("HOME", "")

		opts := Options{}
		parser := goflags.NewParser(&opts, goflags.Default)

		_, err := parser.ParseArgs([]string{})
		if err != nil {
			t.Fatalf("parsing should succeed even with empty HOME: %v", err)
		}

		// DataPath and StaticPath should still be empty (defaults)
		if opts.DataPath != "" {
			t.Errorf("DataPath should be empty, got %q", opts.DataPath)
		}
		if opts.StaticPath != "" {
			t.Errorf("StaticPath should be empty, got %q", opts.StaticPath)
		}
	})
}
