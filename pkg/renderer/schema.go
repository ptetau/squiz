package renderer

import "encoding/json"

// Document is the JSON contract between the agent and the renderer.
// The agent writes one of these to a file, then runs `squiz <file>`.
type Document struct {
	Theme     string  `json:"theme"`     // paper | phosphor | amber | beige
	Density   string  `json:"density"`   // compact | comfortable
	Scanlines bool    `json:"scanlines"` // CRT scanline overlay
	Cursor    *bool   `json:"cursor"`    // blinking cursor (default true)
	Spec      Spec    `json:"spec"`
	Squizzes  []Squiz `json:"squizzes"`
}

type Spec struct {
	Path       string      `json:"path"`       // shown in the topbar
	Title      string      `json:"title"`      // page H1
	Lede       string      `json:"lede"`       // one-line summary
	Paragraphs []Paragraph `json:"paragraphs"` // optional spec narrative
}

type Paragraph struct {
	Text string `json:"text"` // may contain {{squizId}} markers
}

type Squiz struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Desc    string   `json:"desc"`
	Quote   string   `json:"quote"` // optional pull-quote from the spec
	Options []Option `json:"options"`
}

type Option struct {
	ID    string `json:"id"`
	Label string `json:"label"` // e.g. "Option A" — auto-derived from index if empty
	Name  string `json:"name"`  // short display name
	Desc  string `json:"desc"`  // 1-2 sentence trade-off

	// Art is the unified visual-preview field. Five forms supported:
	//   - raw SVG markup:   "<svg viewBox='...'>...</svg>"
	//   - named library:    "wf:calendar-grid"
	//   - DSL primitive:    "grid:7x7@0.55", "spark:[3,5,4,7]", etc.
	//   - explicit hide:    "none"  (no art slot at all)
	//   - omitted/empty:    auto per-letter abstract pattern
	Art string `json:"art"`

	// ArtSVG is the legacy field — treated as raw SVG only. Use Art for new code.
	ArtSVG string `json:"art_svg,omitempty"`
}

// ResolvedArt returns Art when set, otherwise ArtSVG (treated as raw SVG).
func (o Option) ResolvedArt() string {
	if o.Art != "" {
		return o.Art
	}
	return o.ArtSVG
}

func ParseDocument(data []byte) (*Document, error) {
	var d Document
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, err
	}
	// Note: we intentionally don't default Theme here — leave it empty so
	// the renderer's auto-rotation can take over. Set explicitly in JSON
	// (or via --theme) to opt out.
	if d.Density == "" {
		d.Density = "compact"
	}
	if d.Cursor == nil {
		t := true
		d.Cursor = &t
	}
	return &d, nil
}
