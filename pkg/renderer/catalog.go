package renderer

import (
	"encoding/json"
	"fmt"
	"html"
	"sort"
	"strings"
)

// CatalogEntry is the JSON shape emitted by `squiz catalog <name> --json`.
type CatalogEntry struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Category string `json:"category,omitempty"`
}

// ThemeInfo describes one of the 8 ship themes for `catalog themes`.
type ThemeInfo struct {
	Name string `json:"name"`
	Vibe string `json:"vibe"`
	Mode string `json:"mode"` // "light" or "dark"
}

// ThemeCatalog is the canonical lookup of theme name → vibe text + mode.
// Vibe text lifted from skills/squiz/SKILL.md's themes table.
var ThemeCatalog = []ThemeInfo{
	{Name: "paper", Vibe: "Cream + ink + rust accent. Editorial, calm.", Mode: "light"},
	{Name: "phosphor", Vibe: "Green-on-black CRT.", Mode: "dark"},
	{Name: "amber", Vibe: "IBM 3279 amber on near-black.", Mode: "dark"},
	{Name: "beige", Vibe: "PS/2 cream with IBM blue.", Mode: "light"},
	{Name: "rose", Vibe: "Warm pink, plum ink, rose accent.", Mode: "light"},
	{Name: "ocean", Vibe: "Pale blue-grey, deep teal, coral accent.", Mode: "light"},
	{Name: "forest", Vibe: "Oat cream, moss, warm gold.", Mode: "light"},
	{Name: "slate", Vibe: "Cool dark grey, electric blue accent.", Mode: "dark"},
}

// DSLPrimitive describes one DSL primitive for `catalog dsl`.
type DSLPrimitive struct {
	Grammar string `json:"grammar"`
	Desc    string `json:"desc"`
}

// DSLCatalog is the canonical list of the 11 DSL primitives that
// resolveDSL in dsl.go dispatches on.
var DSLCatalog = []DSLPrimitive{
	{Grammar: `grid:NxM[@RATE]`, Desc: "N×M heatmap, RATE in [0,1]"},
	{Grammar: `spark:[v,v,...]`, Desc: "sparkline from data"},
	{Grammar: `bars:[v,v,...]`, Desc: "bar chart"},
	{Grammar: `swatches:#A,#B,...`, Desc: "palette swatches"},
	{Grammar: `pills:A*|B|C*`, Desc: "chip row, * = active"},
	{Grammar: `sample:"text"[@FONT]`, Desc: "styled sample text, FONT = serif/sans/mono"},
	{Grammar: `circle-pack:N`, Desc: "N organically-arranged circles"},
	{Grammar: `text:"..."[@FONT]?attrs`, Desc: "multi-line styled text (size, align, weight, color)"},
	{Grammar: `flow:[a,b,c]`, Desc: "L-to-R pipeline of named boxes (?icon=arch-name)"},
	{Grammar: `box:label[?icon=NAME]`, Desc: "single labeled box"},
	{Grammar: `arrow:"label"[?dir=DIR]`, Desc: "labeled arrow glyph"},
}

// CatalogNames returns the list of supported catalog names for the top-level
// `squiz catalog` no-arg listing.
var CatalogNames = []string{"wf", "arch", "dsl", "themes"}

// WFCatalog returns the WFLibrary as a sorted slice of CatalogEntry. Sort
// is by name so output is deterministic across runs.
func WFCatalog() []CatalogEntry {
	out := make([]CatalogEntry, 0, len(WFLibrary))
	for name := range WFLibrary {
		out = append(out, CatalogEntry{
			Name:     name,
			Desc:     WFDescriptions[name],
			Category: wfCategory[name],
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// ArchCatalog returns the ArchLibrary as a sorted slice of CatalogEntry.
func ArchCatalog() []CatalogEntry {
	out := make([]CatalogEntry, 0, len(ArchLibrary))
	for name := range ArchLibrary {
		out = append(out, CatalogEntry{
			Name:     name,
			Desc:     ArchDescriptions[name],
			Category: archCategory[name],
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// FormatCatalogText renders entries as "name   desc" lines with the names
// padded to a common width so descriptions visually line up. Stable
// regardless of the longest name in the registry.
func FormatCatalogText(entries []CatalogEntry) string {
	maxName := 0
	for _, e := range entries {
		if len(e.Name) > maxName {
			maxName = len(e.Name)
		}
	}
	pad := maxName + 2
	var b strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&b, "%-*s%s\n", pad, e.Name, e.Desc)
	}
	return b.String()
}

// FormatCatalogJSON marshals entries as a JSON array. Pretty-printed with
// 2-space indent so it's diffable and copyable.
func FormatCatalogJSON(entries []CatalogEntry) (string, error) {
	b, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

// FormatDSLText renders DSLCatalog as aligned grammar + description rows.
func FormatDSLText() string {
	maxG := 0
	for _, p := range DSLCatalog {
		if len(p.Grammar) > maxG {
			maxG = len(p.Grammar)
		}
	}
	pad := maxG + 2
	var b strings.Builder
	for _, p := range DSLCatalog {
		fmt.Fprintf(&b, "%-*s%s\n", pad, p.Grammar, p.Desc)
	}
	return b.String()
}

// FormatDSLJSON marshals the DSL primitives as a JSON array.
func FormatDSLJSON() (string, error) {
	b, err := json.MarshalIndent(DSLCatalog, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

// FormatThemesText renders ThemeCatalog as `name vibe (mode)` rows, aligned.
func FormatThemesText() string {
	maxN, maxV := 0, 0
	for _, t := range ThemeCatalog {
		if len(t.Name) > maxN {
			maxN = len(t.Name)
		}
		if len(t.Vibe) > maxV {
			maxV = len(t.Vibe)
		}
	}
	var b strings.Builder
	for _, t := range ThemeCatalog {
		fmt.Fprintf(&b, "%-*s  %-*s  (%s)\n", maxN, t.Name, maxV, t.Vibe, t.Mode)
	}
	return b.String()
}

// FormatThemesJSON marshals ThemeCatalog as a JSON array.
func FormatThemesJSON() (string, error) {
	b, err := json.MarshalIndent(ThemeCatalog, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b) + "\n", nil
}

// FormatNamesText prints the top-level catalog names (one per line) for
// `squiz catalog` with no argument.
func FormatNamesText() string {
	var b strings.Builder
	for _, n := range CatalogNames {
		b.WriteString(n)
		b.WriteByte('\n')
	}
	return b.String()
}

// RenderGalleryHTML builds a self-contained HTML gallery page showing each
// entry's resolved SVG, name, and description. The renderFn parameter lets
// the caller plug in either resolveNamed (for wf:*) or resolveArch (for
// arch:*) without exposing those internals.
//
// `theme` selects which of the 8 ship themes the page uses (default
// "paper"). The full themes bundle is embedded so the page works offline
// and the user can flip `data-theme` if they want to A/B.
//
// `title` is shown in the page <h1>.
func RenderGalleryHTML(title, theme string, entries []CatalogEntry, renderFn func(name string) string) string {
	if !validTheme(theme) {
		theme = "paper"
	}
	var b strings.Builder
	b.WriteString("<!doctype html>\n<html lang='en' data-theme='")
	b.WriteString(html.EscapeString(theme))
	b.WriteString("'>\n<head>\n<meta charset='utf-8'>\n<title>")
	b.WriteString(html.EscapeString(title))
	b.WriteString("</title>\n<style>\n")
	b.WriteString(ThemeCSS())
	b.WriteString(`
body { background: var(--bg); color: var(--ink); font-family: 'IBM Plex Sans', system-ui, sans-serif; margin: 0; padding: 24px; }
h1 { font-family: 'IBM Plex Mono', monospace; font-size: 18px; letter-spacing: 0.04em; margin: 0 0 4px; color: var(--ink); text-transform: uppercase; }
.sub { font-family: 'IBM Plex Mono', monospace; font-size: 11px; color: var(--ink-3); letter-spacing: 0.08em; text-transform: uppercase; margin: 0 0 24px; }
.grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 16px; }
.card { background: var(--surface); border: 1px solid var(--rule); padding: 12px; display: flex; flex-direction: column; gap: 8px; }
.art { aspect-ratio: 5 / 3; background: var(--bg-2); display: flex; align-items: center; justify-content: center; overflow: hidden; }
.art svg { max-width: 95%; max-height: 95%; }
.name { font-family: 'IBM Plex Mono', monospace; font-size: 12px; font-weight: 700; color: var(--accent); }
.desc { font-family: 'IBM Plex Sans', sans-serif; font-size: 12px; color: var(--ink-2); line-height: 1.4; }
.cat { font-family: 'IBM Plex Mono', monospace; font-size: 9px; color: var(--ink-3); letter-spacing: 0.12em; text-transform: uppercase; }
</style>
</head>
<body>
`)
	fmt.Fprintf(&b, "<h1>%s</h1>\n<p class='sub'>%d entries · theme: %s</p>\n",
		html.EscapeString(title), len(entries), html.EscapeString(theme))
	b.WriteString("<div class='grid'>\n")
	for _, e := range entries {
		svg := renderFn(e.Name)
		b.WriteString("  <div class='card'>\n")
		b.WriteString("    <div class='art'>")
		b.WriteString(svg) // SVG is trusted (we generated it).
		b.WriteString("</div>\n")
		fmt.Fprintf(&b, "    <div class='name'>%s</div>\n", html.EscapeString(e.Name))
		if e.Category != "" {
			fmt.Fprintf(&b, "    <div class='cat'>%s</div>\n", html.EscapeString(e.Category))
		}
		fmt.Fprintf(&b, "    <div class='desc'>%s</div>\n", html.EscapeString(e.Desc))
		b.WriteString("  </div>\n")
	}
	b.WriteString("</div>\n</body>\n</html>\n")
	return b.String()
}

// WFRender returns the resolved SVG string for a wf:* name (or a labeled
// placeholder if the name is unknown). Stable, hidden=false ignored.
func WFRender(name string) string {
	svg, _ := resolveNamed(name)
	return svg
}

// ArchRender returns the resolved SVG string for an arch:* name.
func ArchRender(name string) string {
	svg, _ := resolveArch(name)
	return svg
}

// RenderPreviewHTML builds a self-contained page showing one art form at
// large size with the spec string printed below. `spec` is the raw user
// input (e.g. `wf:calendar-grid` or `flow:[client,api,db]`) and is escaped
// before embedding. `svg` is the resolved SVG fragment (trusted output of
// RenderArt).
func RenderPreviewHTML(spec, svg, theme string) string {
	if !validTheme(theme) {
		theme = "paper"
	}
	var b strings.Builder
	b.WriteString("<!doctype html>\n<html lang='en' data-theme='")
	b.WriteString(html.EscapeString(theme))
	b.WriteString("'>\n<head>\n<meta charset='utf-8'>\n<title>preview · ")
	b.WriteString(html.EscapeString(spec))
	b.WriteString("</title>\n<style>\n")
	b.WriteString(ThemeCSS())
	b.WriteString(`
body { background: var(--bg); color: var(--ink); font-family: 'IBM Plex Sans', system-ui, sans-serif; margin: 0; padding: 32px; display: flex; flex-direction: column; align-items: center; gap: 24px; min-height: 100vh; box-sizing: border-box; }
h1 { font-family: 'IBM Plex Mono', monospace; font-size: 14px; letter-spacing: 0.08em; color: var(--ink-3); text-transform: uppercase; margin: 0; }
.stage { background: var(--surface); border: 1px solid var(--rule); padding: 32px; width: 80%; max-width: 720px; aspect-ratio: 5 / 3; display: flex; align-items: center; justify-content: center; }
.stage svg { width: 100%; height: auto; max-height: 100%; }
.spec { font-family: 'IBM Plex Mono', monospace; font-size: 13px; color: var(--accent); background: var(--bg-2); border: 1px solid var(--rule-2); padding: 8px 14px; }
.foot { font-family: 'IBM Plex Mono', monospace; font-size: 10px; color: var(--ink-3); letter-spacing: 0.1em; text-transform: uppercase; }
</style>
</head>
<body>
`)
	b.WriteString("<h1>squiz preview</h1>\n")
	b.WriteString("<div class='stage'>")
	b.WriteString(svg) // trusted
	b.WriteString("</div>\n")
	fmt.Fprintf(&b, "<div class='spec'>%s</div>\n", html.EscapeString(spec))
	fmt.Fprintf(&b, "<div class='foot'>theme · %s</div>\n", html.EscapeString(theme))
	b.WriteString("</body>\n</html>\n")
	return b.String()
}
