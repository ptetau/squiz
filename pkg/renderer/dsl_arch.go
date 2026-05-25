package renderer

import (
	"fmt"
	"strings"
)

// dslFlow parses `flow:[a,b,c]` or `flow:[client?icon=user,api?icon=api,db?icon=database]`
//
// Each comma-separated token is a labeled box, connected left-to-right
// with arrows. Optional `?icon=<name>` attaches an arch icon inside the box.
// Auto-scales horizontally to fit viewBox 0 0 100 60.
// Emits errArt for: empty list, malformed token, unknown icon name.
func dslFlow(args string) string {
	args = strings.TrimSpace(args)
	if !strings.HasPrefix(args, "[") || !strings.HasSuffix(args, "]") || len(args) < 2 {
		return errArt(`flow: need [a,b,...]`)
	}
	inner := strings.TrimSpace(args[1 : len(args)-1])
	if inner == "" {
		return errArt(`flow: empty list`)
	}

	type node struct {
		label string
		icon  string // empty if none
	}
	toks := strings.Split(inner, ",")
	nodes := make([]node, 0, len(toks))
	for _, t := range toks {
		t = strings.TrimSpace(t)
		if t == "" {
			return errArt(`flow: empty token`)
		}
		n := node{}
		if q := strings.Index(t, "?"); q >= 0 {
			n.label = strings.TrimSpace(t[:q])
			// Only one supported attr: icon=<name>
			attr := t[q+1:]
			if !strings.HasPrefix(attr, "icon=") {
				return errArt(`flow: bad token ` + t)
			}
			n.icon = attr[len("icon="):]
			if n.icon == "" {
				return errArt(`flow: bad token ` + t)
			}
			if ArchIcon(n.icon) == "" {
				return errArt(`flow: unknown icon ` + n.icon)
			}
		} else {
			n.label = t
		}
		if n.label == "" {
			return errArt(`flow: bad token ` + t)
		}
		nodes = append(nodes, n)
	}

	const w, h = 100.0, 60.0
	const padX = 4.0
	const arrowW = 6.0
	usableW := w - 2*padX
	n := len(nodes)
	// Total width = n*boxW + (n-1)*arrowW. Solve for boxW.
	boxW := (usableW - float64(n-1)*arrowW) / float64(n)
	if boxW > 26 {
		boxW = 26
	}
	if boxW < 10 {
		boxW = 10
	}
	totalW := boxW*float64(n) + arrowW*float64(n-1)
	startX := (w - totalW) / 2

	// Box dimensions.
	boxH := 32.0
	boxY := (h - boxH) / 2

	var b strings.Builder
	b.WriteString(`<svg viewBox='0 0 100 60' style='width:90%;height:auto'>`)
	for i, nd := range nodes {
		x := startX + float64(i)*(boxW+arrowW)
		// Box.
		fmt.Fprintf(&b,
			`<rect x='%.2f' y='%.2f' width='%.2f' height='%.2f' rx='2' fill='var(--surface)' stroke='var(--accent)' stroke-width='1.4'/>`,
			x, boxY, boxW, boxH,
		)
		// Optional icon, scaled to fit upper portion.
		if nd.icon != "" {
			g := ArchIcon(nd.icon)
			// Source icon sits in a ~40x40 region centered on (50,30). Map
			// that region into the upper ~16px of the box.
			//   scale = 16/40 = 0.40
			//   translate so (50,30) maps to (x+boxW/2, boxY+10)
			scale := 0.40
			tx := x + boxW/2 - 50*scale
			ty := boxY + 10 - 30*scale
			fmt.Fprintf(&b,
				`<g transform='translate(%.2f,%.2f) scale(%.3f)'>%s</g>`,
				tx, ty, scale, g,
			)
			// Label below icon.
			fmt.Fprintf(&b,
				`<text x='%.2f' y='%.2f' text-anchor='middle' font-family='IBM Plex Mono' font-size='5.5' font-weight='600' fill='var(--ink)'>%s</text>`,
				x+boxW/2, boxY+boxH-4, escapeXML(truncateLabel(nd.label, boxW)),
			)
		} else {
			// Label centered in box.
			fmt.Fprintf(&b,
				`<text x='%.2f' y='%.2f' text-anchor='middle' font-family='IBM Plex Mono' font-size='6.5' font-weight='600' fill='var(--ink)'>%s</text>`,
				x+boxW/2, boxY+boxH/2+2.5, escapeXML(truncateLabel(nd.label, boxW)),
			)
		}
		// Arrow to next box.
		if i < n-1 {
			ax := x + boxW
			ay := h / 2
			fmt.Fprintf(&b,
				`<line x1='%.2f' y1='%.2f' x2='%.2f' y2='%.2f' stroke='var(--ink)' stroke-width='1'/>`,
				ax, ay, ax+arrowW-1, ay,
			)
			fmt.Fprintf(&b,
				`<polygon points='%.2f,%.2f %.2f,%.2f %.2f,%.2f' fill='var(--ink)'/>`,
				ax+arrowW-1, ay, ax+arrowW-4, ay-2, ax+arrowW-4, ay+2,
			)
		}
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// dslBox parses `box:label?icon=server`
//
// Single named box with an optional embedded arch icon. Centered in viewBox.
// Emits errArt for: empty label, unknown icon name, unknown attr key.
func dslBox(args string) string {
	args = strings.TrimSpace(args)
	if args == "" {
		return errArt(`box: empty label`)
	}
	label := args
	icon := ""
	if q := strings.Index(args, "?"); q >= 0 {
		label = strings.TrimSpace(args[:q])
		attrs := args[q+1:]
		for _, kv := range strings.Split(attrs, "&") {
			if kv == "" {
				continue
			}
			eq := strings.Index(kv, "=")
			if eq < 0 {
				return errArt(`box: bad attr ` + kv)
			}
			k := kv[:eq]
			v := kv[eq+1:]
			switch k {
			case "icon":
				if ArchIcon(v) == "" {
					return errArt(`box: unknown icon ` + v)
				}
				icon = v
			default:
				return errArt(`box: unknown attr ` + k)
			}
		}
	}
	if label == "" {
		return errArt(`box: empty label`)
	}

	const boxX, boxY, boxW, boxH = 18.0, 10.0, 64.0, 40.0
	var b strings.Builder
	b.WriteString(`<svg viewBox='0 0 100 60' style='width:80%;height:auto'>`)
	fmt.Fprintf(&b,
		`<rect x='%.2f' y='%.2f' width='%.2f' height='%.2f' rx='3' fill='var(--surface)' stroke='var(--accent)' stroke-width='1.6'/>`,
		boxX, boxY, boxW, boxH,
	)
	if icon != "" {
		g := ArchIcon(icon)
		// Source icon sits in ~40x40 region centered on (50,30); scale to
		// ~24px and place in the upper portion of the box.
		scale := 0.60
		tx := boxX + boxW/2 - 50*scale
		ty := boxY + 8 - 30*scale
		fmt.Fprintf(&b,
			`<g transform='translate(%.2f,%.2f) scale(%.3f)'>%s</g>`,
			tx, ty, scale, g,
		)
		fmt.Fprintf(&b,
			`<text x='%.2f' y='%.2f' text-anchor='middle' font-family='IBM Plex Mono' font-size='8' font-weight='600' fill='var(--ink)'>%s</text>`,
			boxX+boxW/2, boxY+boxH-6, escapeXML(label),
		)
	} else {
		fmt.Fprintf(&b,
			`<text x='%.2f' y='%.2f' text-anchor='middle' font-family='IBM Plex Mono' font-size='10' font-weight='600' fill='var(--ink)'>%s</text>`,
			boxX+boxW/2, boxY+boxH/2+4, escapeXML(label),
		)
	}
	b.WriteString(`</svg>`)
	return b.String()
}

// dslArrow parses `arrow:"label"?dir=right`
//
// A standalone labeled arrow glyph, useful as a small inline visualization
// of a transition. dir = right (default) | down | up | left.
// Emits errArt for: missing quotes, empty label, unknown dir.
func dslArrow(args string) string {
	args = strings.TrimSpace(args)
	if !strings.HasPrefix(args, `"`) {
		return errArt(`arrow: missing quotes`)
	}
	closeIdx := strings.Index(args[1:], `"`)
	if closeIdx < 0 {
		return errArt(`arrow: missing quotes`)
	}
	label := args[1 : 1+closeIdx]
	if label == "" {
		return errArt(`arrow: empty label`)
	}
	suffix := args[2+closeIdx:]

	dir := "right"
	if strings.HasPrefix(suffix, "?") {
		for _, kv := range strings.Split(suffix[1:], "&") {
			if kv == "" {
				continue
			}
			eq := strings.Index(kv, "=")
			if eq < 0 {
				return errArt(`arrow: bad attr ` + kv)
			}
			k := kv[:eq]
			v := kv[eq+1:]
			switch k {
			case "dir":
				switch v {
				case "right", "left", "up", "down":
					dir = v
				default:
					return errArt(`arrow: unknown dir ` + v)
				}
			default:
				return errArt(`arrow: unknown attr ` + k)
			}
		}
	}

	var line, head, labelX, labelY, labelAnchor string
	switch dir {
	case "right":
		line = `<line x1='18' y1='30' x2='78' y2='30' stroke='var(--accent)' stroke-width='1.8'/>`
		head = `<polygon points='82,30 76,27 76,33' fill='var(--accent)'/>`
		labelX, labelY, labelAnchor = "50", "22", "middle"
	case "left":
		line = `<line x1='22' y1='30' x2='82' y2='30' stroke='var(--accent)' stroke-width='1.8'/>`
		head = `<polygon points='18,30 24,27 24,33' fill='var(--accent)'/>`
		labelX, labelY, labelAnchor = "50", "22", "middle"
	case "down":
		line = `<line x1='50' y1='8' x2='50' y2='48' stroke='var(--accent)' stroke-width='1.8'/>`
		head = `<polygon points='50,52 47,46 53,46' fill='var(--accent)'/>`
		labelX, labelY, labelAnchor = "58", "32", "start"
	case "up":
		line = `<line x1='50' y1='12' x2='50' y2='52' stroke='var(--accent)' stroke-width='1.8'/>`
		head = `<polygon points='50,8 47,14 53,14' fill='var(--accent)'/>`
		labelX, labelY, labelAnchor = "58", "32", "start"
	}

	return fmt.Sprintf(
		`<svg viewBox='0 0 100 60' style='width:80%%;height:auto'>%s%s<text x='%s' y='%s' text-anchor='%s' font-family='IBM Plex Mono' font-size='8' font-weight='600' fill='var(--ink)'>%s</text></svg>`,
		line, head, labelX, labelY, labelAnchor, escapeXML(label),
	)
}

// truncateLabel shortens a label heuristically to fit a box width (px).
// ~2.2px per char at font-size 6.5 → maxChars ≈ boxW/2.2.
func truncateLabel(s string, boxW float64) string {
	maxChars := int(boxW / 2.2)
	if maxChars < 3 {
		maxChars = 3
	}
	if len(s) <= maxChars {
		return s
	}
	if maxChars <= 1 {
		return s[:maxChars]
	}
	return s[:maxChars-1] + "…"
}
