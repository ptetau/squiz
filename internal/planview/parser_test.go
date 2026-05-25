package planview

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadPlan_Example exercises the worked fixture and pins down
// title/lede plus per-section item counts in canonical order.
func TestLoadPlan_Example(t *testing.T) {
	p, err := LoadPlan(filepath.Join("..", "..", "testdata", "plan-example", "index.json"))
	if err != nil {
		t.Fatalf("LoadPlan: %v", err)
	}

	if p.Title != "ThermoLog — home temperature logger" {
		t.Errorf("Title = %q", p.Title)
	}
	if !strings.HasPrefix(p.Lede, "A small offline-first") {
		t.Errorf("Lede = %q", p.Lede)
	}
	if p.Theme != "paper" {
		t.Errorf("Theme = %q, want %q", p.Theme, "paper")
	}
	if p.Density != "compact" {
		t.Errorf("Density = %q, want %q", p.Density, "compact")
	}

	wantOrder := []string{
		"overview", "functional", "non-functional",
		"cases", "engineering", "build",
	}
	wantCounts := map[string]int{
		"overview":       3,
		"functional":     4,
		"non-functional": 3,
		"cases":          3,
		"engineering":    4,
		"build":          4,
	}

	if len(p.Sections) != len(wantOrder) {
		t.Fatalf("got %d sections, want %d", len(p.Sections), len(wantOrder))
	}
	for i, s := range p.Sections {
		if s.ID != wantOrder[i] {
			t.Errorf("section[%d].ID = %q, want %q", i, s.ID, wantOrder[i])
		}
		if got, want := len(s.Items), wantCounts[s.ID]; got != want {
			t.Errorf("section %q: %d items, want %d", s.ID, got, want)
		}
	}

	// Spot-check a label.
	if p.Sections[2].Label != "Non-functional" {
		t.Errorf("non-functional label = %q", p.Sections[2].Label)
	}
}

// TestLoadPlan_PreservesCanonicalOrder verifies that canonical sections
// are sorted into CanonicalSections order regardless of how the index
// declares them.
func TestLoadPlan_PreservesCanonicalOrder(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, dir, "index.json", `{
		"title": "Reordered",
		"lede": "",
		"sections": ["build", "overview", "functional"]
	}`)
	writeJSON(t, dir, "overview.json", `{"items":[{"id":"OVR-1","title":"M","desc":"d"}]}`)
	writeJSON(t, dir, "functional.json", `{"items":[{"id":"FR-1","title":"F","desc":"d"}]}`)
	writeJSON(t, dir, "build.json", `{"items":[{"id":"BUILD-1","title":"B","desc":"d"}]}`)

	p, err := LoadPlan(filepath.Join(dir, "index.json"))
	if err != nil {
		t.Fatalf("LoadPlan: %v", err)
	}
	gotIDs := sectionIDs(p)
	want := []string{"overview", "functional", "build"}
	if !equalStrings(gotIDs, want) {
		t.Errorf("section order = %v, want %v", gotIDs, want)
	}
}

// TestLoadPlan_AppendsCustomSections verifies that non-canonical
// sections are appended after the canonical block, preserving the
// order they were declared in.
func TestLoadPlan_AppendsCustomSections(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, dir, "index.json", `{
		"title": "Plan with extras",
		"lede": "",
		"sections": ["overview", "glossary", "functional", "risks"]
	}`)
	writeJSON(t, dir, "overview.json", `{"items":[{"id":"OVR-1","title":"O","desc":"d"}]}`)
	writeJSON(t, dir, "functional.json", `{"items":[{"id":"FR-1","title":"F","desc":"d"}]}`)
	// Custom sections: any IDs are fine (no prefix enforcement).
	writeJSON(t, dir, "glossary.json", `{"items":[{"id":"glossary-a","title":"A","desc":"d"}]}`)
	writeJSON(t, dir, "risks.json", `{"items":[{"id":"risk-1","title":"R","desc":"d"}]}`)

	p, err := LoadPlan(filepath.Join(dir, "index.json"))
	if err != nil {
		t.Fatalf("LoadPlan: %v", err)
	}
	gotIDs := sectionIDs(p)
	want := []string{"overview", "functional", "glossary", "risks"}
	if !equalStrings(gotIDs, want) {
		t.Errorf("section order = %v, want %v", gotIDs, want)
	}

	// Custom labels should be title-cased.
	for _, s := range p.Sections {
		if s.ID == "glossary" && s.Label != "Glossary" {
			t.Errorf("glossary label = %q, want %q", s.Label, "Glossary")
		}
		if s.ID == "risks" && s.Label != "Risks" {
			t.Errorf("risks label = %q, want %q", s.Label, "Risks")
		}
	}
}

// TestLoadPlan_RejectsBadRef verifies that an item referencing a
// non-existent ID is rejected, and that the error mentions both the
// offending item and the missing target.
func TestLoadPlan_RejectsBadRef(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, dir, "index.json", `{
		"title": "BadRef", "lede": "",
		"sections": ["overview", "functional"]
	}`)
	writeJSON(t, dir, "overview.json", `{"items":[{"id":"OVR-1","title":"O","desc":"d"}]}`)
	writeJSON(t, dir, "functional.json", `{"items":[
		{"id":"FR-1","title":"F","desc":"d","refs":["OVR-99"]}
	]}`)

	_, err := LoadPlan(filepath.Join(dir, "index.json"))
	if err == nil {
		t.Fatal("LoadPlan: expected error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "FR-1") {
		t.Errorf("error %q should name offending item FR-1", msg)
	}
	if !strings.Contains(msg, "OVR-99") {
		t.Errorf("error %q should name missing target OVR-99", msg)
	}
}

// TestLoadPlan_RejectsDuplicateID checks both same-file and
// cross-file duplicates.
func TestLoadPlan_RejectsDuplicateID(t *testing.T) {
	t.Run("same file", func(t *testing.T) {
		dir := t.TempDir()
		writeJSON(t, dir, "index.json", `{
			"title": "Dup", "lede": "",
			"sections": ["engineering"]
		}`)
		writeJSON(t, dir, "engineering.json", `{"items":[
			{"id":"ENG-1","title":"A","desc":"d"},
			{"id":"ENG-1","title":"B","desc":"d"}
		]}`)
		_, err := LoadPlan(filepath.Join(dir, "index.json"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "ENG-1") {
			t.Errorf("error %q should name duplicate id ENG-1", err)
		}
	})

	t.Run("cross file (custom section avoids prefix collision)", func(t *testing.T) {
		// Use a custom section so we can safely place an OVR-* id there
		// without tripping the prefix check before the dup check runs.
		dir := t.TempDir()
		writeJSON(t, dir, "index.json", `{
			"title": "Dup", "lede": "",
			"sections": ["overview", "glossary"]
		}`)
		writeJSON(t, dir, "overview.json", `{"items":[{"id":"OVR-1","title":"O","desc":"d"}]}`)
		writeJSON(t, dir, "glossary.json", `{"items":[{"id":"OVR-1","title":"clone","desc":"d"}]}`)
		_, err := LoadPlan(filepath.Join(dir, "index.json"))
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		msg := err.Error()
		if !strings.Contains(msg, "OVR-1") {
			t.Errorf("error %q should name duplicate id OVR-1", msg)
		}
	})
}

// TestLoadPlan_RejectsWrongPrefix verifies the section-prefix check.
func TestLoadPlan_RejectsWrongPrefix(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, dir, "index.json", `{
		"title": "Misfiled", "lede": "",
		"sections": ["functional"]
	}`)
	writeJSON(t, dir, "functional.json", `{"items":[
		{"id":"BUILD-1","title":"misplaced","desc":"d"}
	]}`)
	_, err := LoadPlan(filepath.Join(dir, "index.json"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "BUILD-1") {
		t.Errorf("error %q should name offending id BUILD-1", msg)
	}
	if !strings.Contains(msg, "FR") {
		t.Errorf("error %q should name expected prefix FR", msg)
	}
}

// TestLoadPlan_MissingSectionFile names the file we tried to open.
func TestLoadPlan_MissingSectionFile(t *testing.T) {
	dir := t.TempDir()
	writeJSON(t, dir, "index.json", `{
		"title": "Missing", "lede": "",
		"sections": ["overview", "mystery"]
	}`)
	writeJSON(t, dir, "overview.json", `{"items":[{"id":"OVR-1","title":"O","desc":"d"}]}`)
	// no mystery.json
	_, err := LoadPlan(filepath.Join(dir, "index.json"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "mystery.json") {
		t.Errorf("error %q should name missing file mystery.json", err)
	}
}

func TestLookupItem(t *testing.T) {
	p, err := LoadPlan(filepath.Join("..", "..", "testdata", "plan-example", "index.json"))
	if err != nil {
		t.Fatalf("LoadPlan: %v", err)
	}

	t.Run("hit", func(t *testing.T) {
		it, sid, ok := p.LookupItem("FR-2")
		if !ok {
			t.Fatal("expected to find FR-2")
		}
		if sid != "functional" {
			t.Errorf("section = %q, want functional", sid)
		}
		if it.ID != "FR-2" {
			t.Errorf("item.ID = %q", it.ID)
		}
		if it.Title == "" {
			t.Error("item.Title is empty")
		}
	})

	t.Run("miss", func(t *testing.T) {
		_, _, ok := p.LookupItem("ZZZ-999")
		if ok {
			t.Error("expected miss")
		}
	})

	t.Run("case sensitive", func(t *testing.T) {
		_, _, ok := p.LookupItem("fr-2")
		if ok {
			t.Error("LookupItem should be case-sensitive")
		}
	})
}

func TestParseIndex_Defaults(t *testing.T) {
	idx, err := ParseIndex([]byte(`{"sections":["overview"]}`))
	if err != nil {
		t.Fatalf("ParseIndex: %v", err)
	}
	if idx.Density != "compact" {
		t.Errorf("Density = %q, want %q", idx.Density, "compact")
	}
	if idx.Theme != "" {
		t.Errorf("Theme = %q, want empty (auto-rotation)", idx.Theme)
	}
	if len(idx.Sections) != 1 || idx.Sections[0] != "overview" {
		t.Errorf("Sections = %v", idx.Sections)
	}
}

// --- helpers --------------------------------------------------------------

func writeJSON(t *testing.T, dir, name, body string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func sectionIDs(p *Plan) []string {
	out := make([]string, len(p.Sections))
	for i, s := range p.Sections {
		out[i] = s.ID
	}
	return out
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
