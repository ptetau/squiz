package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// buildBinary compiles ./cmd/squiz into a temp dir and returns the path to
// the executable. Skips the test if the go toolchain isn't on PATH (defensive
// — covers minimal CI images).
func buildBinary(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go not on PATH; skipping CLI build test")
	}
	tmp := t.TempDir()
	name := "squiz"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	bin := filepath.Join(tmp, name)

	// Build from the module root (parent dir of this test file's dir).
	// The test runs with cwd = cmd/squiz, so ./ refers to the squiz package.
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
// testdata/smoke.json regardless of the test's working directory.
func repoRoot(t *testing.T) string {
	t.Helper()
	// Walk up from cwd until we find go.mod.
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
// fresh, runs `squiz render --theme paper --out <tmp> <smoke.json>`, and
// asserts the rendered file exists and is plausibly sized.
//
// IMPORTANT — theme cache safety: we pass `--theme paper` rather than
// stubbing HOME/USERPROFILE on the child process. The renderer's
// ResolveTheme short-circuits on any valid override before it ever calls
// loadCache/saveCache, so this guarantees ~/.squiz/themes.json is not
// touched. (Setting HOME would also work but relies on os.UserHomeDir
// honouring the override, which is one more layer to trust.)
//
// IMPORTANT — flag ordering: Go's flag.Parse stops at the first non-flag
// argument. The CLI calls fs.Parse(args), so flags MUST precede the
// positional <input.json>; otherwise --out is silently dropped and the
// CLI overwrites testdata/smoke.html in the repo.
func TestMain_Build(t *testing.T) {
	bin := buildBinary(t)
	root := repoRoot(t)

	input := filepath.Join(root, "testdata", "smoke.json")
	if _, err := os.Stat(input); err != nil {
		t.Fatalf("smoke fixture missing at %s: %v", input, err)
	}

	outDir := t.TempDir()
	outFile := filepath.Join(outDir, "smoke.html")

	cmd := exec.Command(bin, "render", "--theme", "paper", "--out", outFile, input)
	combined, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("squiz render failed: %v\noutput: %s", err, combined)
	}
	if exitCode := cmd.ProcessState.ExitCode(); exitCode != 0 {
		t.Fatalf("squiz render exit code = %d, want 0\noutput: %s", exitCode, combined)
	}

	info, err := os.Stat(outFile)
	if err != nil {
		t.Fatalf("expected output file %s: %v\nstderr: %s", outFile, err, combined)
	}
	if info.Size() < 1024 {
		t.Errorf("output file %s is %d bytes, expected > 1KB", outFile, info.Size())
	}
}

// TestMain_Version asserts `squiz version` returns 0 and prints something
// recognisable on stdout.
func TestMain_Version(t *testing.T) {
	bin := buildBinary(t)

	cmd := exec.Command(bin, "version")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("squiz version failed: %v\nstdout: %s", err, out)
	}
	if exitCode := cmd.ProcessState.ExitCode(); exitCode != 0 {
		t.Fatalf("squiz version exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(string(out), "squiz") {
		t.Errorf("stdout missing %q: %q", "squiz", string(out))
	}
}
