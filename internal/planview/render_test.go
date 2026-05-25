package planview

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// -update regenerates the golden file instead of comparing against it.
// Standard go test idiom: `go test ./internal/planview/ -run TestRender_PlanExample -update`.
var update = flag.Bool("update", false, "update golden files")

// goldenPlanPath is the on-disk location of the expected rendered HTML.
const goldenPlanPath = "testdata/render_golden_plan.html"

// planExampleIndex is the worked plan-tree fixture at the repo root.
const planExampleIndex = "../../testdata/plan-example/index.json"

// templatesReady reports whether the sibling template agent has landed
// plan.html.tmpl yet. Tests that depend on a successful render Skip
// when this returns false so an in-progress sibling agent doesn't
// turn the suite red.
func templatesReady() bool {
	_, err := os.Stat("templates/plan.html.tmpl")
	return err == nil
}

// TestRender_PlanExample is the end-to-end snapshot test for plan
// rendering. It pins:
//   - the embedded template + CSS (shared + plan-specific)
//   - per-item art resolution (raw, wf:*, DSL, "none", omitted)
//   - the SOURCE.file constant embedded from RenderOpts.OutputPath
//     (if the template surfaces it — the template agent's call)
//   - the theme attribute (forced to "paper" so this never depends on
//     cache state)
//
// Any time the renderer's output changes intentionally, regenerate with:
//
//	go test ./internal/planview/ -run TestRender_PlanExample -update
func TestRender_PlanExample(t *testing.T) {
	if !templatesReady() {
		t.Skip("templates/plan.html.tmpl not on disk yet; sibling template agent in progress")
	}

	// Belt-and-braces: redirect HOME/USERPROFILE so even if ThemeOverride
	// were ever bypassed by a future bug, ResolveTheme would write to the
	// temp dir, not the real ~/.squiz/themes.json. The brief mandates
	// BOTH this and --theme paper (defence in depth, not either/or).
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	plan, err := LoadPlan(planExampleIndex)
	if err != nil {
		t.Fatalf("LoadPlan: %v", err)
	}

	got, err := Render(plan, RenderOpts{
		OutputPath:    "/test/plan.html",
		ThemeOverride: "paper",
		WorkDir:       "/test",
	})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenPlanPath), 0o755); err != nil {
			t.Fatalf("mkdir golden dir: %v", err)
		}
		if err := os.WriteFile(goldenPlanPath, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		t.Logf("updated golden: %s (%d bytes)", goldenPlanPath, len(got))
		return
	}

	wantBytes, err := os.ReadFile(goldenPlanPath)
	if err != nil {
		t.Fatalf("read golden %s: %v (run with -update to generate)", goldenPlanPath, err)
	}
	want := string(wantBytes)

	if got == want {
		return
	}

	idx := firstDiff(got, want)
	t.Fatalf(
		"rendered output does not match golden\n"+
			"  golden: %s\n"+
			"  got=%d bytes, want=%d bytes, first diff at byte %d\n"+
			"  got  [%d:%d]: %q\n"+
			"  want [%d:%d]: %q\n"+
			"if intentional, regenerate with: go test ./internal/planview/ -run TestRender_PlanExample -update",
		goldenPlanPath,
		len(got), len(want), idx,
		idx, clamp(idx+200, len(got)), excerpt(got, idx, 200),
		idx, clamp(idx+200, len(want)), excerpt(want, idx, 200),
	)
}

// TestRender_ResolvesRefs verifies that an item's Refs are expanded to
// "<SectionLabel> · <itemID>" labels in the rendered HTML. ENG-1
// references FR-1 (functional) and NFR-1 (non-functional), so both
// formatted labels should appear in the output.
func TestRender_ResolvesRefs(t *testing.T) {
	if !templatesReady() {
		t.Skip("templates/plan.html.tmpl not on disk yet; sibling template agent in progress")
	}
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	plan, err := LoadPlan(planExampleIndex)
	if err != nil {
		t.Fatalf("LoadPlan: %v", err)
	}
	got, err := Render(plan, RenderOpts{
		OutputPath:    "/test/plan.html",
		ThemeOverride: "paper",
		WorkDir:       "/test",
	})
	if err != nil {
		t.Fatalf("Render: %v", err)
	}

	// Middle dot is U+00B7 — same character emitted by RefView.Label.
	wantLabels := []string{
		"Functional · FR-1",
		"Non-functional · NFR-1",
	}
	for _, lbl := range wantLabels {
		if !strings.Contains(got, lbl) {
			t.Errorf("rendered HTML missing ref label %q", lbl)
		}
	}
}

// TestRender_RespectsThemeOverride asserts the ThemeOverride field
// actually influences the rendered output. Smoke check that the override
// plumbing reaches the template's data-theme attribute.
func TestRender_RespectsThemeOverride(t *testing.T) {
	if !templatesReady() {
		t.Skip("templates/plan.html.tmpl not on disk yet; sibling template agent in progress")
	}
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	plan, err := LoadPlan(planExampleIndex)
	if err != nil {
		t.Fatalf("LoadPlan: %v", err)
	}
	got, err := Render(plan, RenderOpts{
		OutputPath:    "/test/plan.html",
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
