// Package planview is the squiz-plan rendering pipeline.
//
// Concept: take a structured plan (overview → functional → non-functional →
// cases → engineering → build) split across multiple JSON files, validate
// the cross-references between layers, and render a single tabbed HTML
// document where every item can link back to the parents that motivated
// it. Uses the shared theme + art system from pkg/renderer.
package planview

// CanonicalSections is the fixed order the renderer enforces for the six
// built-in section types. Custom sections (declared by the agent in
// Index.Sections that are NOT in this list) are appended after these.
var CanonicalSections = []string{
	"overview",
	"functional",
	"non-functional",
	"cases",
	"engineering",
	"build",
}

// SectionDefaultArt is the fallback `art` spec used when an Item in a
// canonical section omits its own art. The intent (per the v0.7.1
// "examples or prototypes from the overview to the build" guidance) is
// that EVERY item in a canonical section gets a meaningful visual anchor
// by default — authors override per-item, or use "none" to suppress.
//
// Custom sections have no default (omitted == no art, same as today).
// Item.Art == "none" still suppresses (explicit-hide always wins).
var SectionDefaultArt = map[string]string{
	"overview":       "wf:avatar-single",  // who the plan is for / what it's about
	"functional":     "wf:phone-blank",    // the thing the system does
	"non-functional": "wf:gauge",          // how it behaves
	"cases":          "wf:phone-card",     // a scenario, narrative-shaped
	"engineering":    "arch:server",       // the build blocks
	"build":          "wf:cmd-palette",    // the steps
}

// SectionLabel turns a section ID into the display label shown in the tab
// strip. Unknown IDs are title-cased verbatim so custom sections still
// look presentable.
var SectionLabel = map[string]string{
	"overview":       "Overview",
	"functional":     "Functional",
	"non-functional": "Non-functional",
	"cases":          "Cases",
	"engineering":    "Engineering",
	"build":          "Build",
}

// SectionPrefix is the ID prefix items in each section MUST use. Used by
// the validator to catch misfiled items (an FR-* item under engineering.json
// is almost certainly a mistake).
var SectionPrefix = map[string]string{
	"overview":       "OVR",
	"functional":     "FR",
	"non-functional": "NFR",
	"cases":          "CASE",
	"engineering":    "ENG",
	"build":          "BUILD",
}

// Index is the top-level descriptor at plan/index.json. It declares
// metadata + which sections this plan includes. Section data lives in
// sibling files: plan/overview.json, plan/functional.json, etc.
type Index struct {
	Title    string   `json:"title"`              // plan H1
	Lede     string   `json:"lede"`               // one-line summary above the tabs
	Theme    string   `json:"theme,omitempty"`    // optional; same precedence as squiz
	Density  string   `json:"density,omitempty"`  // compact | comfortable
	Sections []string `json:"sections"`           // section IDs to load + render order;
	//                                                 canonical IDs sort to declared
	//                                                 position; unknown IDs append.
}

// SectionFile is the per-section payload at plan/<sectionId>.json.
type SectionFile struct {
	Items []Item `json:"items"`
}

// Item is one entry in a section. ID must use the section's prefix
// (validated). Refs are IDs of parent items the validator confirms exist.
//
// When Options is non-empty the item is a *decision* — the renderer shows
// a chooser (same shape as squiz options) and the user's pick is captured
// in the feedback export as `chose: "optionId"`. When Options is empty
// (the v0.3.0 default) the item is a *statement* — flat card, no chooser.
type Item struct {
	ID      string   `json:"id"`                // OVR-1 / FR-3.2 / BUILD-cli-flags
	Title   string   `json:"title"`             // displayed as the card heading
	Desc    string   `json:"desc"`              // 1-3 sentences
	Art     string   `json:"art,omitempty"`     // same forms as squiz (wf:/DSL/raw SVG/"none"/omitted)
	Refs    []string `json:"refs,omitempty"`    // IDs of parent items (validator checks they exist)
	Options []Option `json:"options,omitempty"` // v0.4.0: optional in-item chooser
}

// Option is one branch of an in-item decision. Same shape as squiz's
// Option (id/label/name/desc/art) so authors can copy patterns between
// the two tools. IDs are local to the item — collisions across items are
// fine; collisions within one item are not (validated).
type Option struct {
	ID    string `json:"id"`              // stable slug, comes back in feedback as `chose`
	Label string `json:"label,omitempty"` // OPTIONAL — auto-derived from index ("Option A", "B"…)
	Name  string `json:"name"`            // short display
	Desc  string `json:"desc"`            // 1-2 sentence trade-off
	Art   string `json:"art,omitempty"`   // same forms as Item.Art

	// Recommendation, when non-empty, marks this option as the author's
	// recommended choice and carries the explanation. The renderer shows
	// a "★ RECOMMENDED" chip + the explanation as a small editorial
	// callout under the option's desc. At most one option per item should
	// carry one; multiple renders all of them but is usually an authoring
	// mistake.
	Recommendation string `json:"recommendation,omitempty"`
}

// Section is the loaded form of a SectionFile, augmented with the
// section's identity. Built by the parser.
type Section struct {
	ID    string // "functional"
	Label string // "Functional"
	Items []Item
}

// Plan is the fully-loaded, validated form passed to the renderer.
type Plan struct {
	Title    string
	Lede     string
	Theme    string
	Density  string
	Sections []Section
}

// LookupItem finds an item by its full ID across all sections. Returns
// the item + its section ID, or (Item{}, "", false) if not found. Used by
// the cross-ref validator and by the renderer to expand badge hover text.
func (p *Plan) LookupItem(id string) (Item, string, bool) {
	for _, s := range p.Sections {
		for _, it := range s.Items {
			if it.ID == id {
				return it, s.ID, true
			}
		}
	}
	return Item{}, "", false
}
