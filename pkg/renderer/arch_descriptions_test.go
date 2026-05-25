package renderer

import (
	"strings"
	"testing"
)

// TestArchDescriptionsCoverLibrary asserts ArchDescriptions and ArchLibrary
// have identical key sets — catches drift in both directions.
func TestArchDescriptionsCoverLibrary(t *testing.T) {
	for name := range ArchLibrary {
		if _, ok := ArchDescriptions[name]; !ok {
			t.Errorf("ArchLibrary has %q but ArchDescriptions does not", name)
			continue
		}
		if strings.TrimSpace(ArchDescriptions[name]) == "" {
			t.Errorf("ArchDescriptions[%q] is empty", name)
		}
	}
	for name := range ArchDescriptions {
		if _, ok := ArchLibrary[name]; !ok {
			t.Errorf("ArchDescriptions has stale entry %q (not in ArchLibrary)", name)
		}
	}
}

// TestArchCategoryCoverLibrary asserts archCategory has every ArchLibrary key.
func TestArchCategoryCoverLibrary(t *testing.T) {
	for name := range ArchLibrary {
		if _, ok := archCategory[name]; !ok {
			t.Errorf("ArchLibrary has %q but archCategory does not", name)
		}
	}
	for name := range archCategory {
		if _, ok := ArchLibrary[name]; !ok {
			t.Errorf("archCategory has stale entry %q (not in ArchLibrary)", name)
		}
	}
}
