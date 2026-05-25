package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ptetau/squiz/pkg/renderer"
)

// validateWarning is the shape returned by runLints. It reuses the same
// (Path, Message) pair as validateError so the JSON consumer treats both
// arrays uniformly.
type validateWarning = validateError

// useRefRE matches <use href="<lib>:<name>" ...>. We capture the library
// prefix (wf|arch) and the name separately so the lint can dispatch into
// the right registry.
//
// We accept both " and ' as the quote (the renderer's built-in entries use
// '). The trailing .*? before /> tolerates the x/y/width/height attrs
// surrounding the href.
var useRefRE = regexp.MustCompile(`<use\s+[^>]*href\s*=\s*["'](wf|arch):([a-zA-Z][a-zA-Z0-9_-]*)["'][^>]*/?>`)

// useAttrRE picks every attribute (name="value" or name='value') out of a
// single `<use ... />` tag so we can pull x/y/width/height regardless of
// order. The lint composes a per-tag scan from this.
var useAttrRE = regexp.MustCompile(`(\w+)\s*=\s*["']([^"']*)["']`)

// useTagRE matches one full <use .../> tag — we feed each match through
// useAttrRE to recover its attributes.
var useTagRE = regexp.MustCompile(`<use\s[^>]*/?>`)

// svgViewBoxRE detects a viewBox attribute anywhere in the top-level <svg>
// open tag. Used by missing-viewbox.
var svgViewBoxRE = regexp.MustCompile(`viewBox\s*=\s*["'][^"']*["']`)

// runLints walks the parsed Document and returns soft warnings that nudge
// authors away from the antipatterns identified in the v0.8.0 diagram
// audit. None of these escalate to errors — the caller never bumps the
// exit code on them.
//
// Lints implemented:
//   - unknown-use-ref       (composed <use href="wf:NAME"/> where NAME isn't in the library)
//   - oob-use-box           (<use x=N y=N width=N height=N/> extending beyond the 100x60 viewBox)
//   - missing-viewbox       (top-level raw <svg> with no viewBox attribute)
//   - sibling-art-collision (≥2 options in one squiz sharing the same `art` string)
func runLints(doc *renderer.Document) []validateWarning {
	var out []validateWarning
	for i, s := range doc.Squizzes {
		base := fmt.Sprintf("squizzes[%d]", i)
		if s.ID != "" {
			base = s.ID
		}

		for j, o := range s.Options {
			optBase := fmt.Sprintf("%s.options[%d]", base, j)
			if o.ID != "" {
				optBase = fmt.Sprintf("%s.%s", base, o.ID)
			}
			out = append(out, lintArtString(optBase, o.ResolvedArt())...)
		}

		out = append(out, lintSiblingCollision(base, optionsToArtPairs(s.Options))...)
	}
	return out
}

// artPair is a (label, art) pair used by sibling/section collision lints.
// Label is the option id (or item id) to surface in the warning message.
type artPair struct {
	label string
	art   string
}

func optionsToArtPairs(opts []renderer.Option) []artPair {
	out := make([]artPair, 0, len(opts))
	for i, o := range opts {
		label := o.ID
		if label == "" {
			label = fmt.Sprintf("[%d]", i)
		}
		out = append(out, artPair{label: label, art: o.ResolvedArt()})
	}
	return out
}

// lintArtString runs the per-string lints (unknown-use-ref, oob-use-box,
// missing-viewbox) on one art spec. path is the JSON-pointer-ish location
// used in the warning's Path field.
func lintArtString(path, art string) []validateWarning {
	var out []validateWarning
	trimmed := strings.TrimSpace(art)
	if trimmed == "" {
		return nil
	}

	// missing-viewbox: top-level raw SVG without a viewBox.
	if strings.HasPrefix(trimmed, "<svg") {
		// Inspect only the opening tag.
		end := strings.Index(trimmed, ">")
		openTag := trimmed
		if end > 0 {
			openTag = trimmed[:end+1]
		}
		if !svgViewBoxRE.MatchString(openTag) {
			out = append(out, validateWarning{
				Path:    path,
				Message: "missing-viewbox: top-level <svg> has no viewBox attribute (rendered size will be unpredictable)",
			})
		}
	}

	// unknown-use-ref: every <use href="wf:NAME"/> or arch:NAME → check library.
	for _, m := range useRefRE.FindAllStringSubmatch(trimmed, -1) {
		lib, name := m[1], m[2]
		if !libraryHas(lib, name) {
			out = append(out, validateWarning{
				Path:    path,
				Message: fmt.Sprintf("unknown-use-ref: <use href=%q> not in %s library", lib+":"+name, lib),
			})
		}
	}

	// oob-use-box: per <use ... /> tag, when all four geometry attrs are
	// present and numeric, check x+width > 100 OR y+height > 60.
	for _, tag := range useTagRE.FindAllString(trimmed, -1) {
		attrs := map[string]string{}
		for _, am := range useAttrRE.FindAllStringSubmatch(tag, -1) {
			attrs[am[1]] = am[2]
		}
		x, xok := parseFloatAttr(attrs["x"])
		y, yok := parseFloatAttr(attrs["y"])
		w, wok := parseFloatAttr(attrs["width"])
		h, hok := parseFloatAttr(attrs["height"])
		if !(xok && yok && wok && hok) {
			continue
		}
		if x+w > 100 || y+h > 60 {
			out = append(out, validateWarning{
				Path:    path,
				Message: fmt.Sprintf("oob-use-box: <use> at (%g,%g) %gx%g extends beyond 100x60 viewBox", x, y, w, h),
			})
		}
	}

	return out
}

func parseFloatAttr(s string) (float64, bool) {
	if s == "" {
		return 0, false
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

// libraryHas tells whether (lib, name) resolves in the canonical
// registries. Unknown lib prefixes return true so we don't false-positive
// on future libraries the renderer might add.
func libraryHas(lib, name string) bool {
	switch lib {
	case "wf":
		_, ok := renderer.WFLibrary[name]
		return ok
	case "arch":
		_, ok := renderer.ArchLibrary[name]
		return ok
	}
	return true
}

// lintSiblingCollision groups art-pairs by their normalized art spec and
// flags any group with ≥2 members. "none" and the empty string are treated
// as absences, not collisions, and skipped.
func lintSiblingCollision(scope string, pairs []artPair) []validateWarning {
	groups := map[string][]string{}
	for _, p := range pairs {
		art := strings.TrimSpace(p.art)
		if art == "" || art == "none" {
			continue
		}
		groups[art] = append(groups[art], p.label)
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var out []validateWarning
	for _, art := range keys {
		members := groups[art]
		if len(members) < 2 {
			continue
		}
		out = append(out, validateWarning{
			Path:    scope,
			Message: fmt.Sprintf("sibling-art-collision: options [%s] share art %q", strings.Join(members, ", "), art),
		})
	}
	return out
}
