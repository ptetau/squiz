package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestEmbeddedExampleMatchesRoot guards against drift between the
// embedded sample (used by `squiz example`) and the canonical fixture
// at testdata/smoke.json that the rest of the repo uses for goldens
// and CLI integration tests.
func TestEmbeddedExampleMatchesRoot(t *testing.T) {
	root := repoRoot(t)
	onDisk, err := os.ReadFile(filepath.Join(root, "testdata", "smoke.json"))
	if err != nil {
		t.Fatalf("read root testdata/smoke.json: %v", err)
	}
	if !bytes.Equal(embeddedExample, onDisk) {
		t.Fatalf("embedded example out of sync with testdata/smoke.json\n  embedded=%d bytes, on-disk=%d bytes\n  fix: cp testdata/smoke.json cmd/squiz/example/smoke.json",
			len(embeddedExample), len(onDisk))
	}
}

// TestMain_Example runs `squiz example --out tmp.json` and asserts the
// written file is byte-identical to the embedded sample.
func TestMain_Example(t *testing.T) {
	bin := buildBinary(t)
	outFile := filepath.Join(t.TempDir(), "scaffolded.json")

	cmd := exec.Command(bin, "example", "--out", outFile)
	if combined, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("squiz example failed: %v\noutput: %s", err, combined)
	}
	got, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read scaffolded file: %v", err)
	}
	if !bytes.Equal(got, embeddedExample) {
		t.Errorf("scaffolded content differs from embedded (sizes: got=%d, want=%d)", len(got), len(embeddedExample))
	}
}

// TestMain_ExampleStdout asserts `squiz example --stdout` writes the
// sample to stdout (no file), exit 0.
func TestMain_ExampleStdout(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "example", "--stdout")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz example --stdout failed: %v", err)
	}
	if !bytes.Equal(out, embeddedExample) {
		t.Errorf("stdout differs from embedded (got=%d, want=%d)", len(out), len(embeddedExample))
	}
}
