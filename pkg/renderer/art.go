package renderer

import (
	"fmt"
	"regexp"
	"strings"
)

// RenderArt resolves an option's `art` field to an SVG fragment.
//
// Returns (svg, hidden):
//   - hidden=true means "render no art slot at all" (the caller should
//     conditionally drop the <div class="option-art"> wrapper).
//   - hidden=false with empty svg means "render the slot empty" (shouldn't
//     happen with well-formed input but is harmless if it does).
//
// letterIdx is the 0-based option position within its squiz, used to pick
// the auto-art pattern when `art` is omitted.
func RenderArt(art string, letterIdx int) (svg string, hidden bool) {
	art = strings.TrimSpace(art)

	switch {
	case art == "":
		// Default: per-letter abstract pattern.
		return autoArt(letterIdx), false
	case art == "none":
		return "", true
	case strings.HasPrefix(art, "<svg"):
		// Raw SVG, passed through.
		return art, false
	case strings.HasPrefix(art, "wf:"):
		return resolveNamed(strings.TrimPrefix(art, "wf:"))
	case dslPrefixRE.MatchString(art):
		return resolveDSL(art)
	default:
		// Unknown form — treat as raw text in a labeled placeholder so the
		// author notices their typo without crashing the render.
		safe := strings.ReplaceAll(strings.ReplaceAll(art, "<", "&lt;"), ">", "&gt;")
		return fmt.Sprintf(`<svg viewBox='0 0 100 60' style='width:78%%;height:auto'><rect x='4' y='4' width='92' height='52' fill='none' stroke='var(--accent)' stroke-width='1' stroke-dasharray='3 2'/><text x='50' y='28' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' fill='var(--ink-3)'>unknown art</text><text x='50' y='40' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' fill='var(--accent)'>%s</text></svg>`, safe), false
	}
}

// dslPrefixRE matches strings that start with a lowercase identifier
// followed by ':', e.g. "grid:7x7@0.55". Excludes raw SVG ("<svg") and
// the named-lib prefix ("wf:") which are handled separately.
var dslPrefixRE = regexp.MustCompile(`^[a-z][a-z0-9-]*:`)

func resolveNamed(name string) (svg string, hidden bool) {
	if entry, ok := WFLibrary[name]; ok {
		return entry, false
	}
	return fmt.Sprintf(`<svg viewBox='0 0 100 60' style='width:78%%;height:auto'><rect x='4' y='4' width='92' height='52' fill='none' stroke='var(--accent)' stroke-width='1' stroke-dasharray='3 2'/><text x='50' y='28' text-anchor='middle' font-family='IBM Plex Mono' font-size='7' fill='var(--ink-3)'>wf:</text><text x='50' y='40' text-anchor='middle' font-family='IBM Plex Mono' font-size='8' font-weight='600' fill='var(--accent)'>%s ?</text></svg>`, name), false
}

// ──────────────────────────────────────────────────────────────────────
// Auto-art: subtle per-letter pattern when `art` is omitted.
// ──────────────────────────────────────────────────────────────────────

var autoPatterns = []string{
	// A — hatched
	`<svg viewBox='0 0 100 60' style='width:78%;height:auto'><defs><pattern id='p-a' patternUnits='userSpaceOnUse' width='5' height='5' patternTransform='rotate(45)'><line x1='0' y1='0' x2='0' y2='5' stroke='var(--accent)' stroke-width='1' opacity='0.55'/></pattern></defs><rect x='6' y='6' width='88' height='48' fill='url(#p-a)' stroke='var(--rule-2)' stroke-width='0.5'/></svg>`,
	// B — dotted
	`<svg viewBox='0 0 100 60' style='width:78%;height:auto'><defs><pattern id='p-b' patternUnits='userSpaceOnUse' width='6' height='6'><circle cx='3' cy='3' r='1' fill='var(--accent)' opacity='0.55'/></pattern></defs><rect x='6' y='6' width='88' height='48' fill='url(#p-b)' stroke='var(--rule-2)' stroke-width='0.5'/></svg>`,
	// C — horizontal stripes
	`<svg viewBox='0 0 100 60' style='width:78%;height:auto'><defs><pattern id='p-c' patternUnits='userSpaceOnUse' width='5' height='5'><line x1='0' y1='0' x2='5' y2='0' stroke='var(--accent)' stroke-width='1' opacity='0.55'/></pattern></defs><rect x='6' y='6' width='88' height='48' fill='url(#p-c)' stroke='var(--rule-2)' stroke-width='0.5'/></svg>`,
	// D — small grid
	`<svg viewBox='0 0 100 60' style='width:78%;height:auto'><defs><pattern id='p-d' patternUnits='userSpaceOnUse' width='6' height='6'><path d='M 6 0 L 0 0 0 6' fill='none' stroke='var(--accent)' stroke-width='0.6' opacity='0.55'/></pattern></defs><rect x='6' y='6' width='88' height='48' fill='url(#p-d)' stroke='var(--rule-2)' stroke-width='0.5'/></svg>`,
	// E — cross-hatched
	`<svg viewBox='0 0 100 60' style='width:78%;height:auto'><defs><pattern id='p-e' patternUnits='userSpaceOnUse' width='6' height='6' patternTransform='rotate(45)'><line x1='0' y1='0' x2='0' y2='6' stroke='var(--accent)' stroke-width='0.6' opacity='0.5'/><line x1='3' y1='0' x2='3' y2='6' stroke='var(--accent)' stroke-width='0.6' opacity='0.5'/></pattern></defs><rect x='6' y='6' width='88' height='48' fill='url(#p-e)' stroke='var(--rule-2)' stroke-width='0.5'/></svg>`,
	// F — waves
	`<svg viewBox='0 0 100 60' style='width:78%;height:auto'><defs><pattern id='p-f' patternUnits='userSpaceOnUse' width='12' height='6'><path d='M 0 3 Q 3 0 6 3 T 12 3' fill='none' stroke='var(--accent)' stroke-width='0.8' opacity='0.55'/></pattern></defs><rect x='6' y='6' width='88' height='48' fill='url(#p-f)' stroke='var(--rule-2)' stroke-width='0.5'/></svg>`,
}

func autoArt(letterIdx int) string {
	if letterIdx < 0 {
		letterIdx = 0
	}
	return autoPatterns[letterIdx%len(autoPatterns)]
}

// LetterFor returns the Option's letter ("A", "B", …) for `Label` defaulting.
func LetterFor(idx int) string {
	if idx < 0 {
		return ""
	}
	// A–Z, then AA, AB, … (we'll never see this in practice).
	if idx < 26 {
		return string(rune('A' + idx))
	}
	return string(rune('A'+(idx/26)-1)) + string(rune('A'+(idx%26)))
}
