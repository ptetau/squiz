package renderer

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Annotation primitives (v0.8.0 Layer 1): GLUE marks that overlay or sit
// alongside nouns (charts, icons, screens) to turn them into statements.
// Each primitive emits a complete viewBox 0 0 100 60 SVG so it's usable
// standalone AND composable via <use href="…"/>.
//
// All primitives use theme CSS vars so they inherit the active theme.
// Inner content is centered in the 100×60 canvas — preserveAspectRatio
// on <use> handles rescaling.

// ──────────────────────────────────────────────────────────────────────
// helpers
// ──────────────────────────────────────────────────────────────────────

// parseQuotedHead pulls a leading "…" quoted string off `s` and returns
// (body, rest, ok). `rest` is whatever comes after the closing quote
// (no leading whitespace trim — callers may want the raw byte).
func parseQuotedHead(s string) (body, rest string, ok bool) {
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, `"`) {
		return "", "", false
	}
	closeIdx := strings.Index(s[1:], `"`)
	if closeIdx < 0 {
		return "", "", false
	}
	return s[1 : 1+closeIdx], s[2+closeIdx:], true
}

// colorVar maps an `accent | ink | ink-2` short name to its CSS var.
// Returns ("", false) on unknown name.
func colorVar(name string) (string, bool) {
	switch name {
	case "accent":
		return "var(--accent)", true
	case "ink":
		return "var(--ink)", true
	case "ink-2":
		return "var(--ink-2)", true
	case "ink-3":
		return "var(--ink-3)", true
	}
	return "", false
}

// parseColorAttr scans `?color=NAME` (or no suffix) and returns the
// CSS var to use plus an error placeholder string (empty on success).
// `defaultColor` is the CSS var to use when no attr is present.
func parseColorAttr(suffix, defaultColor, kindLabel string) (cssVar, errStr string) {
	suffix = strings.TrimSpace(suffix)
	if suffix == "" {
		return defaultColor, ""
	}
	if !strings.HasPrefix(suffix, "?") {
		return "", kindLabel + ": bad attrs"
	}
	for _, kv := range strings.Split(suffix[1:], "&") {
		if kv == "" {
			continue
		}
		eq := strings.Index(kv, "=")
		if eq < 0 {
			return "", kindLabel + ": bad attr " + kv
		}
		k, v := kv[:eq], kv[eq+1:]
		switch k {
		case "color":
			cv, ok := colorVar(v)
			if !ok {
				return "", kindLabel + ": unknown color " + v
			}
			cssVar = cv
		default:
			return "", kindLabel + ": unknown attr " + k
		}
	}
	if cssVar == "" {
		cssVar = defaultColor
	}
	return cssVar, ""
}

// parseFloatList splits "x,y" / "x,y,w" into floats. Returns (vals, ok).
func parseFloatList(s string) ([]float64, bool) {
	parts := strings.Split(s, ",")
	out := make([]float64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, false
		}
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return nil, false
		}
		out = append(out, v)
	}
	return out, len(out) > 0
}

func inRange(v, lo, hi float64) bool { return v >= lo && v <= hi }

// ──────────────────────────────────────────────────────────────────────
// 1. callout:"label"@x1,y1->x2,y2[?color=C]
// ──────────────────────────────────────────────────────────────────────

func dslCallout(args string) string {
	label, rest, ok := parseQuotedHead(args)
	if !ok {
		return errArt(`callout: need "label"`)
	}
	if label == "" {
		return errArt(`callout: empty label`)
	}
	if len(label) > 12 {
		return errArt(`callout: label too long`)
	}
	if !strings.HasPrefix(rest, "@") {
		return errArt(`callout: need @x1,y1->x2,y2`)
	}
	rest = rest[1:]
	// Split off optional ?attrs first.
	var attrs string
	if q := strings.Index(rest, "?"); q >= 0 {
		attrs = rest[q:]
		rest = rest[:q]
	}
	// Split source -> dest on "->".
	arrowIdx := strings.Index(rest, "->")
	if arrowIdx < 0 {
		return errArt(`callout: need x1,y1->x2,y2`)
	}
	srcStr := rest[:arrowIdx]
	dstStr := rest[arrowIdx+2:]
	src, ok1 := parseFloatList(srcStr)
	dst, ok2 := parseFloatList(dstStr)
	if !ok1 || !ok2 || len(src) != 2 || len(dst) != 2 {
		return errArt(`callout: bad coords`)
	}
	x1, y1, x2, y2 := src[0], src[1], dst[0], dst[1]
	if !inRange(x1, 0, 100) || !inRange(x2, 0, 100) || !inRange(y1, 0, 60) || !inRange(y2, 0, 60) {
		return errArt(`callout: coords out of range`)
	}
	color, errStr := parseColorAttr(attrs, "var(--accent)", "callout")
	if errStr != "" {
		return errArt(errStr)
	}

	// Arrow head at (x2,y2). Build a small triangle pointing along the
	// (src→dst) direction. Use a normalized perpendicular for the base.
	dx, dy := x2-x1, y2-y1
	length := dx*dx + dy*dy
	var hx1, hy1, hx2, hy2 float64
	const headLen = 3.0
	if length == 0 {
		hx1, hy1, hx2, hy2 = x2-headLen, y2-headLen/2, x2-headLen, y2+headLen/2
	} else {
		// Normalize.
		invLen := 1.0 / math.Sqrt(length)
		ux, uy := dx*invLen, dy*invLen
		// Perpendicular (rotated 90°).
		px, py := -uy, ux
		bx, by := x2-ux*headLen, y2-uy*headLen
		hx1, hy1 = bx+px*headLen*0.5, by+py*headLen*0.5
		hx2, hy2 = bx-px*headLen*0.5, by-py*headLen*0.5
	}

	return fmt.Sprintf(
		`<svg viewBox='0 0 100 60' style='width:80%%;height:auto'><line x1='%.2f' y1='%.2f' x2='%.2f' y2='%.2f' stroke='%s' stroke-width='1'/><polygon points='%.2f,%.2f %.2f,%.2f %.2f,%.2f' fill='%s'/><text x='%.2f' y='%.2f' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' font-weight='600' fill='%s'>%s</text></svg>`,
		x1, y1, x2, y2, color,
		x2, y2, hx1, hy1, hx2, hy2, color,
		x1, y1-2, color, escapeXML(label),
	)
}


// ──────────────────────────────────────────────────────────────────────
// 2. brace:"{label}"@x,y,w[?dir=down]
// ──────────────────────────────────────────────────────────────────────

func dslBrace(args string) string {
	label, rest, ok := parseQuotedHead(args)
	if !ok {
		return errArt(`brace: need "label"`)
	}
	if label == "" {
		return errArt(`brace: empty label`)
	}
	if !strings.HasPrefix(rest, "@") {
		return errArt(`brace: need @x,y,w`)
	}
	rest = rest[1:]
	var attrs string
	if q := strings.Index(rest, "?"); q >= 0 {
		attrs = rest[q:]
		rest = rest[:q]
	}
	vals, ok := parseFloatList(rest)
	if !ok || len(vals) != 3 {
		return errArt(`brace: need x,y,w`)
	}
	x, y, w := vals[0], vals[1], vals[2]
	if !inRange(w, 4, 96) {
		return errArt(`brace: w out of range`)
	}

	// Default dir=up: label above brace; brace opens downward (toward the
	// thing it's grouping). dir=down: label below brace; brace opens
	// upward.
	dir := "up"
	if strings.HasPrefix(attrs, "?") {
		for _, kv := range strings.Split(attrs[1:], "&") {
			if kv == "" {
				continue
			}
			eq := strings.Index(kv, "=")
			if eq < 0 {
				return errArt(`brace: bad attr ` + kv)
			}
			k, v := kv[:eq], kv[eq+1:]
			switch k {
			case "dir":
				if v != "up" && v != "down" {
					return errArt(`brace: bad dir ` + v)
				}
				dir = v
			default:
				return errArt(`brace: unknown attr ` + k)
			}
		}
	}

	// Brace geometry: corner curves at each end of a horizontal rule, with a
	// center "notch" curve pointing toward the grouped content. The whole
	// brace lives in the band [y .. y±2*armH] depending on direction.
	const armH = 3.0
	const stroke = `stroke='var(--accent)' stroke-width='1' fill='none' stroke-linejoin='round' stroke-linecap='round'`
	midX := x + w/2
	// `sgn` flips the y-direction of the brace body for dir=down.
	sgn := 1.0
	if dir == "down" {
		sgn = -1.0
	}
	edgeY := y               // the label-side edge
	baseY := y + sgn*armH    // the horizontal rule the arms reach
	notchY := y + sgn*2*armH // the tip of the center notch
	// Path: M start → Q corner into the rule → L to just before notch →
	// Q notch-down → Q notch-back → L to far end → Q corner out.
	// 7 verbs that take coords: M(2) Q(4) L(2) Q(4) Q(4) L(2) Q(4) = 22 floats.
	// Args below must total 22.
	path := fmt.Sprintf(
		"M %.2f %.2f Q %.2f %.2f %.2f %.2f L %.2f %.2f Q %.2f %.2f %.2f %.2f Q %.2f %.2f %.2f %.2f L %.2f %.2f Q %.2f %.2f %.2f %.2f",
		// M start: top/bottom-left tip of the arm
		x, edgeY,
		// Q corner: control at (x,baseY), end at (x+armH, baseY)
		x, baseY, x+armH, baseY,
		// L to the notch's left foot
		midX-armH, baseY,
		// Q notch left half: control at (midX, baseY), end at (midX, notchY)
		midX, baseY, midX, notchY,
		// Q notch right half: control at (midX, baseY), end at (midX+armH, baseY)
		midX, baseY, midX+armH, baseY,
		// L to the far corner approach
		x+w-armH, baseY,
		// Q corner out: control at (x+w, baseY), end at (x+w, edgeY)
		x+w, baseY, x+w, edgeY,
	)
	var labelY string
	if dir == "up" {
		labelY = fmt.Sprintf("%.2f", y-1)
	} else {
		labelY = fmt.Sprintf("%.2f", y+armH*2+6)
	}

	return fmt.Sprintf(
		`<svg viewBox='0 0 100 60' style='width:80%%;height:auto'><path d='%s' %s/><text x='%.2f' y='%s' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' font-weight='600' fill='var(--ink)'>%s</text></svg>`,
		path, stroke, midX, labelY, escapeXML(label),
	)
}

// ──────────────────────────────────────────────────────────────────────
// 3. divider:vs[@x][?color=C]
// ──────────────────────────────────────────────────────────────────────

func dslDivider(args string) string {
	args = strings.TrimSpace(args)
	// Must start with the literal "vs" glyph.
	if !strings.HasPrefix(args, "vs") {
		return errArt(`divider: need vs`)
	}
	rest := args[2:]
	x := 50.0
	if strings.HasPrefix(rest, "@") {
		rest = rest[1:]
		var coord string
		if q := strings.Index(rest, "?"); q >= 0 {
			coord = rest[:q]
			rest = rest[q:]
		} else {
			coord = rest
			rest = ""
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(coord), 64)
		if err != nil {
			return errArt(`divider: bad x`)
		}
		x = v
	}
	if !inRange(x, 10, 90) {
		return errArt(`divider: x out of range`)
	}
	color, errStr := parseColorAttr(rest, "var(--accent)", "divider")
	if errStr != "" {
		return errArt(errStr)
	}

	// Dashed rule with a "vs" glyph centered vertically, on a small surface
	// chip so it reads against busy content.
	const top, bot = 6.0, 54.0
	const chipR = 5.0
	return fmt.Sprintf(
		`<svg viewBox='0 0 100 60' style='width:80%%;height:auto'><line x1='%.2f' y1='%.2f' x2='%.2f' y2='%.2f' stroke='var(--rule-2)' stroke-width='0.8' stroke-dasharray='3 2'/><circle cx='%.2f' cy='30' r='%.2f' fill='var(--bg)' stroke='%s' stroke-width='0.8'/><text x='%.2f' y='32' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' font-weight='700' fill='%s'>vs</text></svg>`,
		x, top, x, bot,
		x, chipR, color,
		x, color,
	)
}

// ──────────────────────────────────────────────────────────────────────
// 4. badge:KIND@x,y[?color=C]
// ──────────────────────────────────────────────────────────────────────

func dslBadge(args string) string {
	// Split KIND@x,y[?color=C]
	at := strings.Index(args, "@")
	if at < 0 {
		return errArt(`badge: need KIND@x,y`)
	}
	kind := strings.TrimSpace(args[:at])
	rest := args[at+1:]
	switch kind {
	case "tick", "cross", "warn", "star", "dot", "check":
		// "check" is an accepted synonym for "tick" per the spec text.
	default:
		return errArt(`badge: unknown kind ` + kind)
	}
	if kind == "check" {
		kind = "tick"
	}
	var attrs string
	if q := strings.Index(rest, "?"); q >= 0 {
		attrs = rest[q:]
		rest = rest[:q]
	}
	vals, ok := parseFloatList(rest)
	if !ok || len(vals) != 2 {
		return errArt(`badge: bad coords`)
	}
	x, y := vals[0], vals[1]
	if !inRange(x, 0, 100) || !inRange(y, 0, 60) {
		return errArt(`badge: coords out of range`)
	}

	// Default color depends on kind (tick=accent, cross=ink, warn=accent,
	// star=accent, dot=accent) but the user override always wins.
	defaultColor := "var(--accent)"
	if kind == "cross" {
		defaultColor = "var(--ink)"
	}
	color, errStr := parseColorAttr(attrs, defaultColor, "badge")
	if errStr != "" {
		return errArt(errStr)
	}

	// All glyphs sit in a ~12-unit footprint centered on (x,y).
	const r = 6.0
	var glyph string
	switch kind {
	case "tick":
		// Filled circle background, white check stroke on top.
		glyph = fmt.Sprintf(
			`<circle cx='%.2f' cy='%.2f' r='%.2f' fill='%s'/><path d='M %.2f %.2f L %.2f %.2f L %.2f %.2f' fill='none' stroke='var(--bg)' stroke-width='1.4' stroke-linecap='round' stroke-linejoin='round'/>`,
			x, y, r, color,
			x-3, y, x-1, y+2.5, x+3, y-2.5,
		)
	case "cross":
		glyph = fmt.Sprintf(
			`<circle cx='%.2f' cy='%.2f' r='%.2f' fill='%s'/><path d='M %.2f %.2f L %.2f %.2f M %.2f %.2f L %.2f %.2f' stroke='var(--bg)' stroke-width='1.4' stroke-linecap='round'/>`,
			x, y, r, color,
			x-2.5, y-2.5, x+2.5, y+2.5,
			x-2.5, y+2.5, x+2.5, y-2.5,
		)
	case "warn":
		// Triangle outline with "!" inside.
		glyph = fmt.Sprintf(
			`<polygon points='%.2f,%.2f %.2f,%.2f %.2f,%.2f' fill='%s' stroke='%s' stroke-width='0.6' stroke-linejoin='round'/><text x='%.2f' y='%.2f' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' font-weight='700' fill='var(--bg)'>!</text>`,
			x, y-r, x-r, y+r-0.5, x+r, y+r-0.5, color, color,
			x, y+2.5,
		)
	case "star":
		// 5-point star: build the 10 points around (x,y).
		glyph = starPath(x, y, r, color)
	case "dot":
		glyph = fmt.Sprintf(`<circle cx='%.2f' cy='%.2f' r='%.2f' fill='%s'/>`, x, y, r-1.5, color)
	}

	return fmt.Sprintf(
		`<svg viewBox='0 0 100 60' style='width:80%%;height:auto'>%s</svg>`,
		glyph,
	)
}

// starPath emits a 5-point star polygon centered on (cx,cy) with outer
// radius r. Uses pre-computed unit coords so we don't pull math/trig here.
func starPath(cx, cy, r float64, color string) string {
	// 10 points alternating outer (r) / inner (r*0.4) starting at the top.
	// Angles in degrees: 270, 306, 342, 18, 54, 90, 126, 162, 198, 234.
	// Pre-computed (cos,sin) pairs.
	pts := [][2]float64{
		{0, -1},                         // 270 outer
		{0.293892626, -0.404508497},     // 306 inner
		{0.951056516, -0.309016994},     // 342 outer
		{0.475528258, 0.154508497},      // 18  inner
		{0.587785252, 0.809016994},      // 54  outer
		{0, 0.5},                        // 90  inner
		{-0.587785252, 0.809016994},     // 126 outer
		{-0.475528258, 0.154508497},     // 162 inner
		{-0.951056516, -0.309016994},    // 198 outer
		{-0.293892626, -0.404508497},    // 234 inner
	}
	var b strings.Builder
	b.WriteString(`<polygon points='`)
	for i, p := range pts {
		scale := r
		// Inner points (odd indices) shrink to 0.4r.
		if i%2 == 1 {
			scale = r * 0.4
		}
		if i > 0 {
			b.WriteString(" ")
		}
		// Re-scale: the unit pts above mix r and 0.5 indiscriminately —
		// just use them as direction (we already encoded inner shrink
		// implicitly via 0.5 for the y=90 case, but to be safe override).
		dx := p[0]
		dy := p[1]
		// Normalize ensures a clean star regardless of the precomputed mix.
		n := dx*dx + dy*dy
		invN := 1.0 / math.Sqrt(n)
		dx *= invN
		dy *= invN
		fmt.Fprintf(&b, "%.2f,%.2f", cx+dx*scale, cy+dy*scale)
	}
	fmt.Fprintf(&b, `' fill='%s'/>`, color)
	return b.String()
}

// ──────────────────────────────────────────────────────────────────────
// 5. range:LO-HI@x,y,w[?label=TEXT]
// ──────────────────────────────────────────────────────────────────────

func dslRange(args string) string {
	// Split LO-HI@x,y,w[?label=…]
	at := strings.Index(args, "@")
	if at < 0 {
		return errArt(`range: need LO-HI@x,y,w`)
	}
	rangeStr := args[:at]
	rest := args[at+1:]
	// Split LO-HI on the LAST '-' so negatives aren't supported (per spec
	// examples it's always non-negative). Use a simple split on '-'.
	dash := strings.Index(rangeStr, "-")
	if dash < 0 {
		return errArt(`range: need LO-HI`)
	}
	loStr := strings.TrimSpace(rangeStr[:dash])
	hiStr := strings.TrimSpace(rangeStr[dash+1:])
	if loStr == "" || hiStr == "" {
		return errArt(`range: need LO-HI`)
	}
	if _, err := strconv.ParseFloat(loStr, 64); err != nil {
		return errArt(`range: LO not a number`)
	}
	if _, err := strconv.ParseFloat(hiStr, 64); err != nil {
		return errArt(`range: HI not a number`)
	}
	var attrs string
	if q := strings.Index(rest, "?"); q >= 0 {
		attrs = rest[q:]
		rest = rest[:q]
	}
	vals, ok := parseFloatList(rest)
	if !ok || len(vals) != 3 {
		return errArt(`range: need x,y,w`)
	}
	x, y, w := vals[0], vals[1], vals[2]
	if !inRange(w, 10, 90) {
		return errArt(`range: w out of range`)
	}

	label := ""
	if strings.HasPrefix(attrs, "?") {
		for _, kv := range strings.Split(attrs[1:], "&") {
			if kv == "" {
				continue
			}
			eq := strings.Index(kv, "=")
			if eq < 0 {
				return errArt(`range: bad attr ` + kv)
			}
			k, v := kv[:eq], kv[eq+1:]
			switch k {
			case "label":
				label = v
			default:
				return errArt(`range: unknown attr ` + k)
			}
		}
	}

	var b strings.Builder
	b.WriteString(`<svg viewBox='0 0 100 60' style='width:80%;height:auto'>`)
	// Axis rule.
	fmt.Fprintf(&b,
		`<line x1='%.2f' y1='%.2f' x2='%.2f' y2='%.2f' stroke='var(--ink)' stroke-width='0.8'/>`,
		x, y, x+w, y,
	)
	// End tick marks.
	fmt.Fprintf(&b,
		`<line x1='%.2f' y1='%.2f' x2='%.2f' y2='%.2f' stroke='var(--ink)' stroke-width='0.8'/>`,
		x, y-2, x, y+2,
	)
	fmt.Fprintf(&b,
		`<line x1='%.2f' y1='%.2f' x2='%.2f' y2='%.2f' stroke='var(--ink)' stroke-width='0.8'/>`,
		x+w, y-2, x+w, y+2,
	)
	// LO label (left end, above).
	fmt.Fprintf(&b,
		`<text x='%.2f' y='%.2f' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' font-weight='600' fill='var(--ink)'>%s</text>`,
		x, y-3.5, escapeXML(loStr),
	)
	// HI label.
	fmt.Fprintf(&b,
		`<text x='%.2f' y='%.2f' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' font-weight='600' fill='var(--ink)'>%s</text>`,
		x+w, y-3.5, escapeXML(hiStr),
	)
	// Optional caption.
	if label != "" {
		fmt.Fprintf(&b,
			`<text x='%.2f' y='%.2f' text-anchor='middle' font-family='IBM Plex Mono' font-size='6' fill='var(--ink-2)'>%s</text>`,
			x+w/2, y+8, escapeXML(label),
		)
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// ──────────────────────────────────────────────────────────────────────
// 6. baseline:VALUE@x,y,w[?label=TEXT]
// ──────────────────────────────────────────────────────────────────────

func dslBaseline(args string) string {
	at := strings.Index(args, "@")
	if at < 0 {
		return errArt(`baseline: need VALUE@x,y,w`)
	}
	valStr := strings.TrimSpace(args[:at])
	rest := args[at+1:]
	if valStr == "" {
		return errArt(`baseline: empty value`)
	}
	if _, err := strconv.ParseFloat(valStr, 64); err != nil {
		return errArt(`baseline: value not a number`)
	}
	var attrs string
	if q := strings.Index(rest, "?"); q >= 0 {
		attrs = rest[q:]
		rest = rest[:q]
	}
	vals, ok := parseFloatList(rest)
	if !ok || len(vals) != 3 {
		return errArt(`baseline: need x,y,w`)
	}
	x, y, w := vals[0], vals[1], vals[2]
	if !inRange(w, 10, 90) {
		return errArt(`baseline: w out of range`)
	}

	label := ""
	if strings.HasPrefix(attrs, "?") {
		for _, kv := range strings.Split(attrs[1:], "&") {
			if kv == "" {
				continue
			}
			eq := strings.Index(kv, "=")
			if eq < 0 {
				return errArt(`baseline: bad attr ` + kv)
			}
			k, v := kv[:eq], kv[eq+1:]
			switch k {
			case "label":
				label = v
			default:
				return errArt(`baseline: unknown attr ` + k)
			}
		}
	}

	// "50 · p99" style — value first then optional label.
	caption := valStr
	if label != "" {
		caption = valStr + " · " + label
	}

	return fmt.Sprintf(
		`<svg viewBox='0 0 100 60' style='width:80%%;height:auto'><line x1='%.2f' y1='%.2f' x2='%.2f' y2='%.2f' stroke='var(--accent)' stroke-width='1' stroke-dasharray='3 2'/><text x='%.2f' y='%.2f' text-anchor='end' font-family='IBM Plex Mono' font-size='6' font-weight='600' fill='var(--accent)'>%s</text></svg>`,
		x, y, x+w, y,
		x+w, y-1.5, escapeXML(caption),
	)
}

// ──────────────────────────────────────────────────────────────────────
// 7. times:N[?label=TEXT]
// ──────────────────────────────────────────────────────────────────────

func dslTimes(args string) string {
	var attrs string
	rest := args
	if q := strings.Index(args, "?"); q >= 0 {
		attrs = args[q:]
		rest = args[:q]
	}
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return errArt(`times: need N`)
	}
	n, err := strconv.Atoi(rest)
	if err != nil || n < 1 || n > 999 {
		return errArt(`times: bad N`)
	}

	label := ""
	if strings.HasPrefix(attrs, "?") {
		for _, kv := range strings.Split(attrs[1:], "&") {
			if kv == "" {
				continue
			}
			eq := strings.Index(kv, "=")
			if eq < 0 {
				return errArt(`times: bad attr ` + kv)
			}
			k, v := kv[:eq], kv[eq+1:]
			switch k {
			case "label":
				label = v
			default:
				return errArt(`times: unknown attr ` + k)
			}
		}
	}

	// Centered "×N" — big mono glyph. Optional small caption below.
	var b strings.Builder
	b.WriteString(`<svg viewBox='0 0 100 60' style='width:80%;height:auto'>`)
	mainY := 36.0
	if label == "" {
		mainY = 38.0
	}
	fmt.Fprintf(&b,
		`<text x='50' y='%.2f' text-anchor='middle' font-family='IBM Plex Mono' font-size='24' font-weight='700' fill='var(--accent)'>×%d</text>`,
		mainY, n,
	)
	if label != "" {
		fmt.Fprintf(&b,
			`<text x='50' y='52' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' fill='var(--ink-2)'>%s</text>`,
			escapeXML(label),
		)
	}
	b.WriteString(`</svg>`)
	return b.String()
}
