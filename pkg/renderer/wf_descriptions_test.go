package renderer

import (
	"strings"
	"testing"
)

// TestWFDescriptionsCoverLibrary asserts that WFDescriptions and WFLibrary
// have identical key sets — catches drift in both directions when someone
// adds a wireframe without a description (or vice versa).
func TestWFDescriptionsCoverLibrary(t *testing.T) {
	for name := range WFLibrary {
		if _, ok := WFDescriptions[name]; !ok {
			t.Errorf("WFLibrary has %q but WFDescriptions does not", name)
			continue
		}
		if strings.TrimSpace(WFDescriptions[name]) == "" {
			t.Errorf("WFDescriptions[%q] is empty", name)
		}
	}
	for name := range WFDescriptions {
		if _, ok := WFLibrary[name]; !ok {
			t.Errorf("WFDescriptions has stale entry %q (not in WFLibrary)", name)
		}
	}
}

// TestWFCategoryCoverLibrary asserts wfCategory has every WFLibrary key.
// The category map is the source of truth for `catalog wf --json`.
func TestWFCategoryCoverLibrary(t *testing.T) {
	for name := range WFLibrary {
		if _, ok := wfCategory[name]; !ok {
			t.Errorf("WFLibrary has %q but wfCategory does not", name)
		}
	}
	for name := range wfCategory {
		if _, ok := WFLibrary[name]; !ok {
			t.Errorf("wfCategory has stale entry %q (not in WFLibrary)", name)
		}
	}
}
