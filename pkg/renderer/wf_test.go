package renderer

import (
	"strings"
	"testing"
)

// TestWFLibrary_AllNamesResolve walks every entry in WFLibrary and
// asserts resolveNamed returns the (non-empty) SVG with hidden=false.
// Catches drift if someone adds a name without a body.
func TestWFLibrary_AllNamesResolve(t *testing.T) {
	if len(WFLibrary) == 0 {
		t.Fatal("WFLibrary is empty — expected ~50 entries")
	}
	for name, body := range WFLibrary {
		t.Run(name, func(t *testing.T) {
			svg, hidden := resolveNamed(name)
			if hidden {
				t.Errorf("resolveNamed(%q) hidden = true, want false", name)
			}
			if svg == "" {
				t.Errorf("resolveNamed(%q) returned empty SVG", name)
			}
			if !strings.Contains(svg, "<svg") {
				t.Errorf("resolveNamed(%q) output not an SVG fragment: %s", name, svg)
			}
			// Sanity: registry body should match resolved output exactly.
			if svg != body {
				t.Errorf("resolveNamed(%q) returned different body than registry entry", name)
			}
		})
	}
}

// TestWFLibrary_HasExpectedCoreNames spot-checks a handful of names that
// the spec lists as foundational. Catches accidental rename/removal.
func TestWFLibrary_HasExpectedCoreNames(t *testing.T) {
	mustHave := []string{
		"calendar-grid",
		"streak-counter",
		"spark-rising",
		"bars-up",
		"avatar-single",
		"phone-blank",
		"toggle-on",
		"toggle-off",
		"lock",
		"check-large",
	}
	for _, name := range mustHave {
		if _, ok := WFLibrary[name]; !ok {
			t.Errorf("WFLibrary missing expected entry %q", name)
		}
	}
}

// TestWFLibrary_UnknownName — an unknown name returns a non-empty SVG
// placeholder (so layout stays consistent) and is NOT hidden. The
// placeholder must surface the bad name so the author notices.
func TestWFLibrary_UnknownName(t *testing.T) {
	svg, hidden := resolveNamed("does-not-exist")
	if hidden {
		t.Error("resolveNamed(unknown) hidden = true, want false")
	}
	if svg == "" {
		t.Fatal("resolveNamed(unknown) returned empty SVG")
	}
	if !strings.Contains(svg, "<svg") {
		t.Errorf("resolveNamed(unknown) not an SVG: %s", svg)
	}
	// The placeholder echoes the requested name back so typos are visible.
	if !strings.Contains(svg, "does-not-exist") {
		t.Errorf("expected unknown name echoed in placeholder, got:\n%s", svg)
	}
	// And signals the "wf:" prefix context.
	if !strings.Contains(svg, "wf:") {
		t.Errorf("expected 'wf:' marker in placeholder, got:\n%s", svg)
	}
}

// TestWFLibrary_PhoneScreenWrapper validates the phoneScreen helper
// wraps a body fragment in the expected phone-frame chrome.
func TestWFLibrary_PhoneScreenWrapper(t *testing.T) {
	out := phoneScreen(`<rect x='0' y='0' width='1' height='1'/>`)
	if !strings.Contains(out, "<svg") {
		t.Fatalf("phoneScreen not wrapped in SVG: %s", out)
	}
	if !strings.Contains(out, "<rect") {
		t.Errorf("phoneScreen lost body content: %s", out)
	}
	// Phone frame should include the rounded surface rect.
	if !strings.Contains(out, "rx='4'") {
		t.Errorf("phoneScreen missing phone-frame rx='4': %s", out)
	}
}

// TestWFLibrary_GridCellsHelper exercises the gridCells builder used by
// calendar-grid.
func TestWFLibrary_GridCellsHelper(t *testing.T) {
	out := gridCells(3, 2, 0, 0, 4, 4, []int{0, 2})
	if !strings.Contains(out, "<rect") {
		t.Fatalf("gridCells produced no rects: %s", out)
	}
	if !strings.Contains(out, "var(--accent)") {
		t.Errorf("expected at least one filled cell using --accent, got:\n%s", out)
	}
	if !strings.Contains(out, "var(--rule-2)") {
		t.Errorf("expected at least one outlined cell using --rule-2, got:\n%s", out)
	}
}
