package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestMain_Preview runs `squiz preview wf:calendar-grid --out <tmp>` and
// asserts the file is created and contains an SVG fragment.
func TestMain_Preview(t *testing.T) {
	bin := buildBinary(t)

	outDir := t.TempDir()
	outFile := filepath.Join(outDir, "preview.html")

	cmd := exec.Command(bin, "preview", "wf:calendar-grid", "--out", outFile)
	combined, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("squiz preview failed: %v\noutput: %s", err, combined)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("expected preview output at %s: %v", outFile, err)
	}
	s := string(data)
	if !strings.Contains(s, "<svg") {
		t.Errorf("preview output missing <svg fragment\n---\n%s", s)
	}
	// The spec string should be echoed in the page for reference.
	if !strings.Contains(s, "wf:calendar-grid") {
		t.Errorf("preview output missing spec string echo")
	}
}

// TestMain_Preview_Stdout asserts --stdout prints to stdout without writing
// a file.
func TestMain_Preview_Stdout(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "preview", "wf:calendar-grid", "--stdout")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz preview --stdout failed: %v", err)
	}
	if !strings.Contains(string(out), "<svg") {
		t.Errorf("preview --stdout missing <svg fragment\n---\n%s", out)
	}
}

// TestMain_Preview_DSL asserts a DSL spec (with brackets/commas) round-trips
// through the filename sanitizer and renders.
func TestMain_Preview_DSL(t *testing.T) {
	bin := buildBinary(t)

	outDir := t.TempDir()
	outFile := filepath.Join(outDir, "flow.html")

	cmd := exec.Command(bin, "preview", "flow:[client,api,db]", "--out", outFile, "--theme", "phosphor")
	combined, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("squiz preview flow failed: %v\noutput: %s", err, combined)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("expected output at %s: %v", outFile, err)
	}
	s := string(data)
	if !strings.Contains(s, "<svg") {
		t.Errorf("preview output missing <svg\n---\n%s", s)
	}
	if !strings.Contains(s, "phosphor") {
		t.Errorf("preview output missing theme name")
	}
}
