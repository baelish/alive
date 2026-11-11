package server

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/baelish/alive/api"
)

func TestSaveBoxFileKeepsBackups(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	originalFile := boxFile
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
		boxFile = originalFile
	}()

	dir := t.TempDir()
	boxFile = filepath.Join(dir, "boxes.json")

	read := func(name string) []api.Box {
		t.Helper()
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("reading %s: %v", filepath.Base(name), err)
		}
		var boxes []api.Box
		if err := json.Unmarshal(b, &boxes); err != nil {
			t.Fatalf("parsing %s: %v", filepath.Base(name), err)
		}
		return boxes
	}

	for i, ids := range [][]string{{"a"}, {"a", "b"}, {"a", "b", "c"}} {
		resetBoxStore()
		for _, id := range ids {
			boxStore.Add(api.Box{ID: id, Name: id})
		}
		if err := saveBoxFile(); err != nil {
			t.Fatalf("save %d: %v", i+1, err)
		}
	}

	if got := len(read(boxFile)); got != 3 {
		t.Errorf("boxes.json has %d boxes, want 3", got)
	}
	if got := len(read(boxFile + ".bak1")); got != 2 {
		t.Errorf("boxes.json.bak1 has %d boxes, want 2 (the previous save)", got)
	}
	if got := len(read(boxFile + ".bak2")); got != 1 {
		t.Errorf("boxes.json.bak2 has %d boxes, want 1", got)
	}
	if _, err := os.Stat(boxFile + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("temporary file left behind: %v", err)
	}

	// Only nine backups are kept.
	for i := 0; i < 12; i++ {
		if err := saveBoxFile(); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(boxFile + ".bak9"); err != nil {
		t.Errorf("boxes.json.bak9 missing: %v", err)
	}
	if _, err := os.Stat(boxFile + ".bak10"); !os.IsNotExist(err) {
		t.Errorf("boxes.json.bak10 should not exist")
	}
}
