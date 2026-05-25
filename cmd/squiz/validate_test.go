package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain_Validate_Valid runs `validate` against the canonical smoke
// fixture and expects exit 0 + a "valid" summary line.
func TestMain_Validate_Valid(t *testing.T) {
	bin := buildBinary(t)
	root := repoRoot(t)
	input := filepath.Join(root, "testdata", "smoke.json")

	cmd := exec.Command(bin, "validate", input)
	combined, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("validate on smoke.json failed: %v\n%s", err, combined)
	}
	if exit := cmd.ProcessState.ExitCode(); exit != 0 {
		t.Fatalf("exit = %d, want 0\n%s", exit, combined)
	}
	if !strings.Contains(string(combined), "valid (") {
		t.Errorf("missing valid summary in output: %s", combined)
	}
}

// TestMain_Validate_Invalid writes a deliberately-broken doc to a
// tempfile (duplicate squiz id + duplicate option id + bad theme) and
// asserts the binary exits 1 and reports the expected error paths.
func TestMain_Validate_Invalid(t *testing.T) {
	bin := buildBinary(t)
	broken := `{
		"theme": "neon-pink",
		"spec": {"title": "x", "paragraphs": [{"text": "see {{nosuch}}"}]},
		"squizzes": [
			{"id": "dup", "options": [
				{"id": "a", "name": "A", "desc": "x"},
				{"id": "a", "name": "B", "desc": "y"}
			]},
			{"id": "dup", "options": []}
		]
	}`
	path := filepath.Join(t.TempDir(), "broken.json")
	if err := os.WriteFile(path, []byte(broken), 0644); err != nil {
		t.Fatalf("write broken fixture: %v", err)
	}

	cmd := exec.Command(bin, "validate", path)
	combined, _ := cmd.CombinedOutput()
	if exit := cmd.ProcessState.ExitCode(); exit != 1 {
		t.Fatalf("exit = %d, want 1\n%s", exit, combined)
	}
	out := string(combined)
	for _, want := range []string{
		"theme",
		"unknown theme",
		"squizzes[1].id",
		"duplicate id",
		"squizzes[0].options[1].id",
		"invalid (",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\nfull output:\n%s", want, out)
		}
	}
	// Marker check is a WARNING — should not be counted in error total.
	if !strings.Contains(out, "warning:") || !strings.Contains(out, "{{nosuch}}") {
		t.Errorf("expected a warning about {{nosuch}} marker\nfull output:\n%s", out)
	}
}

// TestMain_Validate_JSON asserts --json emits the report struct shape
// rather than the text format.
func TestMain_Validate_JSON(t *testing.T) {
	bin := buildBinary(t)
	root := repoRoot(t)
	input := filepath.Join(root, "testdata", "smoke.json")

	cmd := exec.Command(bin, "validate", input, "--json")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("validate --json failed: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("--json output not valid JSON: %v\n%s", err, out)
	}
	if got["valid"] != true {
		t.Errorf("valid field = %v, want true", got["valid"])
	}
	if _, ok := got["errors"]; !ok {
		t.Errorf("missing errors key in JSON output")
	}
}
