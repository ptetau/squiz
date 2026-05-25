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
type Item struct {
	ID    string   `json:"id"`             // OVR-1 / FR-3.2 / BUILD-cli-flags
	Title string   `json:"title"`          // displayed as the card heading
	Desc  string   `json:"desc"`           // 1-3 sentences
	Art   string   `json:"art,omitempty"`  // same forms as squiz (wf:/DSL/raw SVG/"none"/omitted)
	Refs  []string `json:"refs,omitempty"` // IDs of parent items (validator checks they exist)
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
