package renderer

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// -update regenerates the golden file instead of comparing against it.
// Standard go test idiom: `go test ./pkg/renderer/ -run TestRender_Golden -update`.
var update = flag.Bool("update", false, "update golden files")

// goldenPath is the on-disk location of the expected rendered HTML.
const goldenPath = "testdata/render_golden_smoke.html"

// smokeInputPath is the comprehensive fixture at the repo root that
// exercises every art form (raw SVG, wf:*, every DSL primitive, auto, none).
const smokeInputPath = "../../testdata/smoke.json"

// TestRender_Golden is the end-to-end snapshot test for Render. It pins:
//   - the embedded template + CSS
//   - every art-form path (raw SVG, named library, DSL, auto, none)
//   - the SOURCE.file constant embedded from RenderOpts.OutputPath
//   - the theme attribute (forced to "paper" so this never depends on cache state)
//
// Any time the renderer's output changes intentionally, regenerate with:
//   go test ./pkg/renderer/ -run TestRender_Golden -update
func TestRender_Golden(t *testing.T) {
	// Belt-and-braces: redirect HOME/USERPROFILE so even if ThemeOverride
	// were ever bypassed by a future bug, ResolveTheme would write to the
	// temp dir, not the real ~/.squiz/themes.json.
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	raw, err := os.ReadFile(smokeInputPath)
	if err != nil {
		t.Fatalf("read smoke fixture %s: %v", smokeInputPath, err)
	}

	doc, err := ParseDocument(raw)
	if err != nil {
		t.Fatalf("ParseDocument: %v", err)
	}

	// OutputPath is embedded in the page (SOURCE.file). Hard-code it so the
	// golden is stable across machines and run directories.
	// ThemeOverride: "paper" makes ResolveTheme short-circuit before it
	// ever touches the on-disk theme cache.
	got, err := Render(doc, RenderOpts{
		OutputPath:    "/test/smoke.html",
		ThemeOverride: "paper",
		WorkDir:       "/test",
	})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatalf("mkdir golden dir: %v", err)
		}
		if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		t.Logf("updated golden: %s (%d bytes)", goldenPath, len(got))
		return
	}

	wantBytes, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden %s: %v (run with -update to generate)", goldenPath, err)
	}
	want := string(wantBytes)

	if got == want {
		return
	}

	// Compact diff: byte counts + small window around the first divergence.
	// Don't dump the whole 65KB HTML.
	idx := firstDiff(got, want)
	t.Fatalf(
		"rendered output does not match golden\n"+
			"  golden: %s\n"+
			"  got=%d bytes, want=%d bytes, first diff at byte %d\n"+
			"  got  [%d:%d]: %q\n"+
			"  want [%d:%d]: %q\n"+
			"if intentional, regenerate with: go test ./pkg/renderer/ -run TestRender_Golden -update",
		goldenPath,
		len(got), len(want), idx,
		idx, clamp(idx+200, len(got)), excerpt(got, idx, 200),
		idx, clamp(idx+200, len(want)), excerpt(want, idx, 200),
	)
}

// TestRender_RespectsOverride asserts the ThemeOverride field actually
// influences the rendered output (smoke check that override plumbing works).
func TestRender_RespectsOverride(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	doc, err := ParseDocument([]byte(`{"squizzes": []}`))
	if err != nil {
		t.Fatalf("ParseDocument: %v", err)
	}
	got, err := Render(doc, RenderOpts{
		OutputPath:    "/test/x.html",
		ThemeOverride: "phosphor",
		WorkDir:       "/test",
	})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if !strings.Contains(got, `data-theme="phosphor"`) {
		t.Errorf(`output missing data-theme="phosphor"`)
	}
}

// firstDiff returns the index of the first byte where a and b differ,
// or min(len(a), len(b)) if one is a prefix of the other.
func firstDiff(a, b string) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}

func clamp(i, max int) int {
	if i > max {
		return max
	}
	return i
}

func excerpt(s string, start, n int) string {
	if start >= len(s) {
		return ""
	}
	end := start + n
	if end > len(s) {
		end = len(s)
	}
	return s[start:end]
}
