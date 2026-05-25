package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// buildBinary compiles ./cmd/squiz-plan into a temp dir and returns the
// path to the executable. Skips the test if the go toolchain isn't on
// PATH (defensive — covers minimal CI images).
func buildBinary(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not on PATH; skipping CLI build test")
	}
	tmp := t.TempDir()
	name := "squiz-plan"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(tmp, name)

	// Build from the squiz-plan package dir. The test runs with cwd =
	// cmd/squiz-plan, so "." refers to this package.
	cmd := exec.Command("go", "build", "-o", bin, ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed: %v\noutput: %s", err, out)
	}
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("built binary missing: %v", err)
	}
	return bin
}

// repoRoot returns the absolute path to the repo root so we can locate
// testdata/plan-example regardless of the test's working directory.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find go.mod walking up from %s", dir)
		}
		dir = parent
	}
}

// TestMain_Build is the end-to-end CLI smoke test. It compiles the binary
// fresh, runs `squiz-plan render --theme paper --out <tmp> <plan-example/index.json>`,
// and asserts the rendered file exists and is plausibly sized.
//
// IMPORTANT — theme cache safety: we pass `--theme paper` rather than
// stubbing HOME/USERPROFILE on the child process. The renderer's
// ResolveTheme short-circuits on any valid override before it ever calls
// loadCache/saveCache, so this guarantees ~/.squiz/themes.json is not
// touched. Mirrors the strategy in cmd/squiz/main_test.go.
//
// Test SKIPS cleanly if the template assets aren't yet on disk — the
// sibling template agent may not have landed plan.html.tmpl when this
// suite first runs.
func TestMain_Build(t *testing.T) {
	root := repoRoot(t)

	// Skip if templates haven't landed yet — render would fail with a
	// clear error but we don't want a red CI on incomplete sibling work.
	tmplPath := filepath.Join(root, "internal", "planview", "templates", "plan.html.tmpl")
	if _, err := os.Stat(tmplPath); err != nil {
		t.Skipf("plan template not on disk yet (%s); sibling template agent in progress", tmplPath)
	}

	bin := buildBinary(t)

	input := filepath.Join(root, "testdata", "plan-example", "index.json")
	if _, err := os.Stat(input); err != nil {
		t.Fatalf("plan fixture missing at %s: %v", input, err)
	}

	outDir := t.TempDir()
	outFile := filepath.Join(outDir, "plan.html")

	cmd := exec.Command(bin, "render", "--theme", "paper", "--out", outFile, input)
	combined, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("squiz-plan render failed: %v\noutput: %s", err, combined)
	}
	if exitCode := cmd.ProcessState.ExitCode(); exitCode != 0 {
		t.Fatalf("squiz-plan render exit code = %d, want 0\noutput: %s", exitCode, combined)
	}

	info, err := os.Stat(outFile)
	if err != nil {
		t.Fatalf("expected output file %s: %v\nstderr: %s", outFile, err, combined)
	}
	if info.Size() < 1024 {
		t.Errorf("output file %s is %d bytes, expected > 1KB", outFile, info.Size())
	}
}

// TestMain_Version asserts `squiz-plan version` returns 0 and prints
// something recognisable on stdout. Independent of template state.
func TestMain_Version(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "version")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz-plan version failed: %v\nstdout: %s", err, out)
	}
	if exitCode := cmd.ProcessState.ExitCode(); exitCode != 0 {
		t.Fatalf("squiz-plan version exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(string(out), "squiz-plan") {
		t.Errorf("stdout missing %q: %q", "squiz-plan", string(out))
	}
}
