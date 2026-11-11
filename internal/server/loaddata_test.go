package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Issue #62: load the newest usable backup when boxes.json is missing or
// corrupt rather than failing to start.
func TestLoadNewestBoxFile(t *testing.T) {
	initTestLogger()
	originalFile, originalDataPath := boxFile, options.DataPath
	defer func() { boxFile, options.DataPath = originalFile, originalDataPath }()

	const good1 = `[{"id":"from-main","name":"main"}]`
	const good2 = `[{"id":"from-bak","name":"bak"},{"id":"two","name":"two"}]`
	const corrupt = `[{"id":"half-writ`

	tests := []struct {
		name       string
		files      map[string]string // file suffix ("" is boxes.json) -> content, "-" means missing
		wantSource string
		wantIDs    string
		wantErr    bool
		wantMoved  bool // corrupt boxes.json kept as boxes.json.corrupt-*
	}{
		{name: "valid data file", files: map[string]string{"": good1, ".bak1": good2}, wantSource: "", wantIDs: "from-main"},
		{name: "corrupt data file uses bak1", files: map[string]string{"": corrupt, ".bak1": good2}, wantSource: ".bak1", wantIDs: "from-bak,two", wantMoved: true},
		{name: "empty data file uses bak1", files: map[string]string{"": "", ".bak1": good2}, wantSource: ".bak1", wantIDs: "from-bak,two", wantMoved: true},
		{name: "skips corrupt backups", files: map[string]string{"": corrupt, ".bak1": corrupt, ".bak2": "", ".bak3": good1}, wantSource: ".bak3", wantIDs: "from-main", wantMoved: true},
		{name: "skips missing backups", files: map[string]string{"": corrupt, ".bak4": good2}, wantSource: ".bak4", wantIDs: "from-bak,two", wantMoved: true},
		{name: "missing data file with backups", files: map[string]string{"": "-", ".bak1": good2}, wantSource: ".bak1", wantIDs: "from-bak,two"},
		{name: "fresh install creates empty data file", files: map[string]string{"": "-"}, wantSource: "", wantIDs: ""},
		{name: "nothing usable", files: map[string]string{"": corrupt, ".bak1": corrupt}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			options.DataPath = dir
			boxFile = filepath.Join(dir, "boxes.json")
			for suffix, content := range tt.files {
				if content == "-" {
					continue
				}
				if err := os.WriteFile(boxFile+suffix, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}
			}

			createDataFiles()
			boxes, source, err := loadNewestBoxFile()

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, loaded %d boxes from %s", len(boxes), source)
				}
				if !strings.Contains(err.Error(), "boxes.json:") || !strings.Contains(err.Error(), "boxes.json.bak1:") {
					t.Errorf("error should name every file tried, got: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if source != boxFile+tt.wantSource {
				t.Errorf("loaded from %s, want %s", filepath.Base(source), filepath.Base(boxFile+tt.wantSource))
			}
			var ids []string
			for _, b := range boxes {
				ids = append(ids, b.ID)
			}
			if got := strings.Join(ids, ","); got != tt.wantIDs {
				t.Errorf("loaded boxes %q, want %q", got, tt.wantIDs)
			}

			moved, _ := filepath.Glob(boxFile + ".corrupt-*")
			if tt.wantMoved != (len(moved) == 1) {
				t.Errorf("corrupt data file kept aside = %v, want %v", moved, tt.wantMoved)
			}
			if tt.wantMoved {
				if b, _ := os.ReadFile(moved[0]); string(b) != tt.files[""] {
					t.Errorf("kept file content changed")
				}
			}
		})
	}
}
