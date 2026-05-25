package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain_Validate_Valid runs `validate` against the canonical plan
// fixture and expects exit 0 + a "valid" summary line.
func TestMain_Validate_Valid(t *testing.T) {
	bin := buildBinary(t)
	root := repoRoot(t)
	input := filepath.Join(root, "testdata", "plan-example", "index.json")
	if _, err := os.Stat(input); err != nil {
		t.Fatalf("plan fixture missing at %s: %v", input, err)
	}

	cmd := exec.Command(bin, "validate", input)
	combined, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("validate on plan-example failed: %v\n%s", err, combined)
	}
	if exit := cmd.ProcessState.ExitCode(); exit != 0 {
		t.Fatalf("exit = %d, want 0\n%s", exit, combined)
	}
	if !strings.Contains(string(combined), "valid (") {
		t.Errorf("missing valid summary in output: %s", combined)
	}
}

// TestMain_Validate_Invalid writes a plan tree with a deliberately
// missing section file to a tempdir and asserts the binary exits 1.
func TestMain_Validate_Invalid(t *testing.T) {
	bin := buildBinary(t)

	planDir := t.TempDir()
	// Index references "overview" + "made-up" — made-up.json won't
	// exist so LoadPlan will fail on the second file open.
	index := `{
		"title": "broken plan",
		"sections": ["overview", "made-up"]
	}`
	if err := os.WriteFile(filepath.Join(planDir, "index.json"), []byte(index), 0644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	overview := `{"items": [{"id": "OVR-1", "title": "ok", "desc": "ok"}]}`
	if err := os.WriteFile(filepath.Join(planDir, "overview.json"), []byte(overview), 0644); err != nil {
		t.Fatalf("write overview: %v", err)
	}

	cmd := exec.Command(bin, "validate", filepath.Join(planDir, "index.json"))
	combined, _ := cmd.CombinedOutput()
	if exit := cmd.ProcessState.ExitCode(); exit != 1 {
		t.Fatalf("exit = %d, want 1\n%s", exit, combined)
	}
	out := string(combined)
	if !strings.Contains(out, "made-up") {
		t.Errorf("expected error to mention missing made-up section\nfull output:\n%s", out)
	}
	if !strings.Contains(out, "invalid (") {
		t.Errorf("missing invalid summary line\nfull output:\n%s", out)
	}
}

// TestMain_Validate_JSON asserts --json emits the report struct shape.
func TestMain_Validate_JSON(t *testing.T) {
	bin := buildBinary(t)
	root := repoRoot(t)
	input := filepath.Join(root, "testdata", "plan-example", "index.json")
	if _, err := os.Stat(input); err != nil {
		t.Fatalf("plan fixture missing at %s: %v", input, err)
	}

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
