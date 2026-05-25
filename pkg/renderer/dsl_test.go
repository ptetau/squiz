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

// Empty input must emit the errArt placeholder (consistent with other DSL
// errors). The leading TrimSpace guard in dslSwatches makes this reachable;
// the old code's len==0 check was unreachable because strings.Split("",",")
// returns [""].
func TestDslSwatches_EmptyArg(t *testing.T) {
	out := dslSwatches("")
	assertSVG(t, "swatches empty", out)
	if !strings.Contains(out, "swatches: need") {
		t.Errorf("expected errArt placeholder mentioning 'swatches: need', got:\n%s", out)
	}
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

// Empty input must emit the errArt placeholder. Same fix-pattern as
// dslSwatches: a leading TrimSpace guard catches what the len==0 check
// couldn't.
func TestDslPills_EmptyArg(t *testing.T) {
	out := dslPills("")
	assertSVG(t, "pills empty", out)
	if !strings.Contains(out, "pills: need") {
		t.Errorf("expected errArt placeholder mentioning 'pills: need', got:\n%s", out)
	}
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

// TestResolveDSL_NewPrefixes confirms text/flow/box/arrow are routed
// through resolveDSL rather than falling into the "unknown" placeholder.
func TestResolveDSL_NewPrefixes(t *testing.T) {
	cases := []struct {
		in       string
		wantSub  string
		notUnk   bool
	}{
		{`text:"hi"`, "hi", true},
		{`flow:[a,b]`, "<rect", true},
		{`box:hello`, "hello", true},
		{`arrow:"go"`, "go", true},
	}
	for _, tc := range cases {
		svg, hidden := resolveDSL(tc.in)
		if hidden {
			t.Errorf("resolveDSL(%q) hidden = true, want false", tc.in)
		}
		assertSVG(t, "resolveDSL "+tc.in, svg, tc.wantSub)
		if tc.notUnk && strings.Contains(svg, "unknown dsl") {
			t.Errorf("resolveDSL(%q) fell through to unknown branch: %s", tc.in, svg)
		}
	}
}

// ──────────────────────────────────────────────────────────────────────
// text
// ──────────────────────────────────────────────────────────────────────

func TestDslText_Valid(t *testing.T) {
	// Default font/attrs.
	out := dslText(`"hello world"`)
	assertSVG(t, "text default", out, "hello world", "IBM Plex Sans")

	// @mono with multi-line, attrs.
	out2 := dslText(`"line one\nline two"@mono?size=12&align=center&weight=600&color=accent`)
	assertSVG(t, "text mono multi", out2, "line one", "line two", "IBM Plex Mono",
		"text-anchor='middle'", "font-size='12'", "font-weight='600'", "var(--accent)")

	// @serif.
	out3 := dslText(`"editorial"@serif`)
	assertSVG(t, "text serif", out3, "editorial", "IBM Plex Serif")

	// align=right.
	out4 := dslText(`"end"?align=right`)
	assertSVG(t, "text right", out4, "text-anchor='end'")

	// color=rule-2 (verifies hyphenated color names parse).
	out5 := dslText(`"x"?color=rule-2`)
	assertSVG(t, "text color rule-2", out5, "var(--rule-2)")
}

func TestDslText_Invalid(t *testing.T) {
	cases := map[string]string{
		"no quotes":           `bare-text`,
		"unterminated quote":  `"oops`,
		"empty body":          `""`,
		"unknown font":        `"x"@cursive`,
		"size below range":    `"x"?size=2`,
		"size above range":    `"x"?size=99`,
		"weight out of range": `"x"?weight=900`,
		"bad align":           `"x"?align=middle`,
		"unknown color":       `"x"?color=tomato`,
		"unknown attr":        `"x"?wobble=12`,
		"bad attr no eq":      `"x"?broken`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			out := dslText(in)
			assertSVG(t, "text "+name, out, "text:")
		})
	}
}

// ──────────────────────────────────────────────────────────────────────
// flow
// ──────────────────────────────────────────────────────────────────────

func TestDslFlow_Valid(t *testing.T) {
	out := dslFlow(`[a,b,c]`)
	assertSVG(t, "flow simple", out, "<rect", "a", "b", "c")

	out2 := dslFlow(`[client?icon=user,api?icon=api,db?icon=database]`)
	assertSVG(t, "flow icons", out2, "client", "api", "db")

	// One element is fine (no arrow needed).
	out3 := dslFlow(`[only]`)
	assertSVG(t, "flow single", out3, "only")
}

func TestDslFlow_Invalid(t *testing.T) {
	cases := map[string]string{
		"empty list":      `[]`,
		"no brackets":     `a,b,c`,
		"empty token":     `[a,,b]`,
		"unknown icon":    `[x?icon=nope]`,
		"bad attr":        `[x?color=red]`,
		"icon empty":      `[x?icon=]`,
		"only brackets":   `[`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			out := dslFlow(in)
			assertSVG(t, "flow "+name, out, "flow:")
		})
	}
}

// ──────────────────────────────────────────────────────────────────────
// box
// ──────────────────────────────────────────────────────────────────────

func TestDslBox_Valid(t *testing.T) {
	out := dslBox(`hello`)
	assertSVG(t, "box plain", out, "<rect", "hello")

	out2 := dslBox(`web?icon=browser`)
	assertSVG(t, "box icon", out2, "web", "<g")

	out3 := dslBox(`db?icon=database`)
	assertSVG(t, "box db icon", out3, "db")
}

func TestDslBox_Invalid(t *testing.T) {
	cases := map[string]string{
		"empty":          ``,
		"unknown icon":   `x?icon=nope`,
		"unknown attr":   `x?color=red`,
		"bad attr no eq": `x?broken`,
		"empty label":    `?icon=server`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			out := dslBox(in)
			assertSVG(t, "box "+name, out, "box:")
		})
	}
}

// ──────────────────────────────────────────────────────────────────────
// arrow
// ──────────────────────────────────────────────────────────────────────

func TestDslArrow_Valid(t *testing.T) {
	out := dslArrow(`"go"`)
	assertSVG(t, "arrow default", out, "go", "<polygon", "<line")

	out2 := dslArrow(`"down"?dir=down`)
	assertSVG(t, "arrow down", out2, "down")

	out3 := dslArrow(`"up"?dir=up`)
	assertSVG(t, "arrow up", out3, "up")

	out4 := dslArrow(`"left"?dir=left`)
	assertSVG(t, "arrow left", out4, "left")

	out5 := dslArrow(`"right"?dir=right`)
	assertSVG(t, "arrow right", out5, "right")
}

func TestDslArrow_Invalid(t *testing.T) {
	cases := map[string]string{
		"no quotes":          `bare`,
		"unterminated quote": `"oops`,
		"empty label":        `""`,
		"unknown dir":        `"x"?dir=sideways`,
		"unknown attr":       `"x"?color=red`,
		"bad attr no eq":     `"x"?broken`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			out := dslArrow(in)
			assertSVG(t, "arrow "+name, out, "arrow:")
		})
	}
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
