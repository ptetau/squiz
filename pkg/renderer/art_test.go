package renderer

import (
	"strings"
	"testing"
)

// TestRenderArt_Dispatch covers each branch of the RenderArt switch.
func TestRenderArt_Dispatch(t *testing.T) {
	cases := []struct {
		name      string
		art       string
		letterIdx int
		wantSVG   bool   // expect non-empty SVG output
		wantSub   string // optional substring the output must contain
		hidden    bool
	}{
		{
			name: "empty falls back to autoArt",
			art:  "", letterIdx: 0,
			wantSVG: true, wantSub: "<svg", hidden: false,
		},
		{
			name: "none returns hidden empty",
			art:  "none", letterIdx: 0,
			wantSVG: false, hidden: true,
		},
		{
			name: "raw svg passes through",
			art:  `<svg viewBox='0 0 10 10'><rect width='10' height='10'/></svg>`, letterIdx: 0,
			wantSVG: true, wantSub: "<rect", hidden: false,
		},
		{
			name: "wf prefix resolves named entry",
			art:  "wf:calendar-grid", letterIdx: 0,
			wantSVG: true, wantSub: "<svg", hidden: false,
		},
		{
			name: "dsl grid dispatches",
			art:  "grid:3x3", letterIdx: 0,
			wantSVG: true, wantSub: "<svg", hidden: false,
		},
		{
			name: "garbage produces unknown-art placeholder",
			art:  "foo", letterIdx: 0,
			wantSVG: true, wantSub: "unknown art", hidden: false,
		},
		{
			name: "whitespace around input is trimmed",
			art:  "  none  ", letterIdx: 0,
			wantSVG: false, hidden: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svg, hidden := RenderArt(tc.art, tc.letterIdx)
			if hidden != tc.hidden {
				t.Errorf("hidden = %v, want %v", hidden, tc.hidden)
			}
			if tc.wantSVG && svg == "" {
				t.Errorf("expected non-empty SVG, got empty")
			}
			if !tc.wantSVG && svg != "" {
				t.Errorf("expected empty SVG, got %q", svg)
			}
			if tc.wantSub != "" && !strings.Contains(svg, tc.wantSub) {
				t.Errorf("SVG missing substring %q\ngot: %s", tc.wantSub, svg)
			}
		})
	}
}

// TestRenderArt_UnknownEscapesAngleBrackets makes sure the unknown-art
// fallback HTML-escapes user input so a typo can't inject markup.
func TestRenderArt_UnknownEscapesAngleBrackets(t *testing.T) {
	svg, hidden := RenderArt("<script>", 0)
	if hidden {
		t.Fatalf("hidden = true, want false for unknown form")
	}
	if strings.Contains(svg, "<script>") {
		t.Errorf("unescaped <script> leaked into output:\n%s", svg)
	}
	if !strings.Contains(svg, "&lt;script&gt;") {
		t.Errorf("expected escaped &lt;script&gt; in output, got:\n%s", svg)
	}
}

// TestLetterFor locks in the current letter-defaulting contract.
func TestLetterFor(t *testing.T) {
	cases := []struct {
		idx  int
		want string
	}{
		{0, "A"},
		{1, "B"},
		{25, "Z"},
		// Beyond Z the implementation rolls into a 2-letter code:
		// idx=26 → A + (26/26)-1 = 'A', 26%26 = 'A' → "AA"
		{26, "AA"},
		{27, "AB"},
		{51, "AZ"},
		{52, "BA"},
	}
	for _, tc := range cases {
		got := LetterFor(tc.idx)
		if got != tc.want {
			t.Errorf("LetterFor(%d) = %q, want %q", tc.idx, got, tc.want)
		}
	}
}

// TestLetterFor_NegativeIdx confirms negative indices return "" instead
// of panicking.
func TestLetterFor_NegativeIdx(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("LetterFor(-1) panicked: %v", r)
		}
	}()
	if got := LetterFor(-1); got != "" {
		t.Errorf("LetterFor(-1) = %q, want \"\"", got)
	}
}

// TestLetterFor_HighIdxDoesNotPanic locks in current behaviour for the
// 2-letter range; we only care that it doesn't crash.
func TestLetterFor_HighIdxDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("LetterFor(100) panicked: %v", r)
		}
	}()
	_ = LetterFor(100)
	_ = LetterFor(675) // (26*26)-1, last valid 2-letter idx
}

// TestAutoArt_DeterministicPerLetter — same idx → same SVG; different
// idx → different SVG (within the auto-pattern cycle).
func TestAutoArt_DeterministicPerLetter(t *testing.T) {
	a1, _ := RenderArt("", 0)
	a2, _ := RenderArt("", 0)
	if a1 != a2 {
		t.Errorf("autoArt for idx 0 not deterministic:\n%s\n---\n%s", a1, a2)
	}

	b, _ := RenderArt("", 1)
	if a1 == b {
		t.Errorf("autoArt for idx 0 and 1 should differ, both =\n%s", a1)
	}

	// Cycles modulo len(autoPatterns); confirm the wrap-around matches.
	wrap, _ := RenderArt("", len(autoPatterns))
	if wrap != a1 {
		t.Errorf("autoArt(len) should equal autoArt(0), got different SVGs")
	}
}

// TestAutoArt_NegativeIdx — negative indices clamp to 0 instead of
// panicking on a negative modulus.
func TestAutoArt_NegativeIdx(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("autoArt(-1) panicked: %v", r)
		}
	}()
	zero, _ := RenderArt("", 0)
	neg, _ := RenderArt("", -1)
	if neg != zero {
		t.Errorf("autoArt(-1) should clamp to autoArt(0)")
	}
}
