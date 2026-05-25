package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// writePlan materialises a multi-file plan tree (index.json + one section
// file) in a temp dir and returns the index path. Most tests use the
// "overview" section so item IDs need the OVR-* prefix.
func writePlan(t *testing.T, sectionID, sectionBody string) string {
	t.Helper()
	dir := t.TempDir()
	index := `{"title": "t", "sections": ["` + sectionID + `"]}`
	if err := os.WriteFile(filepath.Join(dir, "index.json"), []byte(index), 0644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, sectionID+".json"), []byte(sectionBody), 0644); err != nil {
		t.Fatalf("write section: %v", err)
	}
	return filepath.Join(dir, "index.json")
}

func runValidate(t *testing.T, args ...string) (string, int) {
	t.Helper()
	bin := buildBinary(t)
	cmd := exec.Command(bin, append([]string{"validate"}, args...)...)
	combined, _ := cmd.CombinedOutput()
	return string(combined), cmd.ProcessState.ExitCode()
}

// TestValidateLint_UnknownUseRef — composed raw SVG with a typo'd library
// ref triggers unknown-use-ref at exit 0.
func TestValidateLint_UnknownUseRef(t *testing.T) {
	body := `{"items": [
		{"id": "OVR-1", "title": "x", "desc": "x",
		 "art": "<svg viewBox='0 0 100 60'><use href=\"wf:nonsense\"/></svg>"}
	]}`
	path := writePlan(t, "overview", body)
	out, exit := runValidate(t, path)
	if exit != 0 {
		t.Fatalf("exit = %d, want 0\n%s", exit, out)
	}
	for _, want := range []string{"warn:", "unknown-use-ref", "wf:nonsense"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

// TestValidateLint_OOBUseBox — x+width > 100 triggers oob-use-box.
func TestValidateLint_OOBUseBox(t *testing.T) {
	body := `{"items": [
		{"id": "OVR-1", "title": "x", "desc": "x",
		 "art": "<svg viewBox='0 0 100 60'><use href=\"wf:gauge\" x=\"80\" y=\"10\" width=\"40\" height=\"60\"/></svg>"}
	]}`
	path := writePlan(t, "overview", body)
	out, exit := runValidate(t, path)
	if exit != 0 {
		t.Fatalf("exit = %d, want 0\n%s", exit, out)
	}
	if !strings.Contains(out, "oob-use-box") {
		t.Errorf("missing oob-use-box warning\n%s", out)
	}
}

// TestValidateLint_MissingViewBox — top-level <svg> with no viewBox warns.
func TestValidateLint_MissingViewBox(t *testing.T) {
	body := `{"items": [
		{"id": "OVR-1", "title": "x", "desc": "x",
		 "art": "<svg><circle cx='10' cy='10' r='3'/></svg>"}
	]}`
	path := writePlan(t, "overview", body)
	out, exit := runValidate(t, path)
	if exit != 0 {
		t.Fatalf("exit = %d, want 0\n%s", exit, out)
	}
	if !strings.Contains(out, "missing-viewbox") {
		t.Errorf("missing missing-viewbox warning\n%s", out)
	}
}

// TestValidateLint_SiblingCollision — two options inside one item-chooser
// sharing the same art string trigger sibling-art-collision.
func TestValidateLint_SiblingCollision(t *testing.T) {
	body := `{"items": [
		{"id": "OVR-1", "title": "x", "desc": "x", "art": "wf:phone-blank",
		 "options": [
			{"id": "sqlite",   "name": "SQLite",   "desc": "x", "art": "arch:database"},
			{"id": "postgres", "name": "Postgres", "desc": "y", "art": "arch:database"}
		 ]}
	]}`
	path := writePlan(t, "overview", body)
	out, exit := runValidate(t, path)
	if exit != 0 {
		t.Fatalf("exit = %d, want 0\n%s", exit, out)
	}
	for _, want := range []string{"sibling-art-collision", "sqlite", "postgres"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

// TestValidateLint_CompositionThin — 5/5 single-token items in one
// section trigger composition-thin (5/5 = 100% > 60%, ≥3 items).
func TestValidateLint_CompositionThin(t *testing.T) {
	body := `{"items": [
		{"id": "OVR-1", "title": "a", "desc": "a", "art": "wf:phone-blank"},
		{"id": "OVR-2", "title": "b", "desc": "b", "art": "wf:phone-list"},
		{"id": "OVR-3", "title": "c", "desc": "c", "art": "wf:phone-card"},
		{"id": "OVR-4", "title": "d", "desc": "d", "art": "wf:phone-input"},
		{"id": "OVR-5", "title": "e", "desc": "e", "art": "wf:phone-tabs"}
	]}`
	path := writePlan(t, "overview", body)
	out, exit := runValidate(t, path)
	if exit != 0 {
		t.Fatalf("exit = %d, want 0\n%s", exit, out)
	}
	if !strings.Contains(out, "composition-thin") {
		t.Errorf("missing composition-thin warning\n%s", out)
	}
}

// TestValidateLint_SectionCollision — 3 items in a section sharing art
// "arch:server" trigger section-art-collision.
func TestValidateLint_SectionCollision(t *testing.T) {
	body := `{"items": [
		{"id": "ENG-1", "title": "a", "desc": "a", "art": "arch:server"},
		{"id": "ENG-2", "title": "b", "desc": "b", "art": "arch:server"},
		{"id": "ENG-3", "title": "c", "desc": "c", "art": "arch:server"}
	]}`
	path := writePlan(t, "engineering", body)
	out, exit := runValidate(t, path)
	if exit != 0 {
		t.Fatalf("exit = %d, want 0\n%s", exit, out)
	}
	for _, want := range []string{"section-art-collision", "ENG-1", "ENG-2", "ENG-3", "arch:server"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

// TestValidateLint_NoWarningsOnGoodFixture loads the real plan-example
// fixture and asserts the lints DON'T error (exit 0). Warnings are
// allowed — we record the current baseline in the test log so a follow-up
// in v0.8.1 can drive it down.
func TestValidateLint_NoWarningsOnGoodFixture(t *testing.T) {
	root := repoRoot(t)
	input := filepath.Join(root, "testdata", "plan-example", "index.json")
	if _, err := os.Stat(input); err != nil {
		t.Fatalf("plan fixture missing at %s: %v", input, err)
	}
	out, exit := runValidate(t, input)
	if exit != 0 {
		t.Fatalf("exit = %d, want 0 on good fixture\n%s", exit, out)
	}
	warnings := strings.Count(out, "warn:")
	t.Logf("plan-example baseline: %d lint warning(s)", warnings)
	if !strings.Contains(out, "valid (") {
		t.Errorf("missing valid summary\n%s", out)
	}
}
