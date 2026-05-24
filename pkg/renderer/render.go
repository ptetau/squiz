package renderer

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed templates/index.html.tmpl templates/styles.css
var assets embed.FS

// {{squizId}} marker pattern. ID is [a-zA-Z][a-zA-Z0-9_-]*.
var markerRE = regexp.MustCompile(`\{\{([a-zA-Z][a-zA-Z0-9_-]*)\}\}`)

// RenderOpts captures everything the agent CAN'T put in the JSON itself:
// the absolute output path (for self-referential anchors) and overrides
// that come from the CLI rather than the document.
type RenderOpts struct {
	OutputPath    string // absolute path of the .html being written
	ThemeOverride string // from --theme; trumps doc.Theme
	WorkDir       string // dir used for repo→theme resolution
}

func Render(d *Document, opts RenderOpts) (string, error) {
	cssBytes, err := assets.ReadFile("templates/styles.css")
	if err != nil {
		return "", err
	}
	tmplBytes, err := assets.ReadFile("templates/index.html.tmpl")
	if err != nil {
		return "", err
	}

	// Theme precedence: CLI flag > doc.Theme > repo-derived auto-rotation.
	override := opts.ThemeOverride
	if override == "" {
		override = d.Theme
	}
	d.Theme = ResolveTheme(opts.WorkDir, override)

	// Lookup map for resolving {{markers}} in spec paragraphs.
	squizByID := make(map[string]Squiz, len(d.Squizzes))
	for _, s := range d.Squizzes {
		squizByID[s.ID] = s
	}

	// Pre-render spec paragraphs: replace {{markers}} with anchor chips.
	specHTML := make([]template.HTML, 0, len(d.Spec.Paragraphs))
	for _, p := range d.Spec.Paragraphs {
		specHTML = append(specHTML, renderParagraph(p.Text, squizByID))
	}

	// Build the view model. Resolves art form per option and computes a
	// default letter label if the author left Label empty.
	squizViews := make([]SquizView, 0, len(d.Squizzes))
	for _, s := range d.Squizzes {
		optViews := make([]OptionView, 0, len(s.Options))
		for i, o := range s.Options {
			art := o.ResolvedArt()
			svg, hidden := RenderArt(art, i)
			label := o.Label
			if label == "" {
				label = "Option " + LetterFor(i)
			}
			optViews = append(optViews, OptionView{
				Option:  o,
				Label:   label,
				ArtHTML: template.HTML(svg),
				HideArt: hidden,
				Index:   i,
				Letter:  LetterFor(i),
			})
		}
		squizViews = append(squizViews, SquizView{Squiz: s, Options: optViews})
	}

	// Client-side data: per-decision metadata for the JS export builder.
	// Includes the anchor (#squiz-id) so the export JSON is self-locating.
	clientSquizzes := make([]map[string]any, 0, len(d.Squizzes))
	for _, s := range d.Squizzes {
		opts := make([]map[string]any, 0, len(s.Options))
		for _, o := range s.Options {
			opts = append(opts, map[string]any{
				"id":   o.ID,
				"name": o.Name,
				"desc": o.Desc,
			})
		}
		clientSquizzes = append(clientSquizzes, map[string]any{
			"id":      s.ID,
			"title":   s.Title,
			"anchor":  "#squiz-" + s.ID,
			"options": opts,
		})
	}
	squizzesJSON, _ := json.Marshal(clientSquizzes)
	specTitleJSON, _ := json.Marshal(d.Spec.Title)

	// Source: absolute path of the rendered HTML, embedded so the export
	// JSON can carry it back to the agent.
	source := map[string]any{
		"file":     opts.OutputPath,
		"basename": filepath.Base(opts.OutputPath),
	}
	sourceJSON, _ := json.Marshal(source)

	scanlines := "off"
	if d.Scanlines {
		scanlines = "on"
	}
	cursor := "on"
	if d.Cursor != nil && !*d.Cursor {
		cursor = "off"
	}

	type tmplData struct {
		Doc            *Document
		CSS            template.CSS
		SpecParagraphs []template.HTML
		Squizzes       []SquizView
		SquizzesJSON   template.JS
		SpecTitleJSON  template.JS
		SourceJSON     template.JS
		ShowSpec       bool
		ScanlinesAttr  string
		CursorAttr     string
		Total          int
	}

	funcs := template.FuncMap{
		"safesvg": func(s template.HTML) template.HTML { return s },
		"pad2":    func(n int) string { return fmt.Sprintf("%02d", n) },
		"inc":     func(i int) int { return i + 1 },
		"plural": func(n int) string {
			if n == 1 {
				return ""
			}
			return "s"
		},
	}

	tmpl, err := template.New("index").Funcs(funcs).Parse(string(tmplBytes))
	if err != nil {
		return "", err
	}

	data := tmplData{
		Doc:            d,
		CSS:            template.CSS(cssBytes),
		SpecParagraphs: specHTML,
		Squizzes:       squizViews,
		SquizzesJSON:   template.JS(squizzesJSON),
		SpecTitleJSON:  template.JS(specTitleJSON),
		SourceJSON:     template.JS(sourceJSON),
		ShowSpec:       len(d.Spec.Paragraphs) > 0,
		ScanlinesAttr:  scanlines,
		CursorAttr:     cursor,
		Total:          len(d.Squizzes),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// SquizView / OptionView: template view models. Carry pre-computed values
// (resolved art, default labels, letter for keyboard nav) so the template
// stays presentation-only.

type SquizView struct {
	Squiz
	Options []OptionView
}

type OptionView struct {
	Option
	Label   string        // resolved (auto-derived if blank)
	ArtHTML template.HTML // resolved SVG (raw / wf: / DSL / auto)
	HideArt bool          // true → skip the art slot entirely
	Index   int           // 0-based position
	Letter  string        // "A", "B", "C", … for keyboard hints
}

// renderParagraph walks `text` and replaces each {{id}} marker with a
// <a class="squiz-mark" data-squiz="id" href="#squiz-id">first few words</a>.
// Unknown markers are emitted as literal text so authors notice typos.
func renderParagraph(text string, squizByID map[string]Squiz) template.HTML {
	var sb strings.Builder
	last := 0
	matches := markerRE.FindAllStringSubmatchIndex(text, -1)
	for _, m := range matches {
		sb.WriteString(template.HTMLEscapeString(text[last:m[0]]))
		id := text[m[2]:m[3]]
		squiz, ok := squizByID[id]
		if !ok {
			sb.WriteString(template.HTMLEscapeString(text[m[0]:m[1]]))
		} else {
			label := shortLabel(squiz.Title)
			sb.WriteString(`<a href="#squiz-`)
			sb.WriteString(template.HTMLEscapeString(id))
			sb.WriteString(`" class="squiz-mark" data-squiz="`)
			sb.WriteString(template.HTMLEscapeString(id))
			sb.WriteString(`">`)
			sb.WriteString(template.HTMLEscapeString(label))
			sb.WriteString(`</a>`)
		}
		last = m[1]
	}
	sb.WriteString(template.HTMLEscapeString(text[last:]))
	return template.HTML(sb.String())
}

// shortLabel turns a squiz title into the chip text shown in the spec.
// "How does a habit visualize over time?" → "how does a habit"
func shortLabel(title string) string {
	cleaned := strings.Map(func(r rune) rune {
		switch r {
		case '?', '.', ',', '!', ':':
			return -1
		}
		return r
	}, strings.ToLower(title))
	words := strings.Fields(cleaned)
	if len(words) > 4 {
		words = words[:4]
	}
	return strings.Join(words, " ")
}
