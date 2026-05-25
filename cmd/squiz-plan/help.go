package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// cmdHelp implements `squiz-plan help [topic]` — the topical-reference
// surface (mirrors `git help <topic>` / `go help <topic>` convention).
//
// Routing nuance: main.go wires BOTH `help` and `--help` to this function.
// `squiz-plan --help` should keep the legacy `printUsage()` behavior; only
// the bare `help` subcommand opts into topical mode. We disambiguate by
// reading os.Args[1] directly, since we cannot modify main.go.
//
// Output discipline:
//   - success (topic list, topic body) → stdout
//   - errors (unknown topic, plus the topic list shown alongside) → stderr
//   - --help-style usage → stderr (matches printUsage)
func cmdHelp(args []string) {
	// Preserve `squiz-plan --help` / `squiz-plan -h` behavior.
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "--help", "-h":
			printUsage()
			return
		}
	}

	// `squiz-plan help --help` explains the help command itself.
	if len(args) > 0 {
		switch args[0] {
		case "--help", "-h":
			printHelpAbout()
			return
		}
	}

	// `squiz-plan help` with no topic → list available topics on stdout.
	if len(args) == 0 {
		printTopicList(os.Stdout)
		return
	}

	topic := args[0]
	body, ok := topics[topic]
	if !ok {
		fmt.Fprintf(os.Stderr, "squiz-plan help: unknown topic %q\n\n", topic)
		printTopicList(os.Stderr)
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, body)
}

// printTopicList writes the "Available topics" block in the format the spec
// shows verbatim (two-space indent, name padded to col ~18, then summary).
func printTopicList(w *os.File) {
	fmt.Fprintln(w, `Available topics. Run "squiz-plan help <topic>" for details.`)
	fmt.Fprintln(w)

	names := make([]string, 0, len(topics))
	for name := range topics {
		names = append(names, name)
	}
	sort.SliceStable(names, func(i, j int) bool {
		return topicOrder[names[i]] < topicOrder[names[j]]
	})

	for _, name := range names {
		summary := topicSummaries[name]
		fmt.Fprintf(w, "  %-16s %s\n", name, summary)
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, `See also: squiz-plan --help (CLI flags), squiz-plan skill (full SKILL.md).`)
}

func printHelpAbout() {
	fmt.Fprintln(os.Stdout, strings.TrimSpace(`
squiz-plan help — topical reference

Usage:
  squiz-plan help                    list available topics
  squiz-plan help <topic>            print the reference block for one topic
  squiz-plan help --help             this message

Topics are stand-alone references — reading `+"`squiz-plan help art`"+` should
tell you everything about art forms without needing to consult SKILL.md.

For the full SKILL.md (the canonical spec the binary was built from), run:

  squiz-plan skill                   # dump SKILL.md to stdout
  squiz-plan skill --out FILE        # write SKILL.md to FILE

For the CLI flag reference (render / example / version / etc.), run:

  squiz-plan --help
`))
}

// topicOrder controls listing order — the headline three first, then the
// plan-specific topics (sections / refs / notes / proposed-items), then the
// option/recommendation pair.
var topicOrder = map[string]int{
	"art":             1,
	"themes":          2,
	"dsl":             3,
	"sections":        4,
	"refs":            5,
	"notes":           6,
	"proposed-items":  7,
	"schema":          8,
	"export":          9,
	"options":         10,
	"recommendations": 11,
}

var topicSummaries = map[string]string{
	"art":             `Art is composition from parts; the grammar, the visual budget, the checklist.`,
	"themes":          `The 8 themes; precedence; auto-rotation behavior.`,
	"dsl":             `The 11 parametric DSL primitives with grammar.`,
	"sections":        `The 6 canonical sections + custom sections; ID prefixes.`,
	"refs":            `Cross-references between items; the upward convention; validation.`,
	"notes":           `The 3 note channels: per-item, per-section, plan-level.`,
	"proposed-items":  `The "+ add item" affordance; how proposed items round-trip.`,
	"schema":          `High-level JSON shape (use "squiz-plan schema" for machine-readable).`,
	"export":          `Shape of the JSON payload the user pastes back.`,
	"options":         `Per-item options: id, label, name, desc, art, recommendation.`,
	"recommendations": `When to set the recommendation field; one-per-item rule.`,
}

// topics holds the long-form reference for each topic. Stand-alone — a
// reader of any one topic should not need to consult SKILL.md.
//
// Parity over abstraction: the squiz-plan topics that overlap with squiz
// (art, themes, dsl, schema, export, options, recommendations) duplicate
// the prose rather than DRY into a shared package. Each binary stays
// self-contained.
var topics = map[string]string{
	"art": `squiz-plan help art — composition with library parts

Art is composition. The library and DSL primitives are parts.

The wf:*, arch:*, and DSL items are pieces you remix inside a custom
drawing — a labeled box of your own words, an arrow you draw between two
icons, a small composition that shows the specific idea the item is
about. The default mode is composition: raw SVG that embeds one or more
library references plus your own text and connecting marks. Treat a single
library token alone as a deliberate exception, not the norm.

THE COMPOSITION MECHANISM (v0.8.0+)

Inside raw SVG, write <use href="wf:NAME" .../> or <use href="arch:NAME" .../>
to drop a library shape in. The renderer inlines it as a <symbol> and you
get a bespoke composition. Library names resolve across both namespaces.

  {
    "art": "<svg viewBox='0 0 100 60' xmlns='http://www.w3.org/2000/svg' style='width:80%;height:auto'>
      <use href='arch:browser' x='2'  y='18' width='22' height='22'/>
      <use href='arch:cache'   x='40' y='18' width='22' height='22'/>
      <use href='arch:database' x='76' y='18' width='22' height='22'/>
      <line x1='25' y1='29' x2='39' y2='29' stroke='var(--ink-3)' stroke-width='1.2'/>
      <line x1='63' y1='29' x2='75' y2='29' stroke='var(--accent)' stroke-width='1.6'/>
      <text x='50' y='12' text-anchor='middle' font-family='IBM Plex Mono' font-size='7'
            fill='var(--ink-2)'>50ms p99</text>
      <text x='87' y='52' text-anchor='middle' font-family='IBM Plex Mono' font-size='7'
            fill='var(--accent)'>hot</text>
    </svg>"
  }

That's one composition. Three arch icons (parts), one custom mono label,
one accented relationship-line. The picture says "browser → cache → db,
the db edge is the hot path" — which is the item's actual desc, not
generic decoration.

Run "squiz-plan catalog wf|arch|dsl" for the authoritative list of parts.
The --json flag reports each entry's naturalBox so you can size <use>
boxes without guessing.

PICK THE KIND OF DIAGRAM FIRST, THEN COMPOSE

Six shapes consistently read at 100×60:

  Shape                 When it fits                      Typical recipe
  ───────────────────── ───────────────────────────────── ───────────────────────────────
  labeled-object        one thing with one part that      wf/arch icon + callout
                        matters                           pointing at the part
  flow                  a process or path                 2-4 boxes with arrows; one
                        (request → ack, raw → cooked)     highlighted
  contrast / vs         before/after, this/that,          two mini-panels with
                        chosen/rejected                   divider:vs between them
  part-whole            scope, subset, "this slice of     container with one cell
                        that"                             emphasized
  metric-with-context   a number that means something     sparkline + baseline:N
                        only vs a reference               reference line + label
  typed-list            a choice among siblings or        pills:a*|b|c* or N rows
                        named modes                       with one accent

VISUAL BUDGET AT 100×60

Four "ink events" max — one silhouette, one accent, one ≤6-char label,
one relationship-line. More than that and the eye gives up.

One accent color in one place only; if two things are accented, neither
is. ~30-40% of the viewBox should be empty. Asymmetric composition
(centered = decorative; off-center = informative).

AUTHORING CHECKLIST (run before committing each art form)

  1. CAPTION IT. In one sentence, what does the picture I'm about to
     write actually show?

  2. KEY-NOUN MATCH. Does the caption reproduce at least two load-bearing
     nouns or verbs from the desc? If desc says "hot path: browser → cache
     → db, 50ms p99", the caption must say browser/cache/db/50ms.

  3. SIBLING DIFF CHECK. If this is one of N options in a chooser, would
     the same single token be defensible on each? If yes, mine isn't
     specific enough.

  4. COMPOSE-VS-SINGLE TEST. Could a label, arrow, or second icon make
     this read more specifically? If yes, compose; if the primitive
     truly IS the picture, single is fine.

  5. "NONE" CHECK. If after all the above the best I can do is a generic
     library icon, use "art": "none" instead. Generic decoration is
     worse than no art.

WHEN A SINGLE LIBRARY TOKEN IS CORRECT

Only when the primitive IS the picture. Three legitimate cases:

  - "text:\"systemctl enable\\nclipsi.service\"@mono?size=10&color=accent"
    — the literal command IS the thing

  - "pills:M*|T|W*|T|F*|S|S"
    — three lit days IS the desc "3× / week"

  - "wf:calendar-grid" alone
    — only when the item is literally about a calendar grid

If you reach for one library token, finish the sentence "the picture I
want to draw IS this exact icon, alone, with no label." If you can't say
that, compose.

ANTI-PATTERNS THE VALIDATOR WATCHES FOR

  - Two options in a chooser share the same "art" — the icon means
    nothing if siblings collide (warning: sibling-art-collision)
  - Generic arch:user / arch:server / arch:database where the desc
    names something specific — give it a label or compose around it
  - A section where every item uses a single library token and no
    composed SVG (warning: composition-thin)
  - Decoration that "looks plausible" but reproduces zero key nouns
    from the desc

THE 5 FORM SHAPES (string syntax)

The renderer detects form from the string:

  starts with "wf:" or "arch:"          named-library token
  starts with "<svg"                    raw SVG (composition lives here)
  matches a DSL prefix (flow:, box:,    DSL primitive
    text:, pills:, spark:, bars:,
    grid:, swatches:, sample:,
    circle-pack:, arrow:, callout:,
    brace:, divider:, badge:, range:,
    baseline:, times:)
  exactly "none"                        slot collapses
  omitted                               plan items: no art shown
                                        (unlike /squiz options' per-letter
                                        fallback)

Raw SVG is no longer the "escape hatch" — it's the host of composition.
Use viewBox='0 0 100 60', style='width:80%;height:auto', and CSS vars
(var(--accent), var(--ink), var(--ink-2), var(--ink-3), var(--rule),
var(--rule-2)) so the picture inherits the active theme.

Items rendered with the per-section default art are flagged by the
validator — treat the defaults as a sign you skipped composing, not as
good output.

SEE ALSO

  squiz-plan help dsl              full DSL grammar (18 primitives)
  squiz-plan help themes           how art inherits theme colors via CSS vars
  squiz-plan catalog wf            list wf:* parts with naturalBox dimensions
  squiz-plan catalog arch          list arch:* parts with naturalBox dimensions
  squiz-plan catalog dsl           list DSL primitives with grammar
  squiz-plan validate              flags sibling-art-collision, composition-thin,
                                   default-art-fallback
`,

	"themes": `squiz-plan help themes — the 8 themes & auto-rotation

squiz-plan ships 8 retro-terminal × editorial themes shared with the /squiz
sibling. Every rendered plan lives in exactly one theme; the binary picks
it for you by default.

THE 8 THEMES

  paper      Cream + ink + rust accent. Editorial, calm.          light
  phosphor   Green-on-black CRT.                                   dark
  amber      IBM 3279 amber on near-black.                         dark
  beige      PS/2 cream with IBM blue.                             light
  rose       Warm pink, plum ink, rose accent.                     light
  ocean      Pale blue-grey, deep teal, coral accent.              light
  forest     Oat cream, moss, warm gold.                           light
  slate      Cool dark grey, electric blue accent.                 dark

AUTO-ROTATION (default behavior)

When you do NOT set a theme, the binary auto-assigns one per repo,
sequentially, persisted across runs. Each new repo you render a plan in
gets the next theme in the rotation. Re-renders of the same repo always
get the same theme — so a project's identity stays stable across multiple
plans and squizzes.

The cache lives at ~/.squiz/themes.json (shared with /squiz). The repo
key is:

  1. The output of "git remote get-url origin" if available, else
  2. The absolute working directory.

PRECEDENCE (highest wins)

  1. --theme <name> on the CLI
  2. "theme" field in index.json
  3. Auto-derived from repo cache

GUIDANCE

Omit "theme" from index.json. Let it auto. Only set it explicitly when you
want a specific identity for one plan — for example, "phosphor" for an
infra plan, "rose" for a brand/marketing one. Setting a theme in JSON does
NOT update the repo cache; it's a one-doc override.

OPTIONAL TOGGLES (in index.json)

  "scanlines": true            CRT-character scanline overlay; best with
                               phosphor / amber / slate.
  "density":  "comfortable"    Roomier line-heights and gutters. Default
                               is "compact".
  "cursor":   true             Blinking cursor in the wordmark.

EXAMPLES

  { /* theme omitted */ }                              auto-rotate
  { "theme": "phosphor" }                              one-doc override
  squiz-plan render plan/index.json --theme amber      CLI override

SEE ALSO

  squiz-plan help art              themed art inherits CSS vars
  squiz-plan help schema           where theme / density / scanlines live
`,

	"dsl": `squiz-plan help dsl — the 11 parametric DSL primitives

Same DSL as /squiz. Use it when the wf:* / arch:* library doesn't ship the
exact shape you need, or when you want to compose a fresh idea.

All DSL output is theme-aware — strokes, fills, and text inherit CSS vars
(--ink, --ink-2, --accent, --rule, --rule-2) so the same string renders
differently across the 8 themes.

THE 11 PRIMITIVES

| Form                              | Example                                              | Renders                                  |
|-----------------------------------|------------------------------------------------------|------------------------------------------|
| grid:NxM[@RATE]                   | "grid:7x7@0.55"                                      | N×M heatmap; RATE in [0,1] fills cells   |
| spark:[V,V,V,...]                 | "spark:[3,5,4,7,6,9,11]"                             | sparkline from data series               |
| bars:[V,V,V,...]                  | "bars:[3,5,4,7,6,9,11]"                              | bar chart                                |
| swatches:#A,#B,...                | "swatches:#f1ebde,#1a1814,#b34a1a"                   | palette swatches (literal hex colors)    |
| pills:A*|B|C*                     | "pills:morning*|midday|evening*"                     | chip row; trailing * marks active        |
| sample:"text"[@FONT]              | "sample:\"Quiet welcome.\"@serif"                    | styled sample text                       |
| circle-pack:N                     | "circle-pack:12"                                     | N organically-arranged circles           |
| text:"a\nb"[@FONT][?opts]         | "text:\"Quiet\\nwelcome.\"@mono?size=18&color=accent" | rich multi-line styled text             |
| flow:[a,b,c]                      | "flow:[client?icon=user,api?icon=api,db?icon=database]" | left-to-right pipeline of named boxes |
| box:label[?icon=ARCH]             | "box:web-tier?icon=server"                           | single labeled box with optional icon    |
| arrow:"label"[?dir=DIR]           | "arrow:\"async\"?dir=down"                           | standalone labeled arrow glyph           |

FONT OPTIONS (sample: and text:)

  serif        editorial Plex Serif
  sans         Plex Sans (default)
  mono         Plex Mono

TEXT: QUERY OPTIONS

  ?size=N             6-36; default 14
  ?align=A            left | center | right; default left
  ?weight=W           300-700; default 400
  ?color=C            ink | ink-2 | ink-3 | accent | rule | rule-2;
                      default ink

  Combine with & — e.g. ?size=18&align=center&color=accent&weight=700.
  Multi-line text uses \n inside the quoted string.

FLOW: ICON EMBEDDING

  flow:[a,b,c]                            three plain boxes with arrows
  flow:[a?icon=user,b?icon=api,c?icon=db] same, each box gets an arch icon

  Icons must be valid arch:* names. Especially useful for engineering and
  build items where the diagram IS the explanation.

ARROW: DIRECTIONS

  ?dir=right    default; left-to-right
  ?dir=down     top-to-bottom
  ?dir=up       bottom-to-top
  ?dir=left     right-to-left

JSON ESCAPING

DSL strings live inside JSON, so backslashes and quotes need JSON-escaping:

  // JSON-source                                  Effective DSL
  "art": "text:\"hi\\nthere\"@mono"            →  text:"hi\nthere"@mono
  "art": "sample:\"Quiet.\"@serif"             →  sample:"Quiet."@serif

PLAN-SPECIFIC TIPS

flow: and box: shine in engineering and build sections, where a small
diagram says more than a paragraph. arrow: is rare on its own; reach for
it when you want to show coupling between two adjacent items in a tab.

SEE ALSO

  squiz-plan help art              when to pick DSL vs library vs raw SVG
  squiz-plan catalog               every wf:* / arch:* name DSL can use
`,

	"sections": `squiz-plan help sections — the 6 canonical sections + custom sections

A plan is a multi-file JSON tree. index.json lists which sections to load;
the binary reads each <sectionId>.json sibling, validates IDs and refs,
and renders one tabbed HTML doc with one tab per section.

THE 6 CANONICAL SECTIONS

  Section         Prefix    Purpose
  ─────────────── ────────  ──────────────────────────────────────────────
  overview        OVR       Mission, audience, hard constraints, success
                            criteria. The "why are we doing this" tab.
  functional      FR        What the system does. User-facing behavior.
  non-functional  NFR       How the system behaves — performance, security,
                            offline, accessibility, etc.
  cases           CASE      Concrete real-world scenarios the system has
                            to handle. These make the plan feel real.
  engineering     ENG       Architecture decisions, components, the
                            shape of the build.
  build           BUILD     Concrete steps, per component, to actually
                            deliver. The "tickets".

Canonical sections always render in canonical order in the tab strip,
regardless of position in index.json's "sections" array.

ID PREFIX CONVENTION

Each item's id MUST start with the section's prefix:

  overview        → OVR-1, OVR-2, OVR-audience
  functional      → FR-1, FR-2.3, FR-export
  non-functional  → NFR-1, NFR-perf
  cases           → CASE-1, CASE-bedroom-freeze
  engineering     → ENG-1, ENG-storage
  build           → BUILD-1, BUILD-firmware

The validator rejects items where the prefix doesn't match the section
filename. This is what makes refs work — a "[Engineering · ENG-2]" badge
is unambiguous because ENG- ids only live in engineering.json.

The suffix after the prefix is a stable slug. Numeric (FR-1, FR-2.3) is
fine for ordered/sequenced items; semantic (CASE-bedroom-freeze) is
better when an item has a name that won't change. Both round-trip back
in the user's feedback payload as feedback[].id, so pick something that
still makes sense in a month.

CUSTOM SECTIONS

You may append custom sections (glossary, risks, appendix, alternatives,
out-of-scope, …) after the canonical six. Add the id to index.json's
"sections" array; the binary will look for <sectionId>.json beside
index.json.

  // index.json
  {
    "sections": [
      "overview", "functional", "non-functional",
      "cases", "engineering", "build",
      "risks", "glossary"
    ]
  }

Custom sections render as tabs AFTER the canonical six, in the order
declared. Their items don't get a built-in ID prefix convention — pick
one that makes sense (RISK-1, GLOSS-database) and the validator will
enforce it as long as you're consistent.

SECTION SIZING TARGETS

  overview         3-5 items   mission, audience, constraints, success
  functional       3-8 items   what it does (user-facing)
  non-functional   2-5 items   how it behaves (cross-cutting)
  cases            2-5 items   concrete scenarios; the plan's pulse
  engineering      3-8 items   architecture; refs functional + non-func
  build            3-10 items  concrete work; refs engineering

Fewer items than the floor → the section feels empty. More than the
ceiling → the section becomes a wall the user has to wade through.
Split or merge if you bust the range.

EXAMPLES

  // plan/overview.json
  {
    "items": [
      { "id": "OVR-1", "title": "Mission",   "desc": "Log home temperatures every minute." },
      { "id": "OVR-2", "title": "Audience",  "desc": "One household; the owner is the operator." },
      { "id": "OVR-3", "title": "Constraint","desc": "Total cost ≤ $50 of hardware; one $5/mo VPS." }
    ]
  }

  // plan/engineering.json
  {
    "items": [
      {
        "id":    "ENG-1",
        "title": "Sensor driver (BLE)",
        "desc":  "Read 4 cheap Bluetooth LE thermometers. Pi scans every 60 s.",
        "art":   "wf:dot-trend",
        "refs":  ["FR-1", "NFR-1"]
      }
    ]
  }

SEE ALSO

  squiz-plan help refs             cross-references between items
  squiz-plan help schema           the per-section JSON shape
  squiz-plan example               an end-to-end sample plan tree
`,

	"refs": `squiz-plan help refs — cross-references between items

Most plans collapse under their own weight: by the time you read step 17 in
the build section, you've forgotten which functional requirement motivated
it. The "refs" field solves that by making the thread between layers
visible and clickable.

THE FIELD

Every item carries an optional "refs" array of parent item IDs:

  {
    "id":    "ENG-1",
    "title": "Sensor driver (BLE)",
    "desc":  "Read from 4 cheap Bluetooth LE thermometers.",
    "refs":  ["FR-1", "NFR-1"]
  }

WHAT THE RENDERER DOES

For each ref ID, the renderer:

  1. Validates the ID exists somewhere in the plan. Missing refs fail
     "squiz-plan validate" and "squiz-plan render".
  2. Renders an inline badge with the parent's section label:
       [Functional · FR-1]   [Non-functional · NFR-1]
  3. Wires the badge as a tab-switch link. Click it and the renderer
     switches tabs, scrolls to the parent item, and highlights it briefly.
  4. The browser back button returns to the previous item — so a user can
     follow a chain (BUILD-3 → ENG-1 → FR-1 → OVR-2) and then pop back
     without losing their place.

THE UPWARD CONVENTION

By convention, refs point UPWARD in the spine:

  build         →  engineering
  engineering   →  functional, non-functional, cases
  functional    →  overview
  non-functional →  overview
  cases         →  functional, overview

You CAN sideways-ref (an ENG-* item refs another ENG-* item) but use it
sparingly — too many sideways refs make the plan feel knotty.

You SHOULD NOT ref downward (an FR-* item refs an ENG-* item). It inverts
the dependency direction the reader expects and the badges become noise.
If you find yourself wanting to ref downward, restate the parent item
instead.

VALIDATION

"squiz-plan validate plan/index.json" walks the tree and rejects:

  - refs to IDs that don't exist anywhere in the plan
  - refs to IDs in section files that aren't listed in index.json's
    "sections" array
  - circular refs (A refs B refs A)

Run validate before render when you change IDs — a typo'd ref is the
most common breakage.

REFS AS DOCUMENTATION

A good "refs" list is itself documentation: it tells the reviewer
"this item exists because of these upstream commitments". When refs are
sparse, the reviewer wonders what motivated the item. When refs are
comprehensive, the reviewer can audit the plan section-by-section
without losing context.

EXAMPLES

  // engineering item refs the functional + non-functional reqs that
  // motivated it. Clicking either badge in the rendered HTML switches
  // to that tab and highlights the parent.
  {
    "id":    "ENG-2",
    "title": "Local storage engine",
    "desc":  "Embedded store the Pi writes readings into.",
    "refs":  ["FR-1", "NFR-2"]
  }

  // build item refs the engineering item it implements.
  {
    "id":    "BUILD-3",
    "title": "Wire SQLite + WAL mode",
    "desc":  "Initialize DB on first boot. WAL for concurrent readers.",
    "refs":  ["ENG-2"]
  }

SEE ALSO

  squiz-plan help sections         what each prefix means
  squiz-plan help schema           where refs live in the JSON shape
  squiz-plan validate              the validator that enforces refs
`,

	"notes": `squiz-plan help notes — the 3 note channels

The rendered HTML gives the user THREE places to leave notes. Each channel
serves a different scope; each one round-trips back in the export payload
under a distinct field. Apply them as instructions for the next iteration.

THE 3 CHANNELS

1. PER-ITEM NOTES   (scope: one item)

   A textarea inside each item's feedback widget, alongside the
   ✓ approve / ? question / ✗ reject buttons. Use for feedback that
   targets exactly one item:

     "20 minutes feels long — most pipe-freezing scenarios are faster."
     "Wire SQLite + WAL mode is fine but rename it; WAL is the default."

   Round-trips as: feedback[].note

2. PER-SECTION NOTES   (scope: one tab)

   A sticky textarea at the top of each tab. Use for feedback that spans
   multiple items in a section, or for gaps where the right answer is to
   add a new item:

     "These FRs are missing the export workflow."
     "Consider folding ENG-3 into ENG-2."
     "The non-functional section needs an offline-availability item."

   Round-trips as: section_notes.<sectionId>

3. PLAN-LEVEL NOTES   (scope: whole plan)

   A single textarea inside the copy-json modal. Use for overall direction
   that doesn't fit any tab:

     "This scope is too ambitious for v1 — see section notes."
     "Looks good overall; let's ship it. Address the questioned items
      in the next sprint."

   Round-trips as: plan_note (single string, may be null)

ROUND-TRIP SHAPE

The export payload bundles all three channels:

  {
    "feedback": [
      {
        "id":     "FR-3",
        "status": "questioned",
        "note":   "20 minutes feels long.",                  // PER-ITEM
        ...
      }
    ],
    "section_notes": {                                       // PER-SECTION
      "functional":  "Missing the export workflow.",
      "engineering": "ENG-3 could fold into ENG-2."
    },
    "plan_note": "Overall: too ambitious for v1."            // PLAN-LEVEL
  }

WHEN THE USER PASTES BACK

Treat each channel as instructions for the next iteration of the plan:

  - PER-ITEM notes → revise the corresponding item, then re-render.
    If the note questions a fact, push back in chat before changing the
    item — the user may be wrong.

  - PER-SECTION notes → consider restructuring the section. Add missing
    items, fold redundant items, reorder. Per-section notes are often
    where new items get proposed (see also: squiz-plan help proposed-items).

  - PLAN-LEVEL note → re-evaluate the plan's scope or shape. If it says
    "too ambitious", consider splitting into v1/v2; if it says "looks
    good", proceed to the build phase.

EMPTY VS NULL

A channel the user didn't touch comes back as:

  - feedback[].note          → null (or the field is omitted)
  - section_notes.<id>       → omitted from the section_notes object
  - plan_note                → null

So you can distinguish "user left it blank intentionally" from "user
never opened the widget" by checking presence vs null. In practice,
treat both the same way (the user has nothing to add).

GUIDANCE FOR THE AGENT

  - Don't ignore notes. Even a single per-item note can flip the
    interpretation of a "?" status.
  - Quote the note back in your reply ("You said FR-3 felt long — I've
    cut it from 20 to 10 minutes; see the regenerated tab"). It shows
    the user their note landed.
  - If notes contradict each other (per-section says "ship as-is",
    plan-level says "too ambitious"), surface the conflict and ask.

SEE ALSO

  squiz-plan help export           full export payload shape
  squiz-plan help proposed-items   the "+ add item" affordance
`,

	"proposed-items": `squiz-plan help proposed-items — the "+ add item" affordance

Each tab in the rendered plan has a "+ add item" button. Clicking it opens
a small form (title, desc, optional refs) where the user can propose a
brand-new item to add to that section.

Proposed items come back in the export as a typed proposed_items[] array.
Apply them as suggestions when you regenerate the plan.

THE EXPORT SHAPE

  {
    ...
    "proposed_items": [
      {
        "section": "functional",
        "title":   "Rate-limit alerts",
        "desc":    "Throttle to at most 1 alert per room per hour.",
        "refs":    ["NFR-2"]
      },
      {
        "section": "build",
        "title":   "CSV export of historical readings",
        "desc":    "Endpoint that streams the last 30 days as CSV.",
        "refs":    ["FR-3", "ENG-2"]
      }
    ]
  }

FIELDS

  section      REQUIRED. The id of the section the user proposed the
               item into. Matches one of the entries in index.json's
               "sections" array.

  title        REQUIRED. The headline. Treat as the user's first draft —
               you may sharpen it when regenerating, but keep the intent.

  desc         OPTIONAL. The 1-2 sentence body. May be empty if the
               user only typed a title.

  refs         OPTIONAL. The parent IDs the user attached. May be empty.
               Validate these against the existing plan before adding —
               a typo here is the user's, not yours.

APPLYING PROPOSED ITEMS

When regenerating:

  1. For each proposed item, generate an ID with the section's prefix.
     Numbering: use the next free integer in that section (if FR-1 …
     FR-5 exist, the new one is FR-6). Semantic IDs are also fine if
     the title suggests one (CASE-bedroom-freeze).

  2. Validate the user's refs. Invalid refs (typos, IDs that no longer
     exist) should be dropped silently OR surfaced in chat — your call
     based on whether the rest of the proposal still makes sense
     without them.

  3. Sharpen the title and desc if the user's draft is rough. Keep the
     intent; tighten the prose. Don't add fields the user didn't ask
     for (no art, no options, no recommendation) unless the item
     obviously calls for one.

  4. Insert into the section file at a sensible position — usually at
     the end, unless the user's "refs" suggest it belongs grouped with
     a particular cluster.

  5. Re-render. Mention the new items in your reply so the user can see
     their proposals landed:

       "Added two items from your proposals:
        - FR-6: Rate-limit alerts (refs NFR-2)
        - BUILD-9: CSV export of historical readings (refs FR-3, ENG-2)"

WHEN TO PUSH BACK

Not every proposal should be accepted as-is. Push back in chat (don't
silently drop) when:

  - The proposal duplicates an existing item with different wording.
    Point at the existing item and ask if the user wants to edit that
    one instead.
  - The proposal contradicts a hard constraint (an OVR- item). Surface
    the conflict; the user may want to relax the constraint or drop
    the proposal.
  - The proposal would land in a section it doesn't belong to (a build
    step proposed into overview). Suggest the right section.

EMPTY ARRAY

If the user proposed nothing, the field is either absent or an empty
array. Either way, nothing to apply.

SEE ALSO

  squiz-plan help notes            section-level notes often pair with
                                   proposed items (the note explains
                                   the rationale; the proposal is the fix)
  squiz-plan help export           full export payload shape
  squiz-plan help sections         section IDs and prefixes
`,

	"schema": `squiz-plan help schema — high-level JSON shape

A plan is a multi-file JSON tree in one directory:

  plan/
  ├── index.json              ← top-level descriptor (mandatory)
  ├── overview.json
  ├── functional.json
  ├── non-functional.json
  ├── cases.json
  ├── engineering.json
  └── build.json

For the authoritative machine-readable JSON Schema, run:

  squiz-plan schema             # write JSON Schema to stdout

For an end-to-end sample plan tree you can copy and modify, run:

  squiz-plan example            # scaffolds squiz-plan-example/ in cwd

INDEX.JSON — the top-level descriptor

  {
    "title":     "ThermoLog — home temperature logger",  // page H1
    "lede":      "Six load-bearing sections, traceable…", // one-liner
    "theme":     "paper",        // OPTIONAL — omit for auto-rotation
    "density":   "compact",      // OPTIONAL — compact | comfortable
    "scanlines": false,          // OPTIONAL — CRT overlay
    "cursor":    true,           // OPTIONAL — blinking cursor in wordmark
    "sections": [
      "overview",
      "functional",
      "non-functional",
      "cases",
      "engineering",
      "build"
      // append "glossary" / "risks" / etc. for custom sections
    ]
  }

Canonical sections render in canonical order regardless of position in
"sections". Custom sections append in declaration order.

<SECTIONID>.JSON — the per-section payload

Each section file declares an "items" array. The shape applies to all six
canonical sections and to custom sections:

  {
    "items": [
      {
        "id":    "ENG-1",                  // REQUIRED. Must start with the
                                           // section's prefix (ENG- for
                                           // engineering.json, FR- for
                                           // functional.json, etc.)
        "title": "Sensor driver (BLE)",    // REQUIRED.
        "desc":  "Read 4 BLE thermometers. Pi scans every 60 s.", // REQUIRED.
        "art":   "wf:dot-trend",           // OPTIONAL. See "help art".
        "refs":  ["FR-1", "NFR-1"],        // OPTIONAL. Parent IDs; validated.
        "options": [ /* … */ ]             // OPTIONAL. See "help options".
      }
    ]
  }

VALIDATION

"squiz-plan validate plan/index.json" rejects:

  - A section listed in index.json but missing its sibling JSON file.
  - An item whose id doesn't start with the section's prefix.
  - Duplicate ids anywhere in the plan.
  - refs to ids that don't exist.
  - Unknown wf:/arch: names or malformed DSL strings in art.

Run validate before render when you change ids — a typo'd ref is the
most common breakage.

WHERE THE EXTRAS LIVE

  - title, lede, theme, density, scanlines, cursor, sections   index.json
  - items[] (per section)                                      <sectionId>.json
  - section_notes / plan_note / proposed_items                 export only

The export payload is documented in "squiz-plan help export".

SEE ALSO

  squiz-plan schema                machine-readable JSON Schema
  squiz-plan example               end-to-end sample plan tree
  squiz-plan help sections         what each section is for
  squiz-plan help refs             how cross-refs work
  squiz-plan help options          per-item options (decisions)
  squiz-plan help export           what the user pastes back
`,

	"export": `squiz-plan help export — the JSON the user pastes back

After the user reads the plan and clicks ✓ / ? / ✗ on items (leaving notes
or proposing new ones), they click the sticky "copy json" button at the
bottom of the modal. That copies a single JSON document to the clipboard;
they paste it back into chat for you to parse.

SHAPE

  {
    "plan": "ThermoLog — home temperature logger",
    "source": {
      "file":     "C:\\dev\\thermolog\\plan\\index.html",
      "basename": "index.html"
    },
    "generatedAt": "2026-05-25T12:34:56Z",
    "feedback": [
      {
        "id":     "ENG-2",
        "status": "approved",
        "anchor": "#item-ENG-2",
        "note":   "Like the SQLite pick.",
        "edits":  null,
        "chose":  "sqlite"
      },
      {
        "id":     "FR-3",
        "status": "questioned",
        "anchor": "#item-FR-3",
        "note":   "20 minutes feels long.",
        "edits":  null,
        "chose":  null
      },
      {
        "id":     "BUILD-2",
        "status": "approved",
        "anchor": "#item-BUILD-2",
        "note":   null,
        "edits":  { "title": "Pi firmware (Go single binary)" },
        "chose":  null
      }
    ],
    "section_notes": {
      "functional":  "Missing the export workflow.",
      "engineering": "ENG-3 could fold into ENG-2."
    },
    "plan_note": "Overall: too ambitious for v1.",
    "proposed_items": [
      {
        "section": "functional",
        "title":   "Rate-limit alerts",
        "desc":    "Throttle to at most 1 alert per room per hour.",
        "refs":    ["NFR-2"]
      }
    ],
    "summary": {
      "total": 21,
      "approved": 18, "questioned": 2, "rejected": 0,
      "withNotes": 4, "withEdits": 1, "withChose": 3,
      "sectionsWithNotes": 2, "hasPlanNote": true, "proposedItems": 1
    }
  }

FIELDS

  plan                 Echo of the plan title.

  source.file          Absolute path of the rendered HTML on the user's
                       machine.
  source.basename      Just the filename.

  generatedAt          ISO-8601 timestamp the user clicked "copy json".

  feedback[]           One entry per item the user TOUCHED. Items the
                       user didn't engage with don't appear.

  feedback[].id        The stable item id you authored.
  feedback[].status    "approved" | "questioned" | "rejected".
  feedback[].anchor    "#item-<id>" — append to source.file to deep-link.
  feedback[].note      Free-text per-item notes (string or null).
  feedback[].edits     Sparse object of field overrides (or null).
                       Apply as SUGGESTIONS, not authoritative truth.
                       Common keys: { "title": "…", "desc": "…" }.
  feedback[].chose     For items with an "options" array, the id of
                       the picked option. null otherwise (item without
                       options, or option not picked).

  section_notes        Map keyed by section id. See "help notes".

  plan_note            Single string or null. See "help notes".

  proposed_items[]     User-proposed new items. See "help proposed-items".

  summary              Aggregate counts. Useful for quick situational
                       awareness ("18/21 approved, 2 questions, 1 edit,
                       1 proposed item, no rejects").

INTERPRETING STATUSES

  approved     The item is good as-is. Don't change it unless `+"`edits`"+`
               carries a sharpening.

  questioned   The item is mostly right but the user has a doubt. The
               `+"`note`"+` field is where they explain — read it carefully
               and either revise or push back in chat.

  rejected     The user thinks the item should not exist or is fundamentally
               wrong. Don't silently drop it; reply in chat to confirm
               (it may have implicit dependents elsewhere in the plan).

ANCHOR USAGE

Combine source.file + anchor to deep-link back to the exact item:

  file:///C:/dev/thermolog/plan/index.html#item-ENG-2

Useful when you want to quote the user's feedback back to them with a
clickable reference.

SEE ALSO

  squiz-plan help schema           input shape (what you author)
  squiz-plan help notes            the 3 note channels
  squiz-plan help proposed-items   the "+ add item" affordance
  squiz-plan help options          how "chose" round-trips item options
`,

	"options": `squiz-plan help options — per-item options (decisions)

Most items in a plan are flat statements: "We will use SQLite." But some
items represent UNSETTLED decisions — "Which embedded store should we use?"
For those, attach an "options" array to the item; the rendered card turns
into a chooser. The user's pick comes back in the export as feedback[].chose.

ITEM-WITH-OPTIONS SHAPE

  {
    "id":    "ENG-2",
    "title": "Local storage engine",
    "desc":  "Pick the embedded store the Pi writes readings into.",
    "refs":  ["FR-1", "NFR-2"],
    "options": [
      {
        "id":             "sqlite",
        "name":           "SQLite (single file)",
        "desc":           "Boring, durable, queryable. ~2 MB footprint.",
        "art":            "wf:file-icons",
        "recommendation": "NFR-2 requires durability without an external dep; SQLite hits both targets and the audience can debug it without learning anything new."
      },
      {
        "id":   "bbolt",
        "name": "bbolt (pure-Go KV)",
        "desc": "No SQL, no cgo. Range scans over time-keyed buckets.",
        "art":  "arch:database"
      },
      {
        "id":   "flatfile",
        "name": "Append-only JSONL",
        "desc": "One file per day. Rotate + gzip nightly. Trivial to ship.",
        "art":  "wf:mono-sample"
      }
    ]
  }

OPTION FIELDS

  id                REQUIRED. Stable slug. Comes back in feedback[].chose.
                    Unique WITHIN this item — collisions across items are
                    fine (two items can both have an "id": "sqlite" option;
                    the parent item's id disambiguates).

  label             OPTIONAL. Display string (typically "Option A",
                    "Option B"). Auto-derived from index if omitted.

  name              REQUIRED. Short headline shown big on the card.

  desc              REQUIRED. The 1-2 sentence trade-off. Lean on concrete
                    numbers and consequences.

  art               OPTIONAL. Same forms as item-level art. See
                    "squiz-plan help art".

  recommendation    OPTIONAL. One or two sentences explaining why this is
                    the preferred pick. Renders a ★ RECOMMENDED chip and
                    editorial callout. At most one per item. See
                    "squiz-plan help recommendations".

WHEN TO USE OPTIONS VS A FLAT ITEM

Use options when the decision is UNSETTLED — you and the user need to
pick from real alternatives. Three signals:

  - You'd write "we COULD use X, Y, or Z" in plain prose.
  - The reviewer should see the alternatives (so they can validate
    that the picked option really is best).
  - The choice has downstream consequences another reviewer would
    want to audit.

Use a flat item (no options) when the decision is SETTLED — you've
already picked the answer based on the spec/constraints. Padding every
item with options the user has no real basis to pick between is
busy-work for them.

THE ROUND-TRIP

When the user picks "sqlite", the export carries:

  { "id": "ENG-2", "status": "approved", "chose": "sqlite", ... }

When the user doesn't pick (or the item has no options), "chose": null.

When you regenerate the plan, replace the options block with a flat
"desc" that locks in the picked option (and removes the others). The
plan becomes more concrete with each iteration.

GUIDANCE

Keep "name" terse, "desc" concrete. Vague descs ("flexible, scalable")
give the user nothing to decide on. Specific descs ("$5/mo ceiling,
100 ms p50, single-region") are what make options work.

Use 2-5 options per item. One option is not a decision; six-plus is a
menu. Pre-filter offline if you have a long list.

SEE ALSO

  squiz-plan help art              the "art" field across all 5 forms
  squiz-plan help recommendations  when and how to set "recommendation"
  squiz-plan help schema           where options live in the JSON shape
  squiz-plan help export           how "chose" round-trips back
`,

	"recommendations": `squiz-plan help recommendations — the "recommendation" field

Any item-option can carry an optional "recommendation" field — a one-or-
two-sentence rationale for why this is the preferred choice given the
plan's overview, refs, or constraints.

  {
    "id":             "sqlite",
    "name":           "SQLite (single file)",
    "desc":           "Boring, durable, queryable. ~2 MB binary footprint.",
    "recommendation": "OVR-3's $5-VPS constraint rules out external deps; SQLite hits NFR-2 (durability) without one."
  }

The renderer responds in two ways:

  1. A ★ RECOMMENDED chip appears next to the option's name.
  2. The recommendation string renders as an editorial callout under desc.

The user is still free to pick a different option. Recommendations are
advisory, not constraints.

WHEN TO SET IT

Set "recommendation" when the plan's spine genuinely points at one option.
Concretely:

  - You can cite a specific ref or constraint in the rationale
    ("OVR-3's $5-VPS constraint", "NFR-2 requires durability").
  - The other options have real downsides that the recommended one
    avoids.
  - A reasonable reviewer would agree the recommendation is earned.

DON'T set it when:

  - You don't have a real preference. Marking every option as
    "recommended" defeats the purpose (the chip stops carrying signal).
  - The rationale is fluff ("flexible, modern, popular"). If you can't
    cite a specific constraint or ref, the recommendation isn't earned.
  - You're hedging. A wishy-washy "either is fine" recommendation is
    worse than none — it nudges the user without committing.

THE ONE-PER-ITEM RULE

Set "recommendation" on AT MOST ONE option per item. Multiple
recommendations within the same item make the signal meaningless.

If two options genuinely look equal, that's a sign the item is
under-constrained — add a ref upward to a constraint that breaks the
tie, or split the item into two distinct decisions, before resorting to
two recommendations.

CITING REFS / CONSTRAINTS

The strongest recommendations cite specific upstream commitments. In a
plan, that almost always means citing OVR-, FR-, or NFR- IDs directly:

  "OVR-3 caps us at one VPS; NFR-2 requires durability without an
  external dep. SQLite hits both."

  "FR-1 says the system must work offline for 24h; ENG-3's local cache
  is the only option that survives a router reboot."

Specific is better than generic. "We recommend SQLite because it's
popular" is unearned; "We recommend SQLite because OVR-3 caps us at
$5/mo and NFR-2 requires durability — SQLite hits both" is earned.

GUIDANCE FOR LENGTH

One or two sentences. Three is too many — the callout overshadows the
desc. If you need more justification, the recommendation isn't earned;
restructure the item.

SEE ALSO

  squiz-plan help options          the field in the broader option shape
  squiz-plan help refs             refs to cite in the rationale
  squiz-plan help schema           where options live in JSON
  squiz-plan help export           recommendations don't show in export;
                                   "chose" reflects what the user picked
`,
}
