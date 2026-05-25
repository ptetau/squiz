package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// planFiles is the list the embedded fixture MUST contain, matching the
// canonical testdata/plan-example/ tree byte-for-byte.
var planFiles = []string{
	"index.json",
	"overview.json",
	"functional.json",
	"non-functional.json",
	"cases.json",
	"engineering.json",
	"build.json",
}

// TestEmbeddedPlanMatchesRoot guards against drift between the embedded
// plan tree (used by `squiz-plan example`) and the canonical fixture at
// testdata/plan-example/ used for goldens and CLI integration tests.
func TestEmbeddedPlanMatchesRoot(t *testing.T) {
	root := repoRoot(t)
	for _, name := range planFiles {
		t.Run(name, func(t *testing.T) {
			onDisk, err := os.ReadFile(filepath.Join(root, "testdata", "plan-example", name))
			if err != nil {
				t.Fatalf("read root testdata/plan-example/%s: %v", name, err)
			}
			embedded, err := embeddedPlan.ReadFile("example/" + name)
			if err != nil {
				t.Fatalf("read embedded example/%s: %v", name, err)
			}
			if !bytes.Equal(embedded, onDisk) {
				t.Fatalf("embedded %s out of sync with testdata copy (embedded=%d bytes, on-disk=%d bytes)\n  fix: cp testdata/plan-example/%s cmd/squiz-plan/example/%s",
					name, len(embedded), len(onDisk), name, name)
			}
		})
	}
}

// TestMain_Example runs `squiz-plan example --out tmp/` and asserts the
// 7 written files exist + match the embedded copies.
func TestMain_Example(t *testing.T) {
	bin := buildBinary(t)
	outDir := filepath.Join(t.TempDir(), "scaffolded")

	cmd := exec.Command(bin, "example", "--out", outDir)
	if combined, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("squiz-plan example failed: %v\noutput: %s", err, combined)
	}

	for _, name := range planFiles {
		got, err := os.ReadFile(filepath.Join(outDir, name))
		if err != nil {
			t.Errorf("scaffolded %s missing: %v", name, err)
			continue
		}
		embedded, _ := embeddedPlan.ReadFile("example/" + name)
		if !bytes.Equal(got, embedded) {
			t.Errorf("scaffolded %s differs from embedded", name)
		}
	}
}
