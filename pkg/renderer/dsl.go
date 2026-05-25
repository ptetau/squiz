package renderer

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// resolveDSL dispatches on the DSL prefix and renders to themed SVG.
// All SVGs render at viewBox 0 0 100 60 (the option-art slot's aspect)
// so they slot in identically alongside named/raw forms.
func resolveDSL(s string) (svg string, hidden bool) {
	colon := strings.Index(s, ":")
	prefix := s[:colon]
	args := s[colon+1:]
	switch prefix {
	case "grid":
		return dslGrid(args), false
	case "spark":
		return dslSpark(args), false
	case "bars":
		return dslBars(args), false
	case "swatches":
		return dslSwatches(args), false
	case "pills":
		return dslPills(args), false
	case "sample":
		return dslSample(args), false
	case "circle-pack":
		return dslCirclePack(args), false
	default:
		return errArt("unknown dsl: " + prefix), false
	}
}

func errArt(msg string) string {
	return fmt.Sprintf(`<svg viewBox='0 0 100 60' style='width:78%%;height:auto'><rect x='4' y='4' width='92' height='52' fill='none' stroke='var(--accent)' stroke-width='1' stroke-dasharray='3 2'/><text x='50' y='34' text-anchor='middle' font-family='IBM Plex Mono' font-size='8' fill='var(--accent)'>%s</text></svg>`, msg)
}

// ──────────────────────────────────────────────────────────────────────
// grid:NxM[@RATE]   — heatmap-style grid, RATE in [0,1] (default 0.5)
// ──────────────────────────────────────────────────────────────────────

func dslGrid(args string) string {
	rate := 0.5
	if i := strings.Index(args, "@"); i >= 0 {
		if v, err := strconv.ParseFloat(args[i+1:], 64); err == nil {
			rate = v
		}
		args = args[:i]
	}
	parts := strings.Split(args, "x")
	if len(parts) != 2 {
		return errArt("grid: need NxM")
	}
	rows, e1 := strconv.Atoi(parts[0])
	cols, e2 := strconv.Atoi(parts[1])
	if e1 != nil || e2 != nil || rows <= 0 || cols <= 0 || rows > 50 || cols > 50 {
		return errArt("grid: bad dims")
	}
	// Square cells, centered in the 100×60 viewBox.
	cellW := 80.0 / float64(cols)
	cellH := 48.0 / float64(rows)
	if cellH < cellW {
		cellW = cellH
	} else {
		cellH = cellW
	}
	startX := (100.0 - float64(cols)*cellW) / 2
	startY := (60.0 - float64(rows)*cellH) / 2

	var b strings.Builder
	b.WriteString(`<svg viewBox='0 0 100 60' style='width:80%;height:auto'>`)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			// Deterministic pseudo-random in [0,1) from (r,c).
			s := math.Sin(float64(r*cols+c)*12.9898+78.233) * 43758.5453
			rnd := s - math.Floor(s)
			x := startX + float64(c)*cellW
			y := startY + float64(r)*cellH
			w := cellW - 0.6
			h := cellH - 0.6
			switch {
			case rnd < rate*0.6:
				fmt.Fprintf(&b, `<rect x='%.2f' y='%.2f' width='%.2f' height='%.2f' fill='var(--accent)'/>`, x, y, w, h)
			case rnd < rate:
				fmt.Fprintf(&b, `<rect x='%.2f' y='%.2f' width='%.2f' height='%.2f' fill='var(--accent)' opacity='0.4'/>`, x, y, w, h)
			default:
				fmt.Fprintf(&b, `<rect x='%.2f' y='%.2f' width='%.2f' height='%.2f' fill='none' stroke='var(--rule-2)' stroke-width='0.4'/>`, x, y, w, h)
			}
		}
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// ──────────────────────────────────────────────────────────────────────
// spark:[V,V,V,...]   — sparkline from a list of numbers
// ──────────────────────────────────────────────────────────────────────

func dslSpark(args string) string {
	vals, ok := parseNumList(args)
	if !ok || len(vals) < 2 {
		return errArt("spark: need [v,v,...]")
	}
	min, max := vals[0], vals[0]
	for _, v := range vals {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	rng := max - min
	if rng == 0 {
		rng = 1
	}
	const w, h = 100.0, 60.0
	const padX, padY = 8.0, 10.0
	usableW := w - 2*padX
	usableH := h - 2*padY

	var path, area strings.Builder
	for i, v := range vals {
		x := padX + (float64(i)/float64(len(vals)-1))*usableW
		y := padY + usableH - ((v-min)/rng)*usableH
		if i == 0 {
			fmt.Fprintf(&path, `M %.2f %.2f`, x, y)
			fmt.Fprintf(&area, `M %.2f %.2f`, x, y)
		} else {
			fmt.Fprintf(&path, ` L %.2f %.2f`, x, y)
			fmt.Fprintf(&area, ` L %.2f %.2f`, x, y)
		}
	}
	fmt.Fprintf(&area, ` L %.2f %.2f L %.2f %.2f Z`, w-padX, h-padY, padX, h-padY)

	lastIdx := len(vals) - 1
	lastX := padX + usableW
	lastY := padY + usableH - ((vals[lastIdx]-min)/rng)*usableH

	return fmt.Sprintf(
		`<svg viewBox='0 0 100 60' style='width:82%%;height:auto'><path d='%s' fill='var(--accent)' opacity='0.15'/><path d='%s' fill='none' stroke='var(--accent)' stroke-width='1.6' stroke-linejoin='round'/><circle cx='%.2f' cy='%.2f' r='2' fill='var(--accent)'/></svg>`,
		area.String(), path.String(), lastX, lastY,
	)
}

// ──────────────────────────────────────────────────────────────────────
// bars:[V,V,V,...]   — bar chart from a list of numbers
// ──────────────────────────────────────────────────────────────────────

func dslBars(args string) string {
	vals, ok := parseNumList(args)
	if !ok || len(vals) < 1 {
		return errArt("bars: need [v,v,...]")
	}
	max := vals[0]
	for _, v := range vals {
		if v > max {
			max = v
		}
	}
	if max == 0 {
		max = 1
	}
	const w, h = 100.0, 60.0
	const padX, padY = 10.0, 8.0
	usableW := w - 2*padX
	usableH := h - 2*padY
	barGap := 2.0
	barW := (usableW - barGap*float64(len(vals)-1)) / float64(len(vals))

	var b strings.Builder
	b.WriteString(`<svg viewBox='0 0 100 60' style='width:80%;height:auto'>`)
	for i, v := range vals {
		barH := (v / max) * usableH
		x := padX + float64(i)*(barW+barGap)
		y := padY + usableH - barH
		fmt.Fprintf(&b, `<rect x='%.2f' y='%.2f' width='%.2f' height='%.2f' fill='var(--accent)' rx='1'/>`, x, y, barW, barH)
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// ──────────────────────────────────────────────────────────────────────
// swatches:#A,#B,#C   — color palette horizontal bars
// ──────────────────────────────────────────────────────────────────────

func dslSwatches(args string) string {
	if strings.TrimSpace(args) == "" {
		return errArt("swatches: need #a,#b,...")
	}
	colors := strings.Split(args, ",")
	const w, h = 100.0, 60.0
	const padX, padY = 8.0, 14.0
	usableW := w - 2*padX
	usableH := h - 2*padY
	gap := 2.0
	cw := (usableW - gap*float64(len(colors)-1)) / float64(len(colors))

	var b strings.Builder
	b.WriteString(`<svg viewBox='0 0 100 60' style='width:80%;height:auto'>`)
	for i, c := range colors {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		x := padX + float64(i)*(cw+gap)
		fmt.Fprintf(&b, `<rect x='%.2f' y='%.2f' width='%.2f' height='%.2f' fill='%s' stroke='var(--ink)' stroke-width='0.3'/>`, x, padY, cw, usableH, c)
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// ──────────────────────────────────────────────────────────────────────
// pills:LABEL*|LABEL|LABEL*   — chip row, `*` suffix = active
// ──────────────────────────────────────────────────────────────────────

func dslPills(args string) string {
	if strings.TrimSpace(args) == "" {
		return errArt("pills: need a|b|c")
	}
	parts := strings.Split(args, "|")
	const w, h = 100.0, 60.0
	const padX, padY = 6.0, 4.0
	gap := 4.0

	// Compute total natural width to scale chips into the viewBox.
	type chip struct {
		text   string
		active bool
		w      float64
	}
	chips := make([]chip, 0, len(parts))
	totalW := 0.0
	for _, p := range parts {
		p = strings.TrimSpace(p)
		active := strings.HasSuffix(p, "*")
		text := strings.TrimSuffix(p, "*")
		// Heuristic width: ~5px per char + 10px padding, capped.
		cw := float64(len(text))*4.5 + 10
		if cw > 40 {
			cw = 40
		}
		chips = append(chips, chip{text, active, cw})
		totalW += cw
	}
	totalW += gap * float64(len(chips)-1)
	usableW := w - 2*padX
	scale := 1.0
	if totalW > usableW {
		scale = usableW / totalW
	}

	// Stack rows of chips if there are many; for now keep one row centered.
	x := (w - totalW*scale) / 2
	y := h/2 - 9

	var b strings.Builder
	b.WriteString(`<svg viewBox='0 0 100 60' style='width:80%;height:auto'>`)
	for _, c := range chips {
		cw := c.w * scale
		fill := "none"
		stroke := "var(--ink)"
		txtFill := "var(--ink)"
		if c.active {
			fill = "var(--accent)"
			stroke = "var(--accent)"
			txtFill = "var(--bg)"
		}
		fmt.Fprintf(&b, `<rect x='%.2f' y='%.2f' width='%.2f' height='14' rx='7' fill='%s' stroke='%s' stroke-width='0.8'/>`, x, y, cw, fill, stroke)
		fmt.Fprintf(&b, `<text x='%.2f' y='%.2f' text-anchor='middle' font-family='IBM Plex Mono' font-size='6.5' fill='%s' font-weight='600'>%s</text>`,
			x+cw/2, y+9, txtFill, escapeXML(c.text))
		x += cw + gap*scale
	}
	b.WriteString(`</svg>`)
	_ = padY
	return b.String()
}

// ──────────────────────────────────────────────────────────────────────
// sample:"TEXT"[@FONT]   — styled text sample. FONT: serif|sans|mono
// ──────────────────────────────────────────────────────────────────────

func dslSample(args string) string {
	font := "sans"
	if i := strings.LastIndex(args, "@"); i >= 0 && i > 0 && strings.HasSuffix(args[:i], `"`) {
		font = strings.TrimSpace(args[i+1:])
		args = args[:i]
	}
	args = strings.TrimSpace(args)
	if !strings.HasPrefix(args, `"`) || !strings.HasSuffix(args, `"`) || len(args) < 2 {
		return errArt(`sample: need "text"`)
	}
	text := args[1 : len(args)-1]

	var family, style, weight string
	switch font {
	case "serif":
		family = "IBM Plex Serif"
		style = "italic"
		weight = "400"
	case "mono":
		family = "IBM Plex Mono"
		style = "normal"
		weight = "500"
	default: // sans
		family = "IBM Plex Sans"
		style = "normal"
		weight = "500"
	}

	// Naive line-wrap at ~26 chars per line, max 3 lines.
	const maxChars = 26
	lines := []string{}
	rest := text
	for len(rest) > 0 && len(lines) < 3 {
		if len(rest) <= maxChars {
			lines = append(lines, rest)
			break
		}
		// Find last space ≤ maxChars
		cut := maxChars
		for j := maxChars; j > 0; j-- {
			if rest[j-1] == ' ' {
				cut = j - 1
				break
			}
		}
		lines = append(lines, strings.TrimSpace(rest[:cut]))
		rest = strings.TrimSpace(rest[cut:])
	}
	if len(rest) > 0 && len(lines) == 3 {
		lines[2] = strings.TrimRight(lines[2], ". ") + "…"
	}

	startY := 30 - float64(len(lines)-1)*5
	var b strings.Builder
	b.WriteString(`<svg viewBox='0 0 100 60' style='width:82%;height:auto'>`)
	for i, ln := range lines {
		fmt.Fprintf(&b, `<text x='50' y='%.2f' text-anchor='middle' font-family='%s' font-style='%s' font-weight='%s' font-size='9' fill='var(--ink)'>%s</text>`,
			startY+float64(i)*10, family, style, weight, escapeXML(ln))
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// ──────────────────────────────────────────────────────────────────────
// circle-pack:N   — N circles arranged organically (deterministic)
// ──────────────────────────────────────────────────────────────────────

func dslCirclePack(args string) string {
	n, err := strconv.Atoi(strings.TrimSpace(args))
	if err != nil || n <= 0 || n > 50 {
		return errArt("circle-pack: need N")
	}
	// Deterministic positions on a Poisson-ish grid.
	const w, h = 100.0, 60.0
	var b strings.Builder
	b.WriteString(`<svg viewBox='0 0 100 60' style='width:80%;height:auto'>`)
	for i := 0; i < n; i++ {
		seed1 := math.Sin(float64(i)*12.9898) * 43758.5453
		seed2 := math.Sin(float64(i)*4.1414+1.7) * 43758.5453
		seed3 := math.Sin(float64(i)*9.81-0.3) * 43758.5453
		rnd1 := seed1 - math.Floor(seed1)
		rnd2 := seed2 - math.Floor(seed2)
		rnd3 := seed3 - math.Floor(seed3)
		x := 10 + rnd1*(w-20)
		y := 8 + rnd2*(h-16)
		r := 2 + rnd3*4
		var fill string
		if i == 0 {
			fill = `fill='var(--accent)'`
		} else if rnd3 < 0.3 {
			fill = `fill='var(--ink)' opacity='0.45'`
		} else {
			fill = `fill='var(--ink)' opacity='0.7'`
		}
		fmt.Fprintf(&b, `<circle cx='%.2f' cy='%.2f' r='%.2f' %s/>`, x, y, r, fill)
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// ──────────────────────────────────────────────────────────────────────
// helpers
// ──────────────────────────────────────────────────────────────────────

// parseNumList parses "[3,5,4,7]" or "3,5,4,7" → []float64
func parseNumList(s string) ([]float64, bool) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	parts := strings.Split(s, ",")
	out := make([]float64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return nil, false
		}
		out = append(out, v)
	}
	return out, len(out) > 0
}

func escapeXML(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return r.Replace(s)
}
