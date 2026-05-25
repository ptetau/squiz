package renderer

import (
	"strings"
	"testing"
)

// TestWFCatalog_NaturalBoxPresent confirms WFCatalog populates NaturalBox
// from the WFNaturalBox table (and not the fallback) for a known entry.
// phone-card is the canary because the v0.8.0 docs use it as the
// composition example.
func TestWFCatalog_NaturalBoxPresent(t *testing.T) {
	entries := WFCatalog()
	var found *CatalogEntry
	for i := range entries {
		if entries[i].Name == "phone-card" {
			found = &entries[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("WFCatalog missing phone-card entry")
	}
	if found.NaturalBox.W <= 0 || found.NaturalBox.H <= 0 {
		t.Errorf("phone-card NaturalBox = %+v, want non-zero W/H", found.NaturalBox)
	}
	if found.NaturalBox.AspectRatio <= 0 {
		t.Errorf("phone-card NaturalBox aspect = %v, want non-zero", found.NaturalBox.AspectRatio)
	}
}

// TestArchCatalog_NaturalBoxPresent — same check for arch entries.
func TestArchCatalog_NaturalBoxPresent(t *testing.T) {
	entries := ArchCatalog()
	var found *CatalogEntry
	for i := range entries {
		if entries[i].Name == "database" {
			found = &entries[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("ArchCatalog missing database entry")
	}
	if found.NaturalBox.W <= 0 || found.NaturalBox.H <= 0 {
		t.Errorf("database NaturalBox = %+v, want non-zero W/H", found.NaturalBox)
	}
}

// TestFormatCatalogJSON_IncludesNaturalBox — once NaturalBox is on
// CatalogEntry, JSON output must surface it (the json tag handles this;
// this is a smoke test that catches accidental tag changes).
func TestFormatCatalogJSON_IncludesNaturalBox(t *testing.T) {
	entries := WFCatalog()
	js, err := FormatCatalogJSON(entries[:1])
	if err != nil {
		t.Fatalf("FormatCatalogJSON: %v", err)
	}
	if !strings.Contains(js, `"naturalBox"`) {
		t.Errorf("JSON missing naturalBox field:\n%s", js)
	}
	if !strings.Contains(js, `"w"`) || !strings.Contains(js, `"h"`) || !strings.Contains(js, `"aspect"`) {
		t.Errorf("JSON naturalBox missing w/h/aspect fields:\n%s", js)
	}
}
