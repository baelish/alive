package server

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/baelish/alive/api"
)

func TestSeedDemoBoxes(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()
	resetBoxStore()

	seedDemoBoxes()

	counts := map[api.BoxSize]int{}
	stats := map[api.BoxSize]int{}
	for _, b := range boxStore.GetAll() {
		counts[b.Size]++
		if b.Stat != "" {
			stats[b.Size]++
		}
	}

	want := []struct {
		size         api.BoxSize
		count, stats int
	}{
		{api.Xlarge, 1, 1},
		{api.Dlarge, 1, 1},
		{api.Large, 1, 1},
		{api.Dmedium, 1, 1},
		{api.Medium, 1, 1},
		{api.Dsmall, 3, 1},
		{api.Small, 5, 2},
		{api.Dmicro, 2, 1},
		{api.Micro, 8, 4},
		{api.Dot, 32, 16},
	}
	for _, w := range want {
		if counts[w.size] != w.count {
			t.Errorf("%s: got %d boxes, want %d", w.size, counts[w.size], w.count)
		}
		if stats[w.size] != w.stats {
			t.Errorf("%s: got %d boxes with a stat, want %d", w.size, stats[w.size], w.stats)
		}
	}
}

func TestCreateMissingDemoBox(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()
	resetBoxStore()

	seedDemoBoxes()
	full := boxStore.Len()

	// Nothing is missing, so nothing is created.
	createMissingDemoBox(boxStore.GetAll())
	if got := boxStore.Len(); got != full {
		t.Fatalf("created a box when none were missing: %d boxes, want %d", got, full)
	}

	// Remove the xlarge box, it is the only size missing so it is replaced.
	for _, b := range boxStore.GetAll() {
		if b.Size == api.Xlarge {
			deleteBox(b.ID, false)
		}
	}
	createMissingDemoBox(boxStore.GetAll())
	var xlarge []api.Box
	for _, b := range boxStore.GetAll() {
		if b.Size == api.Xlarge {
			xlarge = append(xlarge, b)
		}
	}
	if len(xlarge) != 1 {
		t.Fatalf("got %d xlarge boxes after refill, want 1", len(xlarge))
	}
	if xlarge[0].Stat == "" {
		t.Errorf("replacement single xlarge box has no stat")
	}
}

func TestDemoInitials(t *testing.T) {
	tests := map[string]string{
		"Siberian Tiger":              "ST",
		"Badger":                      "BAD",
		"Elk":                         "ELK",
		"Ox":                          "OX",
		"Big-horned sheep":            "BHS",
		"Devi Fish (Giant Sting ray)": "DFG",
		"Bob-Cat":                     "BC",
		"X-ray fish":                  "XRF",
		"Bird of paradise":            "BOP",
		"":                            "",
	}
	for name, want := range tests {
		if got := demoInitials(name); got != want {
			t.Errorf("demoInitials(%q) = %q, want %q", name, got, want)
		}
	}

	// Every demo animal gives 2-3 capital letters.
	valid := regexp.MustCompile(`^[A-Z]{2,3}$`)
	for _, name := range animals {
		if got := demoInitials(name); !valid.MatchString(got) {
			t.Errorf("demoInitials(%q) = %q, want 2-3 capital letters", name, got)
		}
	}
}

func TestSeedDemoMicroDisplayNames(t *testing.T) {
	originalBoxes := boxStore.GetAll()
	defer func() {
		boxStore.mu.Lock()
		boxStore.boxes = originalBoxes
		boxStore.mu.Unlock()
	}()
	resetBoxStore()

	seedDemoBoxes()
	for _, b := range boxStore.GetAll() {
		if b.Size == api.Micro && b.DisplayName != demoInitials(b.Name) {
			t.Errorf("micro box %q has display name %q, want %q", b.Name, b.DisplayName, demoInitials(b.Name))
		}
		if b.Size != api.Micro && b.DisplayName != "" {
			t.Errorf("%s box %q has display name %q, want none", b.Size, b.Name, b.DisplayName)
		}
	}
}

func TestNextStatKeepsFormat(t *testing.T) {
	fraction := regexp.MustCompile(`^(\d+)/10$`)
	tests := []struct {
		start string
		check func(string) error
	}{
		{"3/10", func(s string) error {
			m := fraction.FindStringSubmatch(s)
			if m == nil {
				return fmt.Errorf("not out of 10")
			}
			if n, _ := strconv.Atoi(m[1]); n > 10 {
				return fmt.Errorf("over the total")
			}
			return nil
		}},
		{"0/10", func(s string) error {
			if !fraction.MatchString(s) {
				return fmt.Errorf("not out of 10")
			}
			return nil
		}},
		{"120ms", func(s string) error {
			n, err := strconv.Atoi(strings.TrimSuffix(s, "ms"))
			if !strings.HasSuffix(s, "ms") || err != nil || n < 1 {
				return fmt.Errorf("not a positive ms value")
			}
			return nil
		}},
		{"100%", func(s string) error {
			n, err := strconv.Atoi(strings.TrimSuffix(s, "%"))
			if !strings.HasSuffix(s, "%") || err != nil || n < 0 || n > 100 {
				return fmt.Errorf("not a 0-100 percentage")
			}
			return nil
		}},
		{"42.3k", func(s string) error {
			f, err := strconv.ParseFloat(strings.TrimSuffix(s, "k"), 64)
			if !strings.HasSuffix(s, "k") || err != nil || f < 0 {
				return fmt.Errorf("not a thousands value")
			}
			return nil
		}},
	}

	for _, tt := range tests {
		stat := tt.start
		for i := 0; i < 500; i++ {
			stat = nextStat(stat)
			if err := tt.check(stat); err != nil {
				t.Fatalf("starting from %q, step %d gave %q: %v", tt.start, i, stat, err)
			}
		}
	}
}

func TestNextStatUnknownFormat(t *testing.T) {
	// An unrecognised stat is either left alone or replaced by a known format.
	known := regexp.MustCompile(`^(\d+/\d+|\d+%|\d+ms|\d+\.\dk)$`)
	for i := 0; i < 100; i++ {
		if got := nextStat("hello"); got != "hello" && !known.MatchString(got) {
			t.Fatalf("got %q, want \"hello\" or a known stat format", got)
		}
	}
}
