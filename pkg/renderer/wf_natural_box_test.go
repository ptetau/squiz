package renderer

import "testing"

// TestWFNaturalBoxCoverage asserts WFNaturalBox has a measurement for
// every WFLibrary key (and no stale entries) — keeps the catalog JSON
// honest as new wireframes get added.
func TestWFNaturalBoxCoverage(t *testing.T) {
	for name := range WFLibrary {
		b, ok := WFNaturalBox[name]
		if !ok {
			t.Errorf("WFLibrary has %q but WFNaturalBox does not — fix: add to WFNaturalBox", name)
			continue
		}
		if b.W <= 0 || b.H <= 0 {
			t.Errorf("WFNaturalBox[%q] has non-positive dims %+v — fix: WFNaturalBox", name, b)
		}
		if b.AspectRatio <= 0 {
			t.Errorf("WFNaturalBox[%q] has non-positive aspect %+v — fix: WFNaturalBox", name, b)
		}
	}
	for name := range WFNaturalBox {
		if _, ok := WFLibrary[name]; !ok {
			t.Errorf("WFNaturalBox has stale entry %q (not in WFLibrary) — fix: remove from WFNaturalBox", name)
		}
	}
}
