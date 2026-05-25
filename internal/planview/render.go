package planview

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/ptetau/squiz/pkg/renderer"
)

// Template assets for plan rendering. The sibling "template agent"
// authors plan.html.tmpl, plan.css, and plan.js in templates/. Render()
// still surfaces a clear error if a read fails (e.g. partially-landed
// sibling work) so tests can short-circuit cleanly via templatesReady().
//
//go:embed templates/plan.html.tmpl templates/plan.css templates/plan.js
var assets embed.FS

// RenderOpts mirrors renderer.RenderOpts but is plan-specific. Captures
// everything the plan JSON itself can't carry: the absolute output path
// (used for self-referential anchors) and the CLI/theme overrides.
type RenderOpts struct {
	OutputPath    string // absolute path of the .html being written
	ThemeOverride string // from --theme; trumps plan.Theme
	WorkDir       string // dir used for repo→theme resolution
}

// PlanView is the top-level template view-model. All HTML-safe content
// uses html/template's typed strings so the template can `{{.Field}}`
// without re-escaping. *Attr fields mirror the squiz template's data-*
// attribute pattern so the same theme/density/scanlines/cursor knobs
// work uniformly across both binaries.
type PlanView struct {
	Title         string
	Lede          string
	Theme         string
	Density       string
	ThemeAttr     string
	DensityAttr   string
	ScanlinesAttr string
	CursorAttr    string
	CSS           template.CSS
	Sections      []SectionView
	PlanJSON      template.JS // for the client-side JS (item index + nav)
	SourceJSON    template.JS // {file, basename} — embedded in export payload
}

// SectionView is one tab in the rendered doc. Index is 0-based and
// useful to the template for tab ordering / aria-controls plumbing.
type SectionView struct {
	ID    string
	Label string
	Index int
	Items []ItemView
}

// ItemView is one card within a section. ArtHTML / HideArt mirror the
// squiz OptionView contract so the template agent can reuse the same
// patterns. Refs carries the resolved cross-references.
type ItemView struct {
	ID      string
	Title   string
	Desc    string
	ArtHTML template.HTML
	HideArt bool
	Refs    []RefView
}

// RefView is one resolved cross-reference. Label is the display string
// formatted as "<SectionLabel> · <itemID>" (e.g. "Functional · FR-1"),
// or "<itemID> (missing)" when the target can't be found. Missing should
// never trip in practice — LoadPlan validates refs — but the renderer
// keeps the defensive path so it can't crash on malformed input.
//
// TargetURL is a same-doc fragment ("#item-FR-1") the template uses
// for the <a href> on the rendered ref chip; ID is the bare item ID
// shown to the user.
type RefView struct {
	ID        string // "FR-1" — bare item ID (display + data-target)
	SectionID string // "functional"
	Label     string // "Functional · FR-1" — aria-label / tooltip text
	TargetURL string // "#item-FR-1"
	Missing   bool   // true → target was not in the plan
}

// Render takes a loaded *Plan and emits the final HTML string.
//
// Precedence for theme: opts.ThemeOverride > plan.Theme > auto-rotation
// from renderer.ResolveTheme. Art rendering uses letterIdx=0 (the
// per-letter rotation is squiz-only — every plan item shares the same
// "first option" pattern). CSS is the shared theme bundle concatenated
// with planview's own plan.css.
func Render(p *Plan, opts RenderOpts) (string, error) {
	if p == nil {
		return "", errors.New("Render: nil plan")
	}

	// Theme precedence: CLI flag > plan.Theme > repo-derived auto-rotation.
	override := opts.ThemeOverride
	if override == "" {
		override = p.Theme
	}
	theme := renderer.ResolveTheme(opts.WorkDir, override)

	// Build a lookup so each item's Refs can be enriched with the section
	// label of the target.
	type targetInfo struct {
		sectionID    string
		sectionLabel string
	}
	idToTarget := make(map[string]targetInfo, 32)
	for _, s := range p.Sections {
		for _, it := range s.Items {
			idToTarget[it.ID] = targetInfo{sectionID: s.ID, sectionLabel: s.Label}
		}
	}

	sectionViews := make([]SectionView, 0, len(p.Sections))
	for i, s := range p.Sections {
		itemViews := make([]ItemView, 0, len(s.Items))
		for _, it := range s.Items {
			svg, hidden := renderer.RenderArt(it.Art, 0)
			refs := make([]RefView, 0, len(it.Refs))
			for _, refID := range it.Refs {
				target, ok := idToTarget[refID]
				if !ok {
					refs = append(refs, RefView{
						ID:        refID,
						Label:     fmt.Sprintf("%s (missing)", refID),
						TargetURL: "#item-" + refID,
						Missing:   true,
					})
					continue
				}
				refs = append(refs, RefView{
					ID:        refID,
					SectionID: target.sectionID,
					// "Functional · FR-1" — middle dot (U+00B7) matches the
					// brief and the squiz topbar separator.
					Label:     fmt.Sprintf("%s · %s", target.sectionLabel, refID),
					TargetURL: "#item-" + refID,
				})
			}
			itemViews = append(itemViews, ItemView{
				ID:      it.ID,
				Title:   it.Title,
				Desc:    it.Desc,
				ArtHTML: template.HTML(svg),
				HideArt: hidden,
				Refs:    refs,
			})
		}
		sectionViews = append(sectionViews, SectionView{
			ID:    s.ID,
			Label: s.Label,
			Index: i,
			Items: itemViews,
		})
	}

	// Combined CSS: shared theme bundle + planview's own plan.css.
	planCSS, err := assets.ReadFile("templates/plan.css")
	if err != nil {
		// Templates not yet on disk (sibling agent in progress) — surface a
		// clear error rather than executing a half-baked render.
		return "", fmt.Errorf("planview: template assets not embedded yet: %w", err)
	}
	tmplBytes, err := assets.ReadFile("templates/plan.html.tmpl")
	if err != nil {
		return "", fmt.Errorf("planview: template assets not embedded yet: %w", err)
	}
	combinedCSS := renderer.ThemeCSS() + "\n" + string(planCSS)

	// Client-side payload: minimal index the JS uses for nav + export.
	// Each section carries its items; each item carries the fields the
	// JS reads (id, title, desc) so the export can attribute edits back
	// to the original content.
	clientSections := make([]map[string]any, 0, len(sectionViews))
	for _, s := range sectionViews {
		items := make([]map[string]any, 0, len(s.Items))
		for _, it := range s.Items {
			items = append(items, map[string]any{
				"id":    it.ID,
				"title": it.Title,
				"desc":  it.Desc,
			})
		}
		clientSections = append(clientSections, map[string]any{
			"id":    s.ID,
			"label": s.Label,
			"items": items,
		})
	}
	planClient := map[string]any{
		"title":    p.Title,
		"lede":     p.Lede,
		"sections": clientSections,
	}
	planJSON, _ := json.Marshal(planClient)
	sourceJSON, _ := json.Marshal(map[string]any{
		"file":     opts.OutputPath,
		"basename": filepath.Base(opts.OutputPath),
	})

	density := p.Density
	if density == "" {
		density = "compact"
	}

	view := PlanView{
		Title:    p.Title,
		Lede:     p.Lede,
		Theme:    theme,
		Density:  density,
		// data-* attrs the template wires onto <html>. Scanlines/cursor
		// don't currently surface in the plan JSON schema; defaults
		// mirror squiz's "off/on" semantics so the same CSS rules apply.
		ThemeAttr:     theme,
		DensityAttr:   density,
		ScanlinesAttr: "off",
		CursorAttr:    "on",
		CSS:           template.CSS(combinedCSS),
		Sections:      sectionViews,
		PlanJSON:      template.JS(planJSON),
		SourceJSON:    template.JS(sourceJSON),
	}

	// Mutate the input plan's Theme to the resolved value so callers
	// (the CLI prints "wrote <path> · <theme>") see what was chosen.
	p.Theme = theme

	funcs := template.FuncMap{
		"pad2":   func(n int) string { return fmt.Sprintf("%02d", n) },
		"inc":    func(i int) int { return i + 1 },
		"safesvg": func(s template.HTML) template.HTML { return s },
		"plural": func(n int) string {
			if n == 1 {
				return ""
			}
			return "s"
		},
	}

	tmpl, err := template.New("plan").Funcs(funcs).Parse(string(tmplBytes))
	if err != nil {
		return "", fmt.Errorf("planview: parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, view); err != nil {
		return "", fmt.Errorf("planview: execute template: %w", err)
	}
	return buf.String(), nil
}
