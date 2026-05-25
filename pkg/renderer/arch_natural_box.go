package renderer

// ArchNaturalBox is the measured natural footprint for each ArchLibrary
// entry, in viewBox units. Every arch icon is rendered through
// resolveArch which wraps the `<g>` snippet in
// `<svg viewBox='0 0 100 60' style='width:55%;height:auto'>` — so the
// natural width on a 100-unit canvas is 55 and the natural height
// preserves the 100:60 aspect (55 * 60/100 = 33).
//
// Composing agents use this to size `<use href="arch:…"/>` references
// without guessing — e.g. `width = NaturalBox.W, height = NaturalBox.H`.
//
// Keys MUST match ArchLibrary exactly; TestArchNaturalBoxCoverage
// catches drift in both directions.
var ArchNaturalBox = map[string]NaturalBox{
	// ── Compute ─────────────────────────────────────────────────────
	"server":    archBox,
	"container": archBox,
	"pod":       archBox,
	"function":  archBox,
	"worker":    archBox,
	"scheduler": archBox,

	// ── Data ────────────────────────────────────────────────────────
	"database": archBox,
	"table":    archBox,
	"blob":     archBox,
	"storage":  archBox,
	"cache":    archBox,
	"stream":   archBox,

	// ── Network ─────────────────────────────────────────────────────
	"load-balancer": archBox,
	"gateway":       archBox,
	"cdn":           archBox,
	"dns":           archBox,
	"vpc":           archBox,
	"subnet":        archBox,
	"firewall":      archBox,

	// ── Services ────────────────────────────────────────────────────
	"api":   archBox,
	"queue": archBox,
	"topic": archBox,

	// ── Observability ───────────────────────────────────────────────
	"log":    archBox,
	"metric": archBox,
	"trace":  archBox,

	// ── Identity ────────────────────────────────────────────────────
	"user":    archBox,
	"mobile":  archBox,
	"browser": archBox,

	// ── Security ────────────────────────────────────────────────────
	"secret":   archBox,
	"key-icon": archBox,
}

// archBox is the shared natural footprint for every arch:* icon —
// resolveArch wraps each `<g>` in a `width:55%`-styled SVG, so they all
// share the same W/H/aspect by construction. Pulled into a constant so
// the table stays readable and a one-line edit retunes them all.
var archBox = NaturalBox{W: 55, H: 33, AspectRatio: 55.0 / 33.0}
