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
	"art":             `The 5 art forms; when to pick wf:, arch:, DSL, raw SVG, "none".`,
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
	"art": `squiz help art — the 5 art forms

Every option in a squiz can carry an "art" field. The renderer detects which
of five forms you've used from the string's shape, then renders themed SVG
into the option's art slot.

AUTHORING PREFERENCE ORDER

    wf:* / arch:*   >   DSL primitive   >   "none"   >   raw <svg>

Reach for the named library first; it's the cheapest, ships theme-aware
strokes/fills out of the box, and renders identically across all 8 themes.
Drop to DSL when nothing in the library captures the shape. Use "none" when
the question is purely textual (a name, a string) and an art slot would just
be padding. Reach for raw SVG only when the option needs a bespoke metaphor
(a "living garden" plant, a custom diagram) that nothing else captures.

THE 5 FORMS

1. NAMED LIBRARY  — "art": "wf:<name>" or "art": "arch:<name>"

   The binary ships ~50 curated wireframes (wf:*) and ~30 system-design
   icons (arch:*) baked in. Pick one by name; the renderer themes it via
   CSS vars so it inherits whatever theme is active.

   wf:* names cover: calendars (calendar-grid, day-strip, year-heatmap),
   charts (spark-rising, bars-up, donut, gauge), avatars (avatar-single,
   avatar-feed), phone screens (phone-input, phone-card, phone-tabs),
   controls (toggle-on, slider, dropdown), status (badge-new, lock,
   check-large), typography (serif-sample, sans-sample, mono-sample),
   graphs (graph-force, tree-hier, matrix-heatmap), metaphors (plant-grow,
   garden, paper-fold), and misc (cmd-palette, text-cursor, file-icons).

   arch:* names cover system-architecture primitives: server, database,
   cache, queue, load-balancer, cdn, gateway, api, worker, function,
   scheduler, user, browser, mobile, firewall, storage, blob, table,
   stream, log, metric, trace, container, pod, vpc, subnet, dns, secret,
   key-icon, topic.

   Pick the namespace that matches the kind of picture you're making:
   arch:* for system-design diagrams, wf:* for UI/UX wireframes.

   Run "squiz catalog" for the complete authoritative list with previews.

2. PARAMETRIC DSL  — "art": "<dsl-string>"

   Compact strings the binary parses into themed SVG. 11 primitives:
   grid:, spark:, bars:, swatches:, pills:, sample:, circle-pack:, text:,
   flow:, box:, arrow:. Run "squiz help dsl" for the full grammar.

   DSL is composable — flow: can embed arch:* icons via ?icon=, box: can
   carry an arch icon, text: can mix font/size/color/weight. Reach for DSL
   when you need to compose an idea the library doesn't ship.

3. RAW SVG  — "art": "<svg ...>...</svg>"

   When the library and DSL don't fit, inline raw SVG starting with "<svg".
   Use CSS vars (var(--accent), var(--ink), var(--ink-3), var(--rule-2))
   so the shape inherits the active theme. Use viewBox='0 0 100 60' and
   style='width:80%;height:auto' to match the visual weight of the other
   forms.

4. EXPLICIT HIDE  — "art": "none"

   Drops the art slot entirely; the card collapses to text-only. Use when
   no visual is appropriate (a name/string/free-text question). Better
   than forcing irrelevant art.

5. AUTO PER-LETTER ABSTRACT  — "art" omitted or empty string

   Subtle patterns based on option position: A = hatched, B = dotted,
   C = striped, D = grid, E = cross-hatch, F = waves. Looks intentional
   without authoring effort. Use when you're moving fast and visuals
   don't materially help the decision.

EXAMPLES

  { "art": "wf:phone-input" }                  named library, UI wireframe
  { "art": "arch:database" }                   named library, system icon
  { "art": "spark:[3,5,4,7,6,9,11]" }          DSL sparkline
  { "art": "flow:[client?icon=user,api,db]" }  DSL pipeline with icons
  { "art": "none" }                            no art for this option
  // omit "art" entirely                       auto per-letter pattern

SEE ALSO

  squiz help dsl              full DSL grammar
  squiz help options          how options carry art alongside other fields
  squiz help themes           how art inherits theme colors
  squiz catalog               browse all wf:* / arch:* names with previews
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
