package renderer

import (
	"strings"
	"testing"
)

// ──────────────────────────────────────────────────────────────────────
// callout
// ──────────────────────────────────────────────────────────────────────

func TestDslCallout_Valid(t *testing.T) {
	out := dslCallout(`"v2"@70,12->60,30`)
	assertSVG(t, "callout valid", out, "v2", "<line", "<polygon")

	// With color override.
	out2 := dslCallout(`"hot"@10,10->50,40?color=ink`)
	assertSVG(t, "callout color", out2, "hot", "var(--ink)")

	// Same-point edge case (length=0) shouldn't NaN out the arrow head.
	out3 := dslCallout(`"x"@50,30->50,30`)
	assertSVG(t, "callout zero len", out3, "x")
	if strings.Contains(out3, "NaN") {
		t.Errorf("callout zero-len produced NaN: %s", out3)
	}
}

func TestDslCallout_Invalid(t *testing.T) {
	cases := map[string]string{
		"no quotes":           `bare@1,2->3,4`,
		"empty label":         `""@1,2->3,4`,
		"label too long":      `"this-is-way-too-long"@1,2->3,4`,
		"missing @":           `"x"1,2->3,4`,
		"missing arrow":       `"x"@1,2,3,4`,
		"bad src coords":      `"x"@1->3,4`,
		"bad dst coords":      `"x"@1,2->three,4`,
		"src x out of range":  `"x"@200,2->3,4`,
		"dst y out of range":  `"x"@1,2->3,99`,
		"bad attr":            `"x"@1,2->3,4?color=tomato`,
		"unknown attr":        `"x"@1,2->3,4?weight=400`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			out := dslCallout(in)
			assertSVG(t, "callout "+name, out, "callout")
		})
	}
}

// ──────────────────────────────────────────────────────────────────────
// brace
// ──────────────────────────────────────────────────────────────────────

func TestDslBrace_Valid(t *testing.T) {
	out := dslBrace(`"steady"@20,40,60`)
	assertSVG(t, "brace up", out, "steady", "<path")

	out2 := dslBrace(`"opts"@10,20,80?dir=down`)
	assertSVG(t, "brace down", out2, "opts", "<path")
}

func TestDslBrace_Invalid(t *testing.T) {
	cases := map[string]string{
		"no quotes":     `steady@20,40,60`,
		"empty label":   `""@20,40,60`,
		"missing @":     `"x"20,40,60`,
		"missing w":     `"x"@20,40`,
		"w too small":   `"x"@10,20,2`,
		"w too big":     `"x"@1,20,99`,
		"bad dir":       `"x"@10,20,60?dir=left`,
		"unknown attr":  `"x"@10,20,60?color=accent`,
		"bad attr noeq": `"x"@10,20,60?broken`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			out := dslBrace(in)
			assertSVG(t, "brace "+name, out, "brace")
		})
	}
}

// ──────────────────────────────────────────────────────────────────────
// divider
// ──────────────────────────────────────────────────────────────────────

func TestDslDivider_Valid(t *testing.T) {
	out := dslDivider(`vs`)
	assertSVG(t, "divider default", out, "vs", "<line")

	out2 := dslDivider(`vs@40`)
	assertSVG(t, "divider x=40", out2, "vs")

	out3 := dslDivider(`vs@60?color=ink`)
	assertSVG(t, "divider color", out3, "vs", "var(--ink)")
}

func TestDslDivider_Invalid(t *testing.T) {
	cases := map[string]string{
		"missing vs":     `foo`,
		"x out of range": `vs@5`,
		"bad x":          `vs@foo`,
		"unknown color":  `vs?color=tomato`,
		"unknown attr":   `vs?dir=up`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			out := dslDivider(in)
			assertSVG(t, "divider "+name, out, "divider")
		})
	}
}

// ──────────────────────────────────────────────────────────────────────
// badge
// ──────────────────────────────────────────────────────────────────────

func TestDslBadge_Valid(t *testing.T) {
	out := dslBadge(`tick@88,8`)
	assertSVG(t, "badge tick", out, "<circle", "<path")

	out2 := dslBadge(`cross@10,10`)
	assertSVG(t, "badge cross", out2, "<circle")

	out3 := dslBadge(`warn@50,30`)
	assertSVG(t, "badge warn", out3, "<polygon", "!")

	out4 := dslBadge(`star@50,30`)
	assertSVG(t, "badge star", out4, "<polygon")

	out5 := dslBadge(`dot@50,30?color=ink`)
	assertSVG(t, "badge dot color", out5, "<circle", "var(--ink)")
}

func TestDslBadge_Invalid(t *testing.T) {
	cases := map[string]string{
		"unknown kind":   `wobble@10,10`,
		"missing @":      `tick10,10`,
		"bad coords":     `tick@one,two`,
		"x out of range": `tick@200,10`,
		"y out of range": `tick@10,99`,
		"unknown color":  `tick@10,10?color=tomato`,
		"unknown attr":   `tick@10,10?size=20`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			out := dslBadge(in)
			assertSVG(t, "badge "+name, out, "badge")
		})
	}
}

// ──────────────────────────────────────────────────────────────────────
// range
// ──────────────────────────────────────────────────────────────────────

func TestDslRange_Valid(t *testing.T) {
	out := dslRange(`3-50@10,30,80`)
	assertSVG(t, "range valid", out, "<line", "3", "50")

	out2 := dslRange(`10-100@10,30,80?label=ms`)
	assertSVG(t, "range with label", out2, "10", "100", "ms")
}

func TestDslRange_Invalid(t *testing.T) {
	cases := map[string]string{
		"missing @":      `3-50`,
		"missing range":  `@10,30,80`,
		"bad LO":         `foo-50@10,30,80`,
		"bad HI":         `3-bar@10,30,80`,
		"missing coords": `3-50@10,30`,
		"w too small":    `3-50@10,30,5`,
		"w too big":      `3-50@1,30,99`,
		"unknown attr":   `3-50@10,30,80?color=accent`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			out := dslRange(in)
			assertSVG(t, "range "+name, out, "range")
		})
	}
}

// ──────────────────────────────────────────────────────────────────────
// baseline
// ──────────────────────────────────────────────────────────────────────

func TestDslBaseline_Valid(t *testing.T) {
	out := dslBaseline(`50@10,20,80?label=p99`)
	assertSVG(t, "baseline valid", out, "<line", "50", "p99")

	// No label.
	out2 := dslBaseline(`42@10,30,60`)
	assertSVG(t, "baseline no label", out2, "<line", "42")
}

func TestDslBaseline_Invalid(t *testing.T) {
	cases := map[string]string{
		"missing @":      `50`,
		"empty value":    `@10,20,80`,
		"bad value":      `foo@10,20,80`,
		"missing coords": `50@10,20`,
		"w too small":    `50@10,20,5`,
		"w too big":      `50@10,20,99`,
		"unknown attr":   `50@10,20,80?color=accent`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			out := dslBaseline(in)
			assertSVG(t, "baseline "+name, out, "baseline")
		})
	}
}

// ──────────────────────────────────────────────────────────────────────
// times
// ──────────────────────────────────────────────────────────────────────

func TestDslTimes_Valid(t *testing.T) {
	out := dslTimes(`3`)
	assertSVG(t, "times valid", out, "×3")

	out2 := dslTimes(`4?label=sensors`)
	assertSVG(t, "times with label", out2, "×4", "sensors")
}

func TestDslTimes_Invalid(t *testing.T) {
	cases := map[string]string{
		"empty":         ``,
		"not a number":  `foo`,
		"zero":          `0`,
		"negative":      `-1`,
		"too big":       `1000`,
		"float":         `3.5`,
		"unknown attr":  `3?color=accent`,
		"bad attr noeq": `3?broken`,
	}
	for name, in := range cases {
		t.Run(name, func(t *testing.T) {
			out := dslTimes(in)
			assertSVG(t, "times "+name, out, "times")
		})
	}
}

// ──────────────────────────────────────────────────────────────────────
// dispatcher routing — confirm all 7 new prefixes route through resolveDSL
// ──────────────────────────────────────────────────────────────────────

func TestResolveDSL_NewAnnotationPrefixes(t *testing.T) {
	cases := []struct {
		in      string
		wantSub string
	}{
		{`callout:"v2"@70,12->60,30`, "v2"},
		{`brace:"steady"@20,40,60`, "steady"},
		{`divider:vs`, "vs"},
		{`badge:tick@88,8`, "<svg"},
		{`range:3-50@10,30,80`, "3"},
		{`baseline:50@10,20,80?label=p99`, "p99"},
		{`times:3`, "×3"},
	}
	for _, tc := range cases {
		svg, hidden := resolveDSL(tc.in)
		if hidden {
			t.Errorf("resolveDSL(%q) hidden = true, want false", tc.in)
		}
		assertSVG(t, "resolveDSL "+tc.in, svg, tc.wantSub)
		if strings.Contains(svg, "unknown dsl") {
			t.Errorf("resolveDSL(%q) fell through to unknown branch: %s", tc.in, svg)
		}
	}
}
