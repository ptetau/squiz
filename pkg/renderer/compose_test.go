package renderer

import (
	"strings"
	"testing"
)

// TestResolveUses_NoUses confirms the common case (raw SVG with no
// composition refs) is a byte-identical pass-through — this is what keeps
// the golden render test stable.
func TestResolveUses_NoUses(t *testing.T) {
	in := `<svg viewBox='0 0 100 60'><rect x='10' y='10' width='80' height='40' fill='var(--accent)'/></svg>`
	got := resolveUses(in)
	if got != in {
		t.Errorf("expected pass-through for no-use input\n  in:  %s\n  got: %s", in, got)
	}
}

// TestResolveUses_SingleWF verifies the canonical happy path: one wf:
// reference is inlined as a <symbol> + the <use> rewritten to a local id.
func TestResolveUses_SingleWF(t *testing.T) {
	in := `<svg viewBox='0 0 100 60'><use href="wf:phone-card"/></svg>`
	got := resolveUses(in)

	if !strings.Contains(got, `<defs>`) {
		t.Errorf("expected <defs> block, got:\n%s", got)
	}
	if !strings.Contains(got, `<symbol id="wf-phone-card"`) {
		t.Errorf(`expected <symbol id="wf-phone-card", got:\n%s`, got)
	}
	if !strings.Contains(got, `viewBox="0 0 100 60"`) {
		t.Errorf("expected viewBox on symbol, got:\n%s", got)
	}
	if !strings.Contains(got, `<use href="#wf-phone-card"`) {
		t.Errorf(`expected rewritten <use href="#wf-phone-card", got:\n%s`, got)
	}
	if strings.Contains(got, `href="wf:phone-card"`) {
		t.Errorf("original wf: href should be gone, still present in:\n%s", got)
	}
	// defs must come right after the opening <svg ...> tag.
	if !strings.Contains(got, `60'><defs>`) && !strings.Contains(got, `60"><defs>`) {
		t.Errorf("expected <defs> spliced after opening <svg>, got:\n%s", got)
	}
}

// TestResolveUses_MultipleRefs mixes wf:, arch:, and DSL refs and asserts
// each gets its own <symbol> + each <use> is rewritten.
func TestResolveUses_MultipleRefs(t *testing.T) {
	in := `<svg viewBox='0 0 100 60'>` +
		`<use href="wf:phone-card" x="0" y="0" width="30" height="50"/>` +
		`<use href="arch:database" x="35" y="0" width="30" height="50"/>` +
		`<use href="bars:[1,2,3]" x="70" y="0" width="30" height="50"/>` +
		`</svg>`
	got := resolveUses(in)

	for _, id := range []string{"wf-phone-card", "arch-database", "bars-_1_2_3_"} {
		if !strings.Contains(got, `<symbol id="`+id+`"`) {
			t.Errorf(`expected <symbol id="%s">, got:\n%s`, id, got)
		}
		if !strings.Contains(got, `href="#`+id+`"`) {
			t.Errorf(`expected rewritten href="#%s", got:\n%s`, id, got)
		}
	}
}

// TestResolveUses_DuplicateRef — two uses of the same spec → exactly one
// <symbol> definition.
func TestResolveUses_DuplicateRef(t *testing.T) {
	in := `<svg viewBox='0 0 100 60'>` +
		`<use href="wf:phone-card" x="0" y="0" width="30" height="50"/>` +
		`<use href="wf:phone-card" x="60" y="0" width="30" height="50"/>` +
		`</svg>`
	got := resolveUses(in)

	count := strings.Count(got, `<symbol id="wf-phone-card"`)
	if count != 1 {
		t.Errorf("expected exactly 1 <symbol id=\"wf-phone-card\">, got %d\noutput:\n%s", count, got)
	}
	useCount := strings.Count(got, `href="#wf-phone-card"`)
	if useCount != 2 {
		t.Errorf("expected 2 rewritten <use>s, got %d\noutput:\n%s", useCount, got)
	}
}

// TestResolveUses_UnknownRef — an unknown wf:* name must still produce a
// <symbol> so the <use> doesn't dangle; the body comes from the same
// fail-soft placeholder resolveNamed emits standalone.
func TestResolveUses_UnknownRef(t *testing.T) {
	in := `<svg viewBox='0 0 100 60'><use href="wf:nonsense-thing"/></svg>`
	got := resolveUses(in)

	if !strings.Contains(got, `<symbol id="wf-nonsense-thing"`) {
		t.Errorf("expected placeholder symbol, got:\n%s", got)
	}
	if !strings.Contains(got, `href="#wf-nonsense-thing"`) {
		t.Errorf("expected rewritten href, got:\n%s", got)
	}
	// The placeholder body should still mention the bad name so the author
	// sees their typo in the rendered output.
	if !strings.Contains(got, "nonsense-thing") {
		t.Errorf("expected placeholder body to surface the bad name, got:\n%s", got)
	}
}

// TestResolveUses_PreservesAttrs — every non-href attribute on the
// original <use> must round-trip unchanged (x, y, width, height,
// transform, …).
func TestResolveUses_PreservesAttrs(t *testing.T) {
	in := `<svg viewBox='0 0 100 60'><use href="wf:phone-card" x="10" y="5" width="30" height="50" transform="rotate(10)"/></svg>`
	got := resolveUses(in)

	for _, want := range []string{`x="10"`, `y="5"`, `width="30"`, `height="50"`, `transform="rotate(10)"`} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %s preserved in rewritten use, got:\n%s", want, got)
		}
	}
	if !strings.Contains(got, `href="#wf-phone-card"`) {
		t.Errorf("href not rewritten correctly, got:\n%s", got)
	}
}

// TestResolveUses_IgnoresLocalRefs — a `<use href="#local-defs"/>` (no
// colon-prefix art spec) is left completely untouched and no <defs>
// block is added.
func TestResolveUses_IgnoresLocalRefs(t *testing.T) {
	in := `<svg viewBox='0 0 100 60'><defs><symbol id="my-thing"><rect width="10" height="10"/></symbol></defs><use href="#my-thing" x="5" y="5"/></svg>`
	got := resolveUses(in)
	if got != in {
		t.Errorf("expected local-href input to be untouched\n  in:  %s\n  got: %s", in, got)
	}
}

// TestResolveUses_SingleQuotedHref — agents may quote the href with '
// instead of " (matching the codebase's own style). Both must work.
func TestResolveUses_SingleQuotedHref(t *testing.T) {
	in := `<svg viewBox='0 0 100 60'><use href='wf:phone-card'/></svg>`
	got := resolveUses(in)
	if !strings.Contains(got, `<symbol id="wf-phone-card"`) {
		t.Errorf("single-quoted href not resolved, got:\n%s", got)
	}
	if !strings.Contains(got, `href="#wf-phone-card"`) {
		// Note: rewriteHref preserves the quote style of the *attribute
		// value*; if it was single-quoted originally, we'd expect a
		// single-quoted rewrite. Accept either.
		if !strings.Contains(got, `href='#wf-phone-card'`) {
			t.Errorf("href not rewritten (single or double), got:\n%s", got)
		}
	}
}

// TestResolveUses_MixedComposableAndLocal — when a single SVG has BOTH
// composable refs and plain local-fragment refs, only the composable
// ones are touched.
func TestResolveUses_MixedComposableAndLocal(t *testing.T) {
	in := `<svg viewBox='0 0 100 60'>` +
		`<defs><symbol id="local"><circle r="2"/></symbol></defs>` +
		`<use href="#local" x="5"/>` +
		`<use href="wf:phone-card" x="40"/>` +
		`</svg>`
	got := resolveUses(in)
	if !strings.Contains(got, `href="#local"`) {
		t.Errorf("local href should be untouched, got:\n%s", got)
	}
	if !strings.Contains(got, `href="#wf-phone-card"`) {
		t.Errorf("wf: href should be rewritten, got:\n%s", got)
	}
}
