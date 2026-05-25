package renderer

// ArchDescriptions is the one-line catalog description for each
// ArchLibrary entry. Used by `squiz catalog arch` so agents can pick the
// right system-design icon without reading SVG source.
//
// Keys MUST match ArchLibrary keys exactly — a sync test (see
// arch_descriptions_test.go) catches drift.
//
// One sentence, ~6-12 words, agent-facing. Describe what the icon
// represents (server / database / queue / …), not where to use it.
var ArchDescriptions = map[string]string{
	// Compute
	"server":    "Rack-style server with three slots and power LEDs.",
	"container": "Shipping-container cube with vertical ridge lines.",
	"pod":       "Hexagonal Kubernetes pod with internal spoke lines.",
	"function":  "Rounded square enclosing a lambda glyph.",
	"worker":    "Toothed gear with central hub and spokes.",
	"scheduler": "Clock face flanked by job-bracket pillars.",

	// Data
	"database": "Cylindrical database with stacked rings.",
	"table":    "Ruled table with header row and three columns.",
	"blob":     "Amorphous blob shape with three internal dots.",
	"storage":  "Two stacked disk drives with LED indicators.",
	"cache":    "Boxed lightning bolt indicating fast retrieval.",
	"stream":   "Three parallel wavy lines flowing rightward.",

	// Network
	"load-balancer": "Diamond load-balancer splitting traffic in two directions.",
	"gateway":       "Two pillars with arrow passing between them.",
	"cdn":           "Wireframe globe with four edge node dots.",
	"dns":           "Tag-shaped DNS label with dot and .io text.",
	"vpc":           "Dashed cloud-style boundary enclosing three service boxes.",
	"subnet":        "Inner solid block nested inside dashed subnet outline.",
	"firewall":      "Brick wall with flame icon rising above it.",

	// Services
	"api":   "Curly braces wrapping three dots, JSON-style.",
	"queue": "Three FIFO message envelopes with dequeue arrow.",
	"topic": "Pub/sub source fanning out to three subscribers.",

	// Observability
	"log":    "Lined document page with folded corner.",
	"metric": "Boxed bar chart with four ascending bars.",
	"trace":  "Stacked waterfall trace spans of decreasing width.",

	// Identity
	"user":    "Head-and-shoulders user silhouette outline.",
	"mobile":  "Phone outline with speaker slot and home dot.",
	"browser": "Browser window with tab dots and address bar.",

	// Security
	"secret":   "Padlock with three asterisks inside the body.",
	"key-icon": "Classic key with circular bow and two teeth.",
}

// archCategory maps each ArchLibrary key to its category bucket. Surfaced
// by `squiz catalog arch --json` so agents can filter by group.
var archCategory = map[string]string{
	// Compute
	"server":    "compute",
	"container": "compute",
	"pod":       "compute",
	"function":  "compute",
	"worker":    "compute",
	"scheduler": "compute",

	// Data
	"database": "data",
	"table":    "data",
	"blob":     "data",
	"storage":  "data",
	"cache":    "data",
	"stream":   "data",

	// Network
	"load-balancer": "network",
	"gateway":       "network",
	"cdn":           "network",
	"dns":           "network",
	"vpc":           "network",
	"subnet":        "network",
	"firewall":      "network",

	// Services
	"api":   "services",
	"queue": "services",
	"topic": "services",

	// Observability
	"log":    "observability",
	"metric": "observability",
	"trace":  "observability",

	// Identity
	"user":    "identity",
	"mobile":  "identity",
	"browser": "identity",

	// Security
	"secret":   "security",
	"key-icon": "security",
}
