package renderer

import (
	"strings"
	"testing"
)

// assertSVG checks the output is a non-empty SVG fragment, and (if any
// substrings are provided) that each appears in it.
func assertSVG(t *testing.T, label, svg string, subs ...string) {
	t.Helper()
	if svg == "" {
		t.Fatalf("%s: expected non-empty SVG, got empty", label)
	}
	if !strings.Contains(svg, "<svg") {
		t.Fatalf("%s: expected SVG fragment, got: %s", label, svg)
	}
	for _, s := range subs {
		if !strings.Contains(svg, s) {
			t.Errorf("%s: expected substring %q in output\n%s", label, s, svg)
		}
	}
}

// ──────────────────────────────────────────────────────────────────────
// grid
// ──────────────────────────────────────────────────────────────────────

func TestDslGrid_Valid(t *testing.T) {
	out := dslGrid("3x3@0.5")
	assertSVG(t, "grid valid", out, "<rect")
}

func TestDslGrid_DefaultRate(t *testing.T) {
	// no @rate suffix → default 0.5
	out := dslGrid("4x4")
	assertSVG(t, "grid default rate", out, "<rect")
}

func TestDslGrid_BadDims(t *testing.T) {
	out := dslGrid("foo")
	assertSVG(t, "grid no-x", out, "grid")

	out2 := dslGrid("3xfoo")
	assertSVG(t, "grid bad cols", out2, "grid")

	out3 := dslGrid("999x999") // dims out of range
	assertSVG(t, "grid oversize", out3, "grid")
}

// ──────────────────────────────────────────────────────────────────────
// spark
// ──────────────────────────────────────────────────────────────────────

func TestDslSpark_Valid(t *testing.T) {
	out := dslSpark("[1,2,3,4]")
	assertSVG(t, "spark valid", out, "<path")
}

func TestDslSpark_FlatLine(t *testing.T) {
	// All-equal values: rng=0 → forced to 1, no NaNs.
	out := dslSpark("[5,5,5,5]")
	assertSVG(t, "spark flat", out, "<path")
	if strings.Contains(out, "NaN") {
		t.Errorf("spark flat produced NaN in output: %s", out)
	}
}

func TestDslSpark_Invalid(t *testing.T) {
	out := dslSpark("not-a-list")
	assertSVG(t, "spark invalid", out, "spark")

	out2 := dslSpark("[1]") // need ≥ 2 points
	assertSVG(t, "spark too short", out2, "spark")
}

// ──────────────────────────────────────────────────────────────────────
// bars
// ──────────────────────────────────────────────────────────────────────

func TestDslBars_Valid(t *testing.T) {
	out := dslBars("[3,5,4,7,2]")
	assertSVG(t, "bars valid", out, "<rect")
}

func TestDslBars_AllZero(t *testing.T) {
	// max=0 → forced to 1, all bars height 0, no NaN.
	out := dslBars("[0,0,0]")
	assertSVG(t, "bars zero", out, "<rect")
	if strings.Contains(out, "NaN") {
		t.Errorf("bars zero produced NaN: %s", out)
	}
}

func TestDslBars_Invalid(t *testing.T) {
	out := dslBars("nope")
	assertSVG(t, "bars invalid", out, "bars")
}

// ──────────────────────────────────────────────────────────────────────
// swatches
// ──────────────────────────────────────────────────────────────────────

func TestDslSwatches_Valid(t *testing.T) {
	out := dslSwatches("#fff,#000,#aaa")
	assertSVG(t, "swatches valid", out, "#fff")
	if !strings.Contains(out, "#000") {
		t.Errorf("missing second color in swatches output:\n%s", out)
	}
}

func TestDslSwatches_SingleColor(t *testing.T) {
	out := dslSwatches("#abc")
	assertSVG(t, "swatches single", out, "#abc")
}

// Note: dslSwatches currently has no explicit failure path — strings.Split
// always returns ≥1 element so the len==0 check can't fire. Empty input
// produces an SVG with no <rect> elements but is still a valid frame.
func TestDslSwatches_EmptyArg(t *testing.T) {
	out := dslSwatches("")
	assertSVG(t, "swatches empty", out)
}

// ──────────────────────────────────────────────────────────────────────
// pills
// ──────────────────────────────────────────────────────────────────────

func TestDslPills_Valid(t *testing.T) {
	out := dslPills("one*|two|three*")
	assertSVG(t, "pills valid", out, "one")
	if !strings.Contains(out, "three") {
		t.Errorf("missing 'three' in pills output:\n%s", out)
	}
}

func TestDslPills_SingleChip(t *testing.T) {
	out := dslPills("solo")
	assertSVG(t, "pills single", out, "solo")
}

// Empty arg: dslPills also has no real error path because strings.Split
// returns [""]; locks current behavior (renders an empty-text pill).
func TestDslPills_EmptyArg(t *testing.T) {
	out := dslPills("")
	assertSVG(t, "pills empty", out)
}

// ──────────────────────────────────────────────────────────────────────
// sample
// ──────────────────────────────────────────────────────────────────────

func TestDslSample_Valid(t *testing.T) {
	out := dslSample(`"hello"@mono`)
	assertSVG(t, "sample valid", out, "hello")
	if !strings.Contains(out, "IBM Plex Mono") {
		t.Errorf("expected mono font family, got:\n%s", out)
	}
}

func TestDslSample_DefaultFont(t *testing.T) {
	// No @font → defaults to sans (IBM Plex Sans).
	out := dslSample(`"howdy"`)
	assertSVG(t, "sample default", out, "howdy")
	if !strings.Contains(out, "IBM Plex Sans") {
		t.Errorf("expected sans font family, got:\n%s", out)
	}
}

func TestDslSample_SerifFont(t *testing.T) {
	out := dslSample(`"editorial"@serif`)
	assertSVG(t, "sample serif", out, "editorial")
	if !strings.Contains(out, "IBM Plex Serif") {
		t.Errorf("expected serif font family, got:\n%s", out)
	}
}

func TestDslSample_Invalid(t *testing.T) {
	// No surrounding quotes → error placeholder.
	out := dslSample(`bare-text`)
	assertSVG(t, "sample invalid", out, "sample")
}

// ──────────────────────────────────────────────────────────────────────
// circle-pack
// ──────────────────────────────────────────────────────────────────────

func TestDslCirclePack_Valid(t *testing.T) {
	out := dslCirclePack("8")
	assertSVG(t, "circle-pack valid", out, "<circle")
}

func TestDslCirclePack_Invalid(t *testing.T) {
	out := dslCirclePack("not-a-number")
	assertSVG(t, "circle-pack NaN", out, "circle-pack")

	out2 := dslCirclePack("0") // n<=0 rejected
	assertSVG(t, "circle-pack zero", out2, "circle-pack")

	out3 := dslCirclePack("999") // n>50 rejected
	assertSVG(t, "circle-pack oversize", out3, "circle-pack")
}

// ──────────────────────────────────────────────────────────────────────
// resolveDSL dispatcher
// ──────────────────────────────────────────────────────────────────────

func TestResolveDSL_UnknownPrefix(t *testing.T) {
	svg, hidden := resolveDSL("nosuch:foo")
	if hidden {
		t.Errorf("hidden = true, want false for unknown dsl")
	}
	assertSVG(t, "resolveDSL unknown", svg, "unknown dsl")
}

// ──────────────────────────────────────────────────────────────────────
// helpers
// ──────────────────────────────────────────────────────────────────────

func TestParseNumList(t *testing.T) {
	cases := []struct {
		in     string
		want   []float64
		wantOk bool
	}{
		{"1,2,3", []float64{1, 2, 3}, true},
		{"[1,2,3]", []float64{1, 2, 3}, true},
		{"[ 1.5 , 2.5 ]", []float64{1.5, 2.5}, true},
		{"[1,,2]", []float64{1, 2}, true}, // empty parts skipped
		{"foo", nil, false},
		{"1,foo,3", nil, false},
		{"", nil, false}, // empty → no values → !ok
		{"[]", nil, false},
	}
	for _, tc := range cases {
		got, ok := parseNumList(tc.in)
		if ok != tc.wantOk {
			t.Errorf("parseNumList(%q) ok = %v, want %v", tc.in, ok, tc.wantOk)
			continue
		}
		if !tc.wantOk {
			continue
		}
		if len(got) != len(tc.want) {
			t.Errorf("parseNumList(%q) len = %d, want %d", tc.in, len(got), len(tc.want))
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("parseNumList(%q)[%d] = %v, want %v", tc.in, i, got[i], tc.want[i])
			}
		}
	}
}

func TestEscapeXML(t *testing.T) {
	in := `a<b>c&d"e'f`
	got := escapeXML(in)
	for _, want := range []string{"&lt;", "&gt;", "&amp;", "&quot;", "&apos;"} {
		if !strings.Contains(got, want) {
			t.Errorf("escapeXML missing %q in output: %s", want, got)
		}
	}
	for _, bad := range []string{"<", ">", `"`, `'`} {
		if strings.Contains(got, bad) {
			t.Errorf("escapeXML left raw %q in output: %s", bad, got)
		}
	}
}

func TestErrArt(t *testing.T) {
	out := errArt("boom")
	assertSVG(t, "errArt", out, "boom")
}
