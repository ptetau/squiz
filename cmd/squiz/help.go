package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// cmdHelp implements `squiz help [topic]` — the topical-reference surface
// (mirrors `git help <topic>` / `go help <topic>` convention).
//
// Routing nuance: main.go wires BOTH `help` and `--help` to this function.
// `squiz --help` should keep the legacy `printUsage()` behavior; only the
// bare `help` subcommand opts into topical mode. We disambiguate by reading
// os.Args[1] directly, since we cannot modify main.go.
//
// Output discipline:
//   - success (topic list, topic body) → stdout
//   - errors (unknown topic, plus the topic list shown alongside) → stderr
//   - --help-style usage → stderr (matches printUsage)
func cmdHelp(args []string) {
	// Preserve `squiz --help` / `squiz -h` behavior — those go through the
	// printUsage() flag-summary surface that's been here since v0.1.
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "--help", "-h":
			printUsage()
			return
		}
	}

	// `squiz help --help` (or `-h`) explains the help command itself.
	if len(args) > 0 {
		switch args[0] {
		case "--help", "-h":
			printHelpAbout()
			return
		}
	}

	// `squiz help` with no topic → list available topics on stdout.
	if len(args) == 0 {
		printTopicList(os.Stdout)
		return
	}

	topic := args[0]
	body, ok := topics[topic]
	if !ok {
		fmt.Fprintf(os.Stderr, "squiz help: unknown topic %q\n\n", topic)
		printTopicList(os.Stderr)
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, body)
}

// printTopicList writes the "Available topics" block in the format the spec
// shows verbatim (two-space indent, name padded to col ~18, then summary).
func printTopicList(w *os.File) {
	fmt.Fprintln(w, `Available topics. Run "squiz help <topic>" for details.`)
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
	fmt.Fprintln(w, `See also: squiz --help (CLI flags), squiz skill (full SKILL.md).`)
}

func printHelpAbout() {
	fmt.Fprintln(os.Stdout, strings.TrimSpace(`
squiz help — topical reference

Usage:
  squiz help                    list available topics
  squiz help <topic>            print the reference block for one topic
  squiz help --help             this message

Topics are stand-alone references — reading `+"`squiz help art`"+` should tell
you everything about art forms without needing to consult SKILL.md.

For the full SKILL.md (the canonical spec the binary was built from), run:

  squiz skill                   # dump SKILL.md to stdout
  squiz skill --out FILE        # write SKILL.md to FILE

For the CLI flag reference (render / example / version / etc.), run:

  squiz --help
`))
}

// topicOrder controls listing order — alphabetical-ish but with `art`,
// `themes`, `dsl` (the headline three) first so the most-asked topics land
// at the top of `squiz help`.
var topicOrder = map[string]int{
	"art":             1,
	"themes":          2,
	"dsl":             3,
	"schema":          4,
	"export":          5,
	"options":         6,
	"recommendations": 7,
}

var topicSummaries = map[string]string{
	"art":             `Art is composition from parts; the grammar, the visual budget, the checklist.`,
	"themes":          `The 8 themes; precedence; auto-rotation behavior.`,
	"dsl":             `The 11 parametric DSL primitives with grammar.`,
	"schema":          `High-level JSON shape (use "squiz schema" for machine-readable).`,
	"export":          `Shape of the JSON payload the user pastes back.`,
	"options":         `Per-question options: id, label, name, desc, art, recommendation.`,
	"recommendations": `When to set the recommendation field; one-per-squiz rule.`,
}

// topics holds the long-form reference for each `squiz help <name>` topic.
// Each value is a stand-alone reference (30-100 lines) — readers should not
// need to consult SKILL.md to use the topic correctly.
var topics = map[string]string{
	"art": `squiz help art — composition with library parts

Art is composition. The library and DSL primitives are parts.

The wf:*, arch:*, and DSL items are pieces you remix inside a custom
drawing — a labeled box of your own words, an arrow you draw between two
icons, a small composition that shows the specific idea the option is
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
the db edge is the hot path" — which is the option's actual desc, not
generic decoration.

Run "squiz catalog wf|arch|dsl" for the authoritative list of parts. The
--json flag reports each entry's naturalBox so you can size <use> boxes
without guessing.

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
    — only when the option is literally about a calendar grid

If you reach for one library token, finish the sentence "the picture I
want to draw IS this exact icon, alone, with no label." If you can't say
that, compose.

ANTI-PATTERNS THE VALIDATOR WATCHES FOR

  - Two options in a chooser share the same "art" — the icon means
    nothing if siblings collide (warning: sibling-art-collision)
  - Generic arch:user / arch:server / arch:database where the desc
    names something specific — give it a label or compose around it
  - A section where every option uses a single library token and no
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
  omitted or ""                         auto per-letter abstract pattern
                                        (A=hatched, B=dotted, …)

Raw SVG is no longer the "escape hatch" — it's the host of composition.
Use viewBox='0 0 100 60', style='width:80%;height:auto', and CSS vars
(var(--accent), var(--ink), var(--ink-2), var(--ink-3), var(--rule),
var(--rule-2)) so the picture inherits the active theme.

SEE ALSO

  squiz help dsl              full DSL grammar (18 primitives)
  squiz help options          how options carry art alongside other fields
  squiz help themes           how art inherits theme colors via CSS vars
  squiz catalog wf            list wf:* parts with naturalBox dimensions
  squiz catalog arch          list arch:* parts with naturalBox dimensions
  squiz catalog dsl           list DSL primitives with grammar
  squiz validate              flags sibling-art-collision, composition-thin,
                              default-art-fallback
`,

	"themes": `squiz help themes — the 8 themes & auto-rotation

squiz ships 8 retro-terminal × editorial themes. Every rendered doc lives in
exactly one theme; the binary picks it for you by default.

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
sequentially, persisted across runs. Each new repo you render in gets the
next theme in the rotation. Re-renders of the same repo always get the same
theme — so a project's identity stays stable across multiple squizzes.

The cache lives at ~/.squiz/themes.json. The repo key is:

  1. The output of "git remote get-url origin" if available, else
  2. The absolute working directory.

This makes a forked clone of the same upstream pick the same theme as the
original — usually the right thing.

PRECEDENCE (highest wins)

  1. --theme <name> on the CLI
  2. "theme" field in the JSON spec
  3. Auto-derived from repo cache

GUIDANCE

Omit "theme" from your JSON. Let it auto. Only set it explicitly when you
want a specific identity for one squiz — for example, "phosphor" for an
infra/code-flavored decision, "rose" for a brand/marketing one. Setting a
theme in the JSON does NOT update the repo cache; it's a one-doc override.

OPTIONAL TOGGLES

In addition to the theme, two layout dials sit alongside in the JSON:

  "scanlines": true            CRT-character scanline overlay; best with
                               phosphor / amber / slate.
  "density":  "comfortable"    Roomier line-heights and gutters. Default
                               is "compact".
  "cursor":   true             Blinking cursor in the squiz wordmark.

EXAMPLES

  { /* theme omitted */ }                              auto-rotate
  { "theme": "phosphor" }                              one-doc override
  squiz render foo.json --theme amber                  CLI override (wins
                                                       over JSON)

SEE ALSO

  squiz help art              themed art inherits CSS vars from the theme
  squiz help schema           where theme / density / scanlines live in JSON
`,

	"dsl": `squiz help dsl — the 11 parametric DSL primitives

The DSL is a family of compact strings the binary parses into themed SVG.
Use it when the wf:* / arch:* library doesn't ship the exact shape you need,
or when you want to compose a fresh idea from primitives.

All DSL output is theme-aware — strokes, fills, and text inherit CSS vars
(--ink, --ink-2, --accent, --rule, --rule-2) so the same string renders
differently across the 8 themes without authoring effort.

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

  Icons must be valid arch:* names (server, database, queue, …). See
  "squiz help art" for the full arch:* list, or run "squiz catalog".

ARROW: DIRECTIONS

  ?dir=right    default; left-to-right
  ?dir=down     top-to-bottom
  ?dir=up       bottom-to-top
  ?dir=left     right-to-left

JSON ESCAPING

Because the DSL strings live inside JSON, backslashes and quotes need
JSON-escaping. The text: form is the chief offender:

  // JSON-source                                  Effective DSL
  "art": "text:\"hi\\nthere\"@mono"            →  text:"hi\nthere"@mono
  "art": "sample:\"Quiet.\"@serif"             →  sample:"Quiet."@serif

If your JSON parser rejects a DSL string, count backslashes.

SEE ALSO

  squiz help art              when to pick DSL vs library vs raw SVG
  squiz help themes           how DSL output inherits theme CSS vars
  squiz catalog               every wf:* / arch:* name DSL can reference
`,

	"schema": `squiz help schema — high-level JSON shape

A squiz input file is a single JSON document with three top-level concerns:
theme/layout dials (optional), a "spec" block (optional narrative), and a
"squizzes" array (the questions). For the authoritative machine-readable
JSON Schema, run:

  squiz schema                 # write JSON Schema to stdout

For an end-to-end sample input you can copy and modify, run:

  squiz example                # writes squiz-example.json next to cwd

TOP-LEVEL SHAPE

  {
    // OPTIONAL — theme / layout dials. Omit "theme" for auto-rotation.
    "theme":     "paper",        // see "squiz help themes"
    "density":   "compact",      // compact | comfortable
    "scanlines": false,          // CRT overlay
    "cursor":    true,           // blinking cursor in the wordmark

    // OPTIONAL — narrative shown above the squizzes.
    "spec": {
      "path":  "/usr/specs/tide.md",          // shown in the topbar
      "title": "Tide — a habit tracker",      // page H1
      "lede":  "Eight decisions to lock in…", // one-line summary
      "paragraphs": [
        { "text": "Users land on a {{onboarding}} the very first time…" }
      ]
    },

    // REQUIRED — the questions themselves.
    "squizzes": [
      {
        "id":    "onboarding",                     // stable slug; in anchors
        "title": "First-launch experience",
        "desc":  "The first 60 seconds set expectations.",
        "quote": "the bar to 'doing the thing' should be embarrassingly low",
        "options": [
          {
            "id":    "jumpin",                     // stable slug
            "label": "Option A",                   // OPTIONAL; auto from index
            "name":  "Name one habit, go",         // short display
            "desc":  "Single text field. No setup.",
            "art":   "wf:phone-input"              // see "squiz help art"
          }
        ]
      }
    ]
  }

FIELD NOTES

  spec.paragraphs[].text — supports {{markers}} that link down to squiz
  cards by id. The renderer turns each marker into a clickable badge.

  squizzes[].quote — optional one-line pull-quote shown above the card.
  Use it when you can point to a specific spec line that motivates the
  question.

  squizzes[].options[].id — stable; comes back in the export payload as
  decisions[].choice.id. Pick short kebab-or-camel slugs that still
  make sense in a week.

  squizzes[].options[].label — optional. Auto-derived from option index
  (Option A, Option B, …) if omitted.

  squizzes[].options[].recommendation — see "squiz help recommendations".
  Renders a ★ RECOMMENDED chip + editorial callout. At most one per squiz.

VALIDATION

Run "squiz validate <input.json>" to check structure before rendering. The
validator catches: missing required fields, duplicate ids within a squiz,
unknown wf:/arch: names, malformed DSL strings, and {{markers}} that
don't resolve to a squiz id.

SEE ALSO

  squiz schema                 machine-readable JSON Schema
  squiz example                end-to-end sample input
  squiz help art               the "art" field in detail
  squiz help options           per-question option fields
  squiz help themes            theme / density / scanlines dials
  squiz help export            JSON the user pastes back after answering
`,

	"export": `squiz help export — the JSON the user pastes back

After the user fills in the rendered HTML, they click the sticky "copy json"
button at the bottom. That copies a single JSON document to the clipboard;
they paste it back into chat for you to parse.

SHAPE

  {
    "spec": "Tide — a habit-tracking app",
    "source": {
      "file":     "C:\\Users\\User\\code\\tide\\squiz.html",
      "basename": "squiz.html"
    },
    "generatedAt": "2026-05-25T12:34:56Z",
    "decisions": [
      {
        "id":       "onboarding",
        "question": "First-launch experience",
        "anchor":   "#squiz-onboarding",
        "choice": {
          "id":      "jumpin",
          "name":    "Name one habit, go",
          "summary": "Single text field. No setup. Customize later."
        },
        "notes": "Use a placeholder hint that rotates each launch."
      }
    ],
    "summary": { "total": 8, "resolved": 8, "withNotes": 3 }
  }

FIELDS

  spec               Echo of the rendered spec title.

  source.file        Absolute path of the rendered HTML on the user's
                     machine. Append an anchor to navigate to a specific
                     decision (e.g. file:///C:/.../squiz.html#squiz-onboarding).

  source.basename    Just the filename. Handy for log lines.

  generatedAt        ISO-8601 timestamp the user clicked "copy json".

  decisions[]        One entry per squiz, ordered as the doc renders them.

  decisions[].id           The stable squiz id you authored.
  decisions[].question     The squiz title (for human readability).
  decisions[].anchor       "#squiz-<id>" — append to source.file to link.
  decisions[].choice       The picked option, or null if the user skipped.
  decisions[].choice.id    Stable option id you authored.
  decisions[].choice.name  Short display name.
  decisions[].choice.summary  The option's desc (trimmed).
  decisions[].notes        Free-text notes the user left (string or null).

  summary.total       Number of decisions in the doc.
  summary.resolved    Number with a non-null choice.
  summary.withNotes   Number with a non-empty notes field.

INTERPRETING SKIPS

choice: null means the user skipped that decision. Treat it as "you decide"
unless the notes field carries an instruction. Don't re-render a second
squiz to re-ask — follow up in plain chat.

ANCHOR USAGE

Combine source.file + anchor to deep-link back to the exact card the user
was responding to:

  file:///C:/Users/.../squiz.html#squiz-onboarding

Useful when you want to quote the user's answer back to them with a
clickable reference.

SEE ALSO

  squiz help schema            input shape (what you author)
  squiz help options           per-option fields that round-trip into "choice"
  squiz help recommendations   how recommendations show up vs. influence choice
`,

	"options": `squiz help options — per-question option fields

Each squiz carries an "options" array. Each option is the smallest unit the
user picks between — a card with a name, a 1-2 sentence trade-off, and
optionally a piece of art.

OPTION FIELDS

  id                REQUIRED. Stable slug. Comes back in the export as
                    decisions[].choice.id — pick short kebab-or-camel
                    slugs that still make sense in a week.

  label             OPTIONAL. Display string (typically "Option A",
                    "Option B"). Auto-derived from the option's index if
                    omitted (A, B, C, …). Override when you want
                    semantic labels ("Cheapest", "Fastest").

  name              REQUIRED. Short headline shown big on the card.
                    Keep it punchy — 2-6 words is the sweet spot.

  desc              REQUIRED. The 1-2 sentence trade-off. This is where
                    the user actually decides. Lean on concrete numbers,
                    constraints, and the consequence of picking this
                    option.

  art               OPTIONAL. Mini-visual rendered into the card's art
                    slot. Five forms: wf:* / arch:* (library), DSL
                    primitive, raw <svg>, "none" (collapse the slot),
                    or omit for an auto per-letter abstract. See
                    "squiz help art" for the full picker.

  recommendation    OPTIONAL. One or two sentences explaining why this
                    option is the preferred pick. Renders a ★ RECOMMENDED
                    chip on the option and an editorial callout under
                    desc. At most one per squiz. See "squiz help
                    recommendations" for when to use it.

EXAMPLE

  {
    "id":             "sqlite",
    "label":          "Option A",
    "name":           "SQLite (single file)",
    "desc":           "Boring, durable, queryable. ~2 MB binary footprint.",
    "art":            "wf:file-icons",
    "recommendation": "NFR-2 requires durability without an external dep; SQLite hits both targets and the audience can debug it without learning anything new."
  }

GUIDANCE

Keep "name" terse and "desc" concrete. Vague descs ("flexible, scalable,
easy to use") give the user nothing to decide on. Specific descs ("$5/mo
ceiling, 100 ms p50, single-region") are what make squizzes work.

Use 2-5 options per squiz. One option is not a decision; six-plus is a
menu. If you have six candidates, pre-filter offline and present the top
three with the rationale.

Avoid asymmetric option sets where one option is clearly serious and the
others are throwaways. The user reads it as you padding the page. Either
make all options viable or drop the question.

SEE ALSO

  squiz help art               the "art" field across all 5 forms
  squiz help recommendations   when and how to set "recommendation"
  squiz help schema            where options live in the JSON shape
  squiz help export            how the picked option round-trips back
`,

	"recommendations": `squiz help recommendations — the "recommendation" field

Any option in a squiz can carry an optional "recommendation" field — a one-
or-two-sentence rationale for why this is the preferred choice given the
spec or constraints.

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

Set "recommendation" when the spec, constraints, audience, or prior
decisions genuinely point at one option. Concretely:

  - You can cite a specific ref or constraint in the rationale
    ("OVR-3's $5-VPS constraint", "NFR-2 requires durability").
  - The other options have real downsides that the recommended one
    avoids.
  - A reasonable reviewer would agree the recommendation is earned.

DON'T set it when:

  - You don't have a real preference. Marking every option as
    "recommended" defeats the purpose (the chip stops carrying signal).
  - The rationale is fluff ("flexible, modern, popular"). If you can't
    cite a specific constraint, the recommendation isn't earned.
  - You're hedging. A wishy-washy "either is fine" recommendation is
    worse than none — it nudges the user without committing.

THE ONE-PER-SQUIZ RULE

Set "recommendation" on AT MOST ONE option per squiz. Multiple
recommendations within the same question make the signal meaningless.

If two options genuinely look equal, that's a sign the question is
under-constrained — add a constraint to your overview, or restate the
question more sharply, before resorting to a recommendation.

CITING REFS / CONSTRAINTS

The strongest recommendations cite specific upstream commitments. In a
plain squiz (no plan), cite the spec narrative:

  "The spec's 'embarrassingly low setup' line rules out the multi-step
  onboarding; that leaves A vs C, and A's single field is a tighter fit."

In a squiz that follows a /squiz-plan run, cite the plan's IDs directly:

  "OVR-3 caps us at one VPS; NFR-2 requires durability without an
  external dep. SQLite hits both."

GUIDANCE FOR LENGTH

One or two sentences. Three is too many — the callout overshadows the
desc. If you need more justification, the recommendation isn't earned;
restructure the question.

SEE ALSO

  squiz help options           the field in the broader option shape
  squiz help schema            where options live in JSON
  squiz help export            recommendations don't show in export;
                               choice.id reflects what the user picked
`,
}
