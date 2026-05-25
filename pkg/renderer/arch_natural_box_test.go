package renderer

import "testing"

// TestArchNaturalBoxCoverage asserts ArchNaturalBox has a measurement for
// every ArchLibrary key (and no stale entries) — the catalog JSON would
// silently fall back to DefaultNaturalBox if an entry were missing,
// which is misleading to composing agents.
func TestArchNaturalBoxCoverage(t *testing.T) {
	for name := range ArchLibrary {
		b, ok := ArchNaturalBox[name]
		if !ok {
			t.Errorf("ArchLibrary has %q but ArchNaturalBox does not — fix: add to ArchNaturalBox", name)
			continue
		}
		if b.W <= 0 || b.H <= 0 {
			t.Errorf("ArchNaturalBox[%q] has non-positive dims %+v — fix: ArchNaturalBox", name, b)
		}
		if b.AspectRatio <= 0 {
			t.Errorf("ArchNaturalBox[%q] has non-positive aspect %+v — fix: ArchNaturalBox", name, b)
		}
	}
	for name := range ArchNaturalBox {
		if _, ok := ArchLibrary[name]; !ok {
			t.Errorf("ArchNaturalBox has stale entry %q (not in ArchLibrary) — fix: remove from ArchNaturalBox", name)
		}
	}
}
