package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/ptetau/squiz/internal/planview"
	"github.com/ptetau/squiz/pkg/renderer"
)

// validateWarning is the shape returned by runLints. Same (Path, Message)
// pair as validateError so the JSON consumer treats both arrays uniformly.
type validateWarning = validateError

// useRefRE matches <use href="<lib>:<name>" ...>. Captures library + name.
// Accepts " or ' as the quote char.
var useRefRE = regexp.MustCompile(`<use\s+[^>]*href\s*=\s*["'](wf|arch):([a-zA-Z][a-zA-Z0-9_-]*)["'][^>]*/?>`)

// useAttrRE / useTagRE — together they walk every <use .../> tag and
// recover its (x, y, width, height) attrs regardless of order.
var useAttrRE = regexp.MustCompile(`(\w+)\s*=\s*["']([^"']*)["']`)
var useTagRE = regexp.MustCompile(`<use\s[^>]*/?>`)

// svgViewBoxRE — viewBox attr anywhere in a top-level <svg> open tag.
var svgViewBoxRE = regexp.MustCompile(`viewBox\s*=\s*["'][^"']*["']`)

// runLints walks the loaded Plan and returns soft warnings nudging the
// author away from antipatterns flagged in the v0.8.0 diagram audit.
// Warnings never escalate to errors — the caller never bumps exit code.
//
// Lints implemented:
//   - unknown-use-ref       (composed <use href="wf:NAME"/> where NAME isn't in the library)
//   - oob-use-box           (<use> extends beyond standard 100x60 viewBox)
//   - missing-viewbox       (top-level raw <svg> with no viewBox)
//   - sibling-art-collision (≥2 options in one Item.Options sharing the same art)
//   - composition-thin      (per section: >60% items use single-token art)
//   - default-art-fallback  (per section: ≥2 items rely on the section default)
//   - section-art-collision (per section: ≥3 items share the same art string)
func runLints(plan *planview.Plan) []validateWarning {
	var out []validateWarning
	for _, sec := range plan.Sections {
		out = append(out, runSectionLints(sec)...)
	}
	return out
}

func runSectionLints(sec planview.Section) []validateWarning {
	var out []validateWarning
	sectionPath := sec.ID + ".json"

	// Per-item lints + per-chooser sibling collision.
	for _, it := range sec.Items {
		itemPath := fmt.Sprintf("%s:%s", sectionPath, it.ID)
		out = append(out, lintArtString(itemPath, it.Art)...)

		for _, o := range it.Options {
			optPath := fmt.Sprintf("%s.%s", itemPath, o.ID)
			out = append(out, lintArtString(optPath, o.Art)...)
		}

		if len(it.Options) >= 2 {
			out = append(out, lintSiblingCollision(itemPath, planOptionsToArtPairs(it.Options))...)
		}
	}

	// Section-wide lints.
	out = append(out, lintCompositionThin(sectionPath, sec)...)
	out = append(out, lintDefaultArtFallback(sectionPath, sec)...)
	out = append(out, lintSectionCollision(sectionPath, sec)...)

	return out
}

// artPair labels an art spec with its owning option/item id so the
// collision message can list members.
type artPair struct {
	label string
	art   string
}

func planOptionsToArtPairs(opts []planview.Option) []artPair {
	out := make([]artPair, 0, len(opts))
	for i, o := range opts {
		label := o.ID
		if label == "" {
			label = fmt.Sprintf("[%d]", i)
		}
		out = append(out, artPair{label: label, art: o.Art})
	}
	return out
}

// lintArtString runs the per-string lints (unknown-use-ref, oob-use-box,
// missing-viewbox) on one art spec.
func lintArtString(path, art string) []validateWarning {
	var out []validateWarning
	trimmed := strings.TrimSpace(art)
	if trimmed == "" {
		return nil
	}

	// missing-viewbox: top-level raw SVG without a viewBox attr in its
	// opening tag.
	if strings.HasPrefix(trimmed, "<svg") {
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

	// unknown-use-ref: every wf:/arch: <use> ref must resolve.
	for _, m := range useRefRE.FindAllStringSubmatch(trimmed, -1) {
		lib, name := m[1], m[2]
		if !libraryHas(lib, name) {
			out = append(out, validateWarning{
				Path:    path,
				Message: fmt.Sprintf("unknown-use-ref: <use href=%q> not in %s library", lib+":"+name, lib),
			})
		}
	}

	// oob-use-box: per <use ... /> tag, when all four geometry attrs
	// resolve numerically, check x+width > 100 OR y+height > 60.
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

// lintSiblingCollision flags any group of ≥2 options inside the same
// chooser sharing the same art string. Empty/"none" specs aren't
// collisions (they're absences); skip them.
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

// lintCompositionThin counts items per section and flags sections where
// >60% of items use single-token art. Sections with <3 items are skipped
// (too small to draw conclusions). Single-token = empty, "none", wf:NAME,
// arch:NAME, or any non-<svg-prefixed spec — i.e. NOT composed raw SVG.
func lintCompositionThin(scope string, sec planview.Section) []validateWarning {
	total := len(sec.Items)
	if total < 3 {
		return nil
	}
	single := 0
	for _, it := range sec.Items {
		if isSingleToken(it.Art) {
			single++
		}
	}
	ratio := float64(single) / float64(total)
	if ratio > 0.6 {
		return []validateWarning{{
			Path:    scope,
			Message: fmt.Sprintf("composition-thin: section %q: %d/%d items use single-token art — likely under-composed", sec.ID, single, total),
		}}
	}
	return nil
}

// isSingleToken classifies an art spec as "single token" (uncomposed):
// empty, "none", wf:NAME, arch:NAME, or a DSL primitive (anything that
// doesn't start with `<svg`). Composition is signalled by raw SVG, so any
// non-<svg input is treated as single-token here.
func isSingleToken(art string) bool {
	trimmed := strings.TrimSpace(art)
	if trimmed == "" || trimmed == "none" {
		return true
	}
	if strings.HasPrefix(trimmed, "<svg") {
		return false
	}
	return true
}

// lintDefaultArtFallback counts items where Art == "" (they get the
// section default at render). ≥2 such items per section is a soft hint to
// compose something more specific.
func lintDefaultArtFallback(scope string, sec planview.Section) []validateWarning {
	fallback := 0
	for _, it := range sec.Items {
		if strings.TrimSpace(it.Art) == "" {
			fallback++
		}
	}
	if fallback >= 2 {
		return []validateWarning{{
			Path:    scope,
			Message: fmt.Sprintf("default-art-fallback: section %q: %d items rely on the section default — consider composing", sec.ID, fallback),
		}}
	}
	return nil
}

// lintSectionCollision flags art specs shared by ≥3 items in one section.
// Catches "every engineering item is arch:server". Empty + "none" are
// absences and are skipped.
func lintSectionCollision(scope string, sec planview.Section) []validateWarning {
	groups := map[string][]string{}
	for i, it := range sec.Items {
		art := strings.TrimSpace(it.Art)
		if art == "" || art == "none" {
			continue
		}
		label := it.ID
		if label == "" {
			label = fmt.Sprintf("[%d]", i)
		}
		groups[art] = append(groups[art], label)
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var out []validateWarning
	for _, art := range keys {
		members := groups[art]
		if len(members) < 3 {
			continue
		}
		out = append(out, validateWarning{
			Path:    scope,
			Message: fmt.Sprintf("section-art-collision: section %q: items [%s] all use art %q", sec.ID, strings.Join(members, ", "), art),
		})
	}
	return out
}
