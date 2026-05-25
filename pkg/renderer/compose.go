package renderer

import (
	"fmt"
	"regexp"
	"strings"
)

// Composition support (v0.8.0): raw SVG can embed library + DSL primitives
// via <use href="wf:phone-card"/> / <use href="arch:server"/> /
// <use href="callout:..."/>. The renderer pre-processes raw SVG, finds
// these refs, inlines each referenced body as a <symbol> in a <defs>
// block, and rewrites hrefs to local fragment ids.

// NaturalBox describes the natural footprint of a library entry inside
// its 100×60 viewBox. width/height are in viewBox units; AspectRatio
// = width/height. Used by `squiz catalog --json` and by composing
// agents to pick sensible <use> bounding boxes.
type NaturalBox struct {
	W           float64 `json:"w"`
	H           float64 `json:"h"`
	AspectRatio float64 `json:"aspect"`
}

// DefaultNaturalBox is the conservative fallback for entries that
// haven't been measured: the full viewBox. Agents using <use> with
// these dimensions will get the full-size rendering; the inner content
// will center naturally if it's smaller.
var DefaultNaturalBox = NaturalBox{W: 100, H: 60, AspectRatio: 100.0 / 60.0}

// NaturalBoxOf returns the measured (or default) natural box for a
// library entry. Lookup order: WFNaturalBox → ArchNaturalBox →
// DefaultNaturalBox.
func NaturalBoxOf(name string) NaturalBox {
	if b, ok := WFNaturalBox[name]; ok {
		return b
	}
	if b, ok := ArchNaturalBox[name]; ok {
		return b
	}
	return DefaultNaturalBox
}

// ──────────────────────────────────────────────────────────────────────
// resolveUses — composition pipeline
// ──────────────────────────────────────────────────────────────────────

// useRefRE matches a complete <use ...> tag (self-closing or with </use>).
// We don't try to parse hrefs here — that's done with a separate regex on
// the matched tag — because attribute ordering inside <use ...> is free.
var useRefRE = regexp.MustCompile(`<use\b[^>]*/?>`)

// hrefAttrRE extracts the href value (double- or single-quoted) from a
// single <use ...> tag. Supports `href="…"`, `href='…'`, `xlink:href="…"`,
// `xlink:href='…'`.
var hrefAttrRE = regexp.MustCompile(`(?:xlink:)?href\s*=\s*(?:"([^"]*)"|'([^']*)')`)

// composableSpecRE matches the allowed art-spec prefixes that resolveUses
// will route through RenderArt. Anything else (e.g. `#local-defs`,
// `https://…`, plain ids) is left untouched.
var composableSpecRE = regexp.MustCompile(`^(?:wf|arch|grid|spark|bars|swatches|pills|sample|circle-pack|text|flow|box|arrow|callout|brace|divider|badge|range|baseline|times):`)

// svgOpenTagRE matches the opening <svg ...> tag of a raw SVG body so
// we can splice <defs> right after it. Greedy enough to handle attrs
// with quoted values; intentionally simple.
var svgOpenTagRE = regexp.MustCompile(`<svg\b[^>]*>`)

// viewBoxAttrRE pulls the viewBox value from a rendered SVG so the
// generated <symbol> inherits the same coordinate system.
var viewBoxAttrRE = regexp.MustCompile(`viewBox\s*=\s*(?:"([^"]*)"|'([^']*)')`)

// resolveUses preprocesses raw SVG, finds <use href="<art-spec>"/>
// references, inlines each as a <symbol> in a <defs> block, rewrites
// the hrefs to local fragment ids.
//
// No-op when the input contains no composable <use> refs — preserves
// byte-for-byte output for the common (non-composing) raw-SVG case.
// Failures (unknown wf:/arch:/dsl names) still emit a symbol so the
// <use> doesn't dangle; the symbol body is the same fail-soft placeholder
// RenderArt would have produced standalone.
func resolveUses(svg string) string {
	// Quick reject: nothing to do if there's no <use ... href=…> we care about.
	matches := useRefRE.FindAllStringIndex(svg, -1)
	if len(matches) == 0 {
		return svg
	}

	// Walk all <use> tags, identify the unique composable specs we need to
	// resolve, and remember the rewrite map (spec → sanitized id).
	type useHit struct {
		start, end int    // span of the full <use ...> tag in svg
		spec       string // the original href value (e.g. "wf:phone-card")
		id         string // sanitized symbol id (e.g. "wf-phone-card")
	}
	var hits []useHit
	uniqueSpecs := make(map[string]string) // spec → sanitized id, insertion-stable
	var specOrder []string                 // deterministic <defs> ordering

	for _, m := range matches {
		tag := svg[m[0]:m[1]]
		hAttr := hrefAttrRE.FindStringSubmatch(tag)
		if hAttr == nil {
			continue
		}
		href := hAttr[1]
		if href == "" {
			href = hAttr[2]
		}
		if !composableSpecRE.MatchString(href) {
			continue // local fragment or unrelated href, leave alone
		}
		id, seen := uniqueSpecs[href]
		if !seen {
			id = sanitizeSymbolID(href)
			uniqueSpecs[href] = id
			specOrder = append(specOrder, href)
		}
		hits = append(hits, useHit{start: m[0], end: m[1], spec: href, id: id})
	}

	if len(hits) == 0 {
		return svg
	}

	// Build the <defs> block of <symbol>s, in first-seen order.
	var defs strings.Builder
	defs.WriteString(`<defs>`)
	for _, spec := range specOrder {
		body, _ := RenderArt(spec, 0)
		inner, vb := stripSVGWrapper(body)
		if vb == "" {
			vb = "0 0 100 60"
		}
		fmt.Fprintf(&defs, `<symbol id="%s" viewBox="%s">%s</symbol>`,
			uniqueSpecs[spec], vb, inner)
	}
	defs.WriteString(`</defs>`)

	// Rewrite each captured <use> tag to point at the local fragment id,
	// preserving every other attribute (x, y, width, height, transform, …).
	// Walk back-to-front so byte offsets remain valid as we splice.
	rewritten := svg
	for i := len(hits) - 1; i >= 0; i-- {
		h := hits[i]
		oldTag := rewritten[h.start:h.end]
		newTag := rewriteHref(oldTag, "#"+h.id)
		rewritten = rewritten[:h.start] + newTag + rewritten[h.end:]
	}

	// Splice <defs> right after the opening <svg ...> tag. If no opening
	// <svg> is found (e.g. the input is a bare fragment), prepend the defs
	// block — the resulting HTML still works because <defs> is a no-op
	// container without rendering.
	loc := svgOpenTagRE.FindStringIndex(rewritten)
	if loc == nil {
		return defs.String() + rewritten
	}
	return rewritten[:loc[1]] + defs.String() + rewritten[loc[1]:]
}

// stripSVGWrapper returns the inner content of a `<svg ...>…</svg>` body
// plus the outer viewBox (if any). If the input isn't wrapped in <svg>,
// it's returned verbatim and viewBox is "".
func stripSVGWrapper(body string) (inner string, viewBox string) {
	open := svgOpenTagRE.FindStringIndex(body)
	if open == nil {
		return body, ""
	}
	// Pull the viewBox out of the opening tag.
	openTag := body[open[0]:open[1]]
	if vb := viewBoxAttrRE.FindStringSubmatch(openTag); vb != nil {
		viewBox = vb[1]
		if viewBox == "" {
			viewBox = vb[2]
		}
	}
	closeIdx := strings.LastIndex(body, "</svg>")
	if closeIdx < 0 || closeIdx < open[1] {
		// Mismatched/self-closing; treat the post-open-tag tail as the body.
		return body[open[1]:], viewBox
	}
	return body[open[1]:closeIdx], viewBox
}

// rewriteHref replaces the href value inside a single <use ...> tag with
// newVal. Preserves quote style and surrounding attributes. If the tag
// has both href and xlink:href, only the first encountered is rewritten;
// in practice we only emit one. If no href is found, returns the input
// unchanged.
func rewriteHref(tag, newVal string) string {
	loc := hrefAttrRE.FindStringSubmatchIndex(tag)
	if loc == nil {
		return tag
	}
	// loc layout: [match-start, match-end, g1-start, g1-end, g2-start, g2-end]
	// One of g1 / g2 is the value range; the other is -1/-1.
	var valStart, valEnd int
	if loc[2] >= 0 {
		valStart, valEnd = loc[2], loc[3]
	} else {
		valStart, valEnd = loc[4], loc[5]
	}
	return tag[:valStart] + newVal + tag[valEnd:]
}

// sanitizeSymbolID converts an art spec like "wf:phone-card" or
// `text:"hi"@mono?size=14` into a stable, HTML-valid id. Strategy:
// keep [A-Za-z0-9_-], replace ':' with '-', everything else with '_'.
func sanitizeSymbolID(spec string) string {
	var b strings.Builder
	b.Grow(len(spec))
	for _, r := range spec {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-', r == '_':
			b.WriteRune(r)
		case r == ':':
			b.WriteByte('-')
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}
