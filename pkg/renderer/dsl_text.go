package renderer

import (
	"fmt"
	"strconv"
	"strings"
)

// dslText parses `text:"line 1\nline 2"@font?size=14&align=center&weight=700&color=accent`
//
// Mandatory: text body (must be in double quotes).
// Optional @font: mono | serif | sans (default sans).
// Optional ?attrs:
//
//	size:   integer 6..36 (default 14)
//	align:  left | center | right (default left)
//	weight: 300..700 (default 400)
//	color:  ink | ink-2 | ink-3 | accent | rule | rule-2  (default ink)
//
// Multi-line: `\n` in the body becomes a new line in the SVG.
// Emits an errArt placeholder for: missing quotes, empty body, unknown @font,
// out-of-range size, unknown attr key, unknown color name.
func dslText(args string) string {
	// 1) Body must start with a double quote.
	args = strings.TrimSpace(args)
	if !strings.HasPrefix(args, `"`) {
		return errArt(`text: missing quotes`)
	}
	// 2) Find the matching closing quote. The body is literal — no escape
	//    support; whatever follows the first close-quote is suffix metadata.
	closeIdx := strings.Index(args[1:], `"`)
	if closeIdx < 0 {
		return errArt(`text: missing quotes`)
	}
	body := args[1 : 1+closeIdx]
	suffix := args[2+closeIdx:]
	if body == "" {
		return errArt(`text: empty body`)
	}

	// 3) Optional @font then optional ?attrs.
	font := "sans"
	if strings.HasPrefix(suffix, "@") {
		// Trim the leading @ and read until either '?' or end.
		rest := suffix[1:]
		qIdx := strings.Index(rest, "?")
		if qIdx >= 0 {
			font = rest[:qIdx]
			suffix = "?" + rest[qIdx+1:]
		} else {
			font = rest
			suffix = ""
		}
	}
	var family string
	switch font {
	case "sans":
		family = "IBM Plex Sans"
	case "mono":
		family = "IBM Plex Mono"
	case "serif":
		family = "IBM Plex Serif"
	default:
		return errArt(`text: unknown @font ` + font)
	}

	// 4) Defaults.
	size := 14
	align := "left"
	weight := 400
	color := "ink"

	// 5) Parse attrs.
	if strings.HasPrefix(suffix, "?") {
		for _, kv := range strings.Split(suffix[1:], "&") {
			if kv == "" {
				continue
			}
			eq := strings.Index(kv, "=")
			if eq < 0 {
				return errArt(`text: bad attr ` + kv)
			}
			k := kv[:eq]
			v := kv[eq+1:]
			switch k {
			case "size":
				n, err := strconv.Atoi(v)
				if err != nil || n < 6 || n > 36 {
					return errArt(`text: size out of range`)
				}
				size = n
			case "align":
				if v != "left" && v != "center" && v != "right" {
					return errArt(`text: bad align ` + v)
				}
				align = v
			case "weight":
				n, err := strconv.Atoi(v)
				if err != nil || n < 300 || n > 700 {
					return errArt(`text: weight out of range`)
				}
				weight = n
			case "color":
				switch v {
				case "ink", "ink-2", "ink-3", "accent", "rule", "rule-2":
					color = v
				default:
					return errArt(`text: unknown color ` + v)
				}
			default:
				return errArt(`text: unknown attr ` + k)
			}
		}
	}

	// 6) Split on literal "\n" (two chars: backslash + n).
	lines := strings.Split(body, `\n`)

	// 7) Pick x + text-anchor for alignment.
	var x int
	var anchor string
	switch align {
	case "center":
		x = 50
		anchor = "middle"
	case "right":
		x = 94
		anchor = "end"
	default: // left
		x = 6
		anchor = "start"
	}

	// 8) Vertical layout: center the block of lines around y=30.
	lineH := float64(size) * 1.2
	startY := 30.0 - (float64(len(lines)-1) * lineH / 2) + float64(size)*0.35

	var b strings.Builder
	b.WriteString(`<svg viewBox='0 0 100 60' style='width:80%;height:auto'>`)
	for i, ln := range lines {
		y := startY + float64(i)*lineH
		fmt.Fprintf(&b,
			`<text x='%d' y='%.2f' text-anchor='%s' font-family='%s' font-size='%d' font-weight='%d' fill='var(--%s)'>%s</text>`,
			x, y, anchor, family, size, weight, color, escapeXML(ln),
		)
	}
	b.WriteString(`</svg>`)
	return b.String()
}
