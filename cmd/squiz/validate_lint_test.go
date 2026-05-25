package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// writeTempDoc serialises a tiny Document fixture to a temp file and
// returns its path. Used by every TestValidateLint_* test below — keeps
// the fixture inline so the trigger pattern is visible at the test site.
func writeTempDoc(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "doc.json")
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}

// runValidate runs `squiz validate <path>` and returns combined output +
// exit code. Lints are warnings only — exit must stay 0.
func runValidate(t *testing.T, args ...string) (string, int) {
	t.Helper()
	bin := buildBinary(t)
	cmd := exec.Command(bin, append([]string{"validate"}, args...)...)
	combined, _ := cmd.CombinedOutput()
	return string(combined), cmd.ProcessState.ExitCode()
}

// TestValidateLint_UnknownUseRef — a composed raw SVG that <use href="wf:nonsense"/>s
// a name not in the library must produce an unknown-use-ref warning at
// exit 0.
func TestValidateLint_UnknownUseRef(t *testing.T) {
	doc := `{
		"spec": {"title": "t"},
		"squizzes": [
			{"id": "s1", "options": [
				{"id": "a", "name": "A", "desc": "x",
				 "art": "<svg viewBox='0 0 100 60'><use href=\"wf:nonsense\"/></svg>"}
			]}
		]
	}`
	path := writeTempDoc(t, doc)
	out, exit := runValidate(t, path)
	if exit != 0 {
		t.Fatalf("exit = %d, want 0 (lints are warnings)\n%s", exit, out)
	}
	for _, want := range []string{"warn:", "unknown-use-ref", "wf:nonsense"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

// TestValidateLint_OOBUseBox — <use x="80" y="10" width="40" height="60"/>
// has x+width=120 > 100, so oob-use-box must fire at exit 0.
func TestValidateLint_OOBUseBox(t *testing.T) {
	doc := `{
		"spec": {"title": "t"},
		"squizzes": [
			{"id": "s1", "options": [
				{"id": "a", "name": "A", "desc": "x",
				 "art": "<svg viewBox='0 0 100 60'><use href=\"wf:gauge\" x=\"80\" y=\"10\" width=\"40\" height=\"60\"/></svg>"}
			]}
		]
	}`
	path := writeTempDoc(t, doc)
	out, exit := runValidate(t, path)
	if exit != 0 {
		t.Fatalf("exit = %d, want 0\n%s", exit, out)
	}
	if !strings.Contains(out, "oob-use-box") {
		t.Errorf("missing oob-use-box warning\n%s", out)
	}
}

// TestValidateLint_MissingViewBox — a top-level <svg> with no viewBox attr
// must trigger missing-viewbox at exit 0.
func TestValidateLint_MissingViewBox(t *testing.T) {
	doc := `{
		"spec": {"title": "t"},
		"squizzes": [
			{"id": "s1", "options": [
				{"id": "a", "name": "A", "desc": "x",
				 "art": "<svg><circle cx='10' cy='10' r='3'/></svg>"}
			]}
		]
	}`
	path := writeTempDoc(t, doc)
	out, exit := runValidate(t, path)
	if exit != 0 {
		t.Fatalf("exit = %d, want 0\n%s", exit, out)
	}
	if !strings.Contains(out, "missing-viewbox") {
		t.Errorf("missing missing-viewbox warning\n%s", out)
	}
}

// TestValidateLint_SiblingCollision — two options in one squiz with the
// same art string must trigger sibling-art-collision naming both ids.
func TestValidateLint_SiblingCollision(t *testing.T) {
	doc := `{
		"spec": {"title": "t"},
		"squizzes": [
			{"id": "store", "options": [
				{"id": "sqlite",   "name": "SQLite",   "desc": "x", "art": "arch:database"},
				{"id": "postgres", "name": "Postgres", "desc": "y", "art": "arch:database"}
			]}
		]
	}`
	path := writeTempDoc(t, doc)
	out, exit := runValidate(t, path)
	if exit != 0 {
		t.Fatalf("exit = %d, want 0\n%s", exit, out)
	}
	for _, want := range []string{"sibling-art-collision", "sqlite", "postgres", "arch:database"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

// TestValidateLint_JSONReportShape — --json output must carry warnings in
// the `warnings` array alongside errors. valid must remain true (lints
// don't make a doc invalid).
func TestValidateLint_JSONReportShape(t *testing.T) {
	doc := `{
		"spec": {"title": "t"},
		"squizzes": [
			{"id": "s", "options": [
				{"id": "a", "name": "A", "desc": "x", "art": "arch:database"},
				{"id": "b", "name": "B", "desc": "y", "art": "arch:database"}
			]}
		]
	}`
	path := writeTempDoc(t, doc)
	bin := buildBinary(t)
	out, err := exec.Command(bin, "validate", path, "--json").Output()
	if err != nil {
		t.Fatalf("validate --json: %v", err)
	}
	var rep struct {
		Valid    bool `json:"valid"`
		Errors   []struct{ Path, Message string }
		Warnings []struct{ Path, Message string }
	}
	if err := json.Unmarshal(out, &rep); err != nil {
		t.Fatalf("parse JSON: %v\n%s", err, out)
	}
	if !rep.Valid {
		t.Errorf("valid = false, want true (lints don't invalidate)")
	}
	if len(rep.Warnings) == 0 {
		t.Errorf("warnings array empty, want ≥1 sibling-art-collision\n%s", out)
	}
}
