package renderer

// WFDescriptions is the one-line catalog description for each WFLibrary
// entry. Used by `squiz catalog wf` and `squiz-plan catalog wf` so agents
// can pick the right wireframe without reading SVG source.
//
// Keys MUST match WFLibrary keys exactly — a sync test (see
// wf_descriptions_test.go) catches drift if someone adds a wireframe
// without a description (or vice versa).
//
// One sentence, ~6-12 words, agent-facing. Describe what the icon SHOWS,
// not where to use it. Author the descriptions when adding new wireframes;
// `squiz catalog wf` is meaningless without them.
var WFDescriptions = map[string]string{
	// calendars / dates
	"calendar-grid":  "Monthly heatmap, 7×5 cells with varying intensity.",
	"calendar-week":  "Seven labeled day columns with mixed activity bars.",
	"streak-counter": "Big number with row of filled streak dots.",
	"day-strip":      "Row of seven day-letter chips, last one empty.",
	"year-heatmap":   "Wide GitHub-style year contribution heatmap.",
	"time-of-day":    "Sun, half-moon, and crescent moon glyphs side by side.",
	"clock":          "Analog clock face with hour and minute hands.",

	// charts
	"spark-rising": "Upward-trending sparkline with end marker.",
	"spark-flat":   "Steady, level sparkline hovering around the midline.",
	"spark-noisy":  "Volatile zig-zag sparkline with high variance.",
	"bars-up":      "Seven ascending bars climbing left to right.",
	"donut":        "Donut chart with 64% accent slice and label.",
	"gauge":        "Semicircular gauge at 72% with end dot.",
	"dot-trend":    "Scatter dots rising to upper right with trend line.",

	// identities / avatars
	"avatar-single":  "Single circled initial labeled JUST YOU.",
	"avatar-pair":    "Two side-by-side avatar circles, one filled accent.",
	"avatar-circle":  "Five-person accountability circle of avatar dots.",
	"avatar-feed":    "Feed rows of avatar + line + heart count.",
	"avatar-private": "Solo avatar with small padlock badge.",

	// phone screens
	"phone-blank":   "Empty phone frame with notch speaker line.",
	"phone-list":    "Phone screen showing a list of text rows.",
	"phone-card":    "Phone screen with hero card and image thumb.",
	"phone-input":   "Phone screen with text input and accent submit button.",
	"phone-tabs":    "Phone screen with bottom tab bar, first tab active.",
	"phone-onboard": "Phone screen with progress bar and accent CTA button.",
	"phone-stats":   "Phone screen with big stat number and tiny chart.",

	// controls
	"toggle-on":     "Toggle switch in ON position, accent filled.",
	"toggle-off":    "Toggle switch in OFF position, outline only.",
	"button-accent": "Solid accent CONTINUE primary button.",
	"button-ghost":  "Dashed-outline ghost 'maybe later' secondary button.",
	"slider":        "Horizontal slider at 72%, accent fill and knob.",
	"dropdown":      "Select dropdown with placeholder text and chevron.",

	// status
	"badge-new":   "Solid accent rectangle with NEW label.",
	"pill-row":    "Three status pills: active, pending, draft.",
	"snowflake":   "Six-fold geometric snowflake glyph.",
	"lock":        "Solid padlock icon with shackle.",
	"check-large": "Large accent circle with bold checkmark.",

	// typography samples
	"serif-sample": "Italic serif word 'Tide' over EDITORIAL · SERIF label.",
	"sans-sample":  "Bold sans word 'Tide' over MODERN · SANS label.",
	"mono-sample":  "Monospace word 'tide_' over TECHNICAL · MONO label.",

	// connections / graphs
	"graph-force":    "Force-directed graph: central node with linked satellites.",
	"tree-hier":      "Three-level hierarchical tree with branching children.",
	"radial-burst":   "Central node with six radial spokes to outer dots.",
	"matrix-heatmap": "12×7 heatmap matrix of mixed-intensity cells.",

	// metaphors
	"plant-grow": "Three plants of decreasing height showing growth stages.",
	"garden":     "Row of mixed-size circle plants along a baseline.",
	"paper-fold": "Folded-corner document with horizontal text lines.",

	// misc
	"cmd-palette": "Two keycaps ⌘ and K labeled FIND ANYTHING.",
	"text-cursor": "Text input box with 'type here' and blinking caret.",
	"file-icons":  "Three labeled file icons stacked vertically.",
}

// wfCategory maps each WFLibrary key to its category bucket. Surfaced by
// `squiz catalog wf --json` so agents can filter by group.
var wfCategory = map[string]string{
	// calendars / dates
	"calendar-grid":  "calendars",
	"calendar-week":  "calendars",
	"streak-counter": "calendars",
	"day-strip":      "calendars",
	"year-heatmap":   "calendars",
	"time-of-day":    "calendars",
	"clock":          "calendars",

	// charts
	"spark-rising": "charts",
	"spark-flat":   "charts",
	"spark-noisy":  "charts",
	"bars-up":      "charts",
	"donut":        "charts",
	"gauge":        "charts",
	"dot-trend":    "charts",

	// identities
	"avatar-single":  "identities",
	"avatar-pair":    "identities",
	"avatar-circle":  "identities",
	"avatar-feed":    "identities",
	"avatar-private": "identities",

	// phone screens
	"phone-blank":   "phone-screens",
	"phone-list":    "phone-screens",
	"phone-card":    "phone-screens",
	"phone-input":   "phone-screens",
	"phone-tabs":    "phone-screens",
	"phone-onboard": "phone-screens",
	"phone-stats":   "phone-screens",

	// controls
	"toggle-on":     "controls",
	"toggle-off":    "controls",
	"button-accent": "controls",
	"button-ghost":  "controls",
	"slider":        "controls",
	"dropdown":      "controls",

	// status
	"badge-new":   "status",
	"pill-row":    "status",
	"snowflake":   "status",
	"lock":        "status",
	"check-large": "status",

	// typography
	"serif-sample": "typography",
	"sans-sample":  "typography",
	"mono-sample":  "typography",

	// connections
	"graph-force":    "connections",
	"tree-hier":      "connections",
	"radial-burst":   "connections",
	"matrix-heatmap": "connections",

	// metaphors
	"plant-grow": "metaphors",
	"garden":     "metaphors",
	"paper-fold": "metaphors",

	// misc
	"cmd-palette": "misc",
	"text-cursor": "misc",
	"file-icons":  "misc",
}
