package planview

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ParseIndex parses raw JSON bytes into an Index (no file I/O). Useful
// in tests and for callers that get JSON from somewhere other than disk.
//
// Defaults applied:
//   - Density defaults to "compact" when empty (mirrors squiz).
//   - Theme is left empty so the renderer's auto-rotation can take over.
func ParseIndex(data []byte) (*Index, error) {
	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parse index: %w", err)
	}
	if idx.Density == "" {
		idx.Density = "compact"
	}
	return &idx, nil
}

// LoadPlan reads plan/index.json at the given path, walks its `sections`
// list, loads each <section>.json sibling, validates the result, and
// returns a fully-populated Plan ready for rendering.
//
// Validates:
//   - index.json parses
//   - every section file referenced by Index.Sections exists and parses
//   - every Item.ID starts with the section's prefix (SectionPrefix map);
//     unknown-section IDs are not prefix-checked
//   - every Item.ID is unique across the whole plan
//   - every Item.Refs entry is the ID of an item that actually exists in
//     the plan (forward refs allowed — order within the section list does
//     not constrain ref direction)
//
// Section render order: canonical sections appear in CanonicalSections
// order regardless of their position in index.Sections; custom sections
// (not in CanonicalSections) are appended in the order they were declared.
func LoadPlan(indexPath string) (*Plan, error) {
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("read index %s: %w", indexPath, err)
	}
	idx, err := ParseIndex(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", indexPath, err)
	}
	if len(idx.Sections) == 0 {
		return nil, fmt.Errorf("%s: sections list is empty", indexPath)
	}

	dir := filepath.Dir(indexPath)

	// Track declared sections, deduped, preserving first-occurrence order.
	declared := make([]string, 0, len(idx.Sections))
	seen := make(map[string]bool, len(idx.Sections))
	for _, sid := range idx.Sections {
		if seen[sid] {
			continue
		}
		seen[sid] = true
		declared = append(declared, sid)
	}

	// Load every declared section file.
	loaded := make(map[string]Section, len(declared))
	for _, sid := range declared {
		sectionPath := filepath.Join(dir, sid+".json")
		sectionRel := relForError(indexPath, sectionPath)

		raw, err := os.ReadFile(sectionPath)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", sectionRel, err)
		}
		var sf SectionFile
		if err := json.Unmarshal(raw, &sf); err != nil {
			return nil, fmt.Errorf("%s: parse: %w", sectionRel, err)
		}

		label, ok := SectionLabel[sid]
		if !ok {
			label = titleCase(sid)
		}

		// Prefix check (only for canonical sections).
		if prefix, ok := SectionPrefix[sid]; ok {
			for _, it := range sf.Items {
				if !hasPrefix(it.ID, prefix) {
					return nil, fmt.Errorf(
						"%s: item %q doesn't match section prefix %q",
						sectionRel, it.ID, prefix,
					)
				}
			}
		}

		loaded[sid] = Section{
			ID:    sid,
			Label: label,
			Items: sf.Items,
		}
	}

	// Build render order: canonical sections in canonical order, then
	// custom sections in declared order.
	ordered := make([]Section, 0, len(declared))
	canonical := make(map[string]bool, len(CanonicalSections))
	for _, cid := range CanonicalSections {
		canonical[cid] = true
		if s, ok := loaded[cid]; ok {
			ordered = append(ordered, s)
		}
	}
	for _, sid := range declared {
		if canonical[sid] {
			continue
		}
		ordered = append(ordered, loaded[sid])
	}

	// Validate uniqueness of IDs across the whole plan.
	// idOwner maps each item ID to the section file it was found in,
	// for friendly error messages on duplicates.
	idOwner := make(map[string]string, 32)
	for _, sid := range declared {
		sectionRel := relForError(indexPath, filepath.Join(dir, sid+".json"))
		for _, it := range loaded[sid].Items {
			if prev, dup := idOwner[it.ID]; dup {
				if prev == sectionRel {
					return nil, fmt.Errorf(
						"%s: item id %q used twice",
						sectionRel, it.ID,
					)
				}
				return nil, fmt.Errorf(
					"%s: item id %q already declared in %s",
					sectionRel, it.ID, prev,
				)
			}
			idOwner[it.ID] = sectionRel
		}
	}

	// Validate refs: every Item.Refs entry must point at an existing ID.
	for _, sid := range declared {
		sectionRel := relForError(indexPath, filepath.Join(dir, sid+".json"))
		for _, it := range loaded[sid].Items {
			for _, ref := range it.Refs {
				if _, ok := idOwner[ref]; !ok {
					return nil, fmt.Errorf(
						"%s: item %s references missing parent %s",
						sectionRel, it.ID, ref,
					)
				}
			}
		}
	}

	return &Plan{
		Title:    idx.Title,
		Lede:     idx.Lede,
		Theme:    idx.Theme,
		Density:  idx.Density,
		Sections: ordered,
	}, nil
}

// hasPrefix reports whether id starts with prefix followed by "-".
// Items use the form PREFIX-suffix (e.g. "FR-1", "BUILD-cli-flags").
func hasPrefix(id, prefix string) bool {
	return strings.HasPrefix(id, prefix+"-")
}

// titleCase produces a presentable label for an unknown section ID.
// Mirrors what SectionLabel would have provided.
func titleCase(s string) string {
	if s == "" {
		return ""
	}
	parts := strings.Split(s, "-")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "-")
}

// relForError tries to express sectionPath relative to the index file's
// parent directory's parent (so errors read "plan/foo.json"). Falls back
// to just the basename if anything goes wrong.
func relForError(indexPath, sectionPath string) string {
	dir := filepath.Dir(indexPath)
	parent := filepath.Dir(dir)
	if parent == "" || parent == "." || parent == dir {
		return filepath.Base(sectionPath)
	}
	rel, err := filepath.Rel(parent, sectionPath)
	if err != nil {
		return filepath.Base(sectionPath)
	}
	// Normalise to forward slashes so error messages are stable across OSes.
	return filepath.ToSlash(rel)
}
