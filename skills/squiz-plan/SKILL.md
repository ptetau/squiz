---
name: squiz-plan
description: >
  Visual structured-plan renderer with the Apple //e × IBM Plex aesthetic.
  Sibling to /squiz. The agent splits a plan or spec into 6 canonical
  sections (overview → functional → non-functional → cases → engineering →
  build) plus optional custom sections, writes them as a multi-file JSON
  tree, and the `squiz-plan` CLI renders a tabbed interactive HTML doc
  where every item carries clickable badges back to the parents that
  motivated it. The user reads, approves/questions/rejects per item, leaves
  notes, and pastes a JSON payload back. Use when the user asks for a
  "plan", "spec", or says "/squiz-plan" — and when traceability between
  layers matters (audit, code review, design doc). Prefer /squiz for
  clarification flows.
---

# squiz-plan: structured plans, visible threads

## Purpose

Most plans collapse under their own weight: by the time you read step 17 in the build section, you've forgotten which functional requirement motivated it. `squiz-plan` solves that by making the thread between layers **visible and clickable**. Every item carries `refs:` to the parent items that motivated it; the rendered HTML turns those refs into inline badges that, when clicked, switch tabs and highlight the parent.

The 6 canonical sections are the spine:

1. **Overview** — mission, audience, constraints
2. **Functional requirements** — what the system does
3. **Non-functional requirements** — how it behaves (perf, security, offline, etc.)
4. **Cases** — real-world scenarios the system has to handle
5. **Engineering requirements** — architecture decisions, components
6. **Build** — concrete steps, per component, to actually deliver

Agents may append custom sections (`glossary`, `risks`, `appendix`, etc.) after these six.

## When to use vs. siblings

- **`/squiz`** — clarifying questions BEFORE you know what to build
- **`/squiz-plan`** (this skill) — structured plan AFTER you know enough to commit
- **`/quiz`** — quick inline back-and-forth in chat, no document

A natural pairing: run `/squiz` to gather requirements, then `/squiz-plan` to turn the resolved decisions into a structured plan with traceability.

## The flow (3 phases, 1 binary)

### Phase 1 — Structure

Take the conversation/spec/prior squiz output and split it into the 6 sections. Aim for:
- **3-5 overview items** (mission, audience, hard constraints, success criteria)
- **3-8 functional reqs** (what it does)
- **2-5 non-functional reqs** (how it behaves)
- **2-5 cases** (concrete scenarios — these make the plan feel real)
- **3-8 engineering reqs** (the architecture — refs the FRs and NFRs)
- **3-10 build steps** (the work — refs the engineering reqs)

Each item gets a **stable ID** with the section's prefix:

| Section | Prefix | Example |
|---|---|---|
| overview | `OVR` | `OVR-1`, `OVR-2` |
| functional | `FR` | `FR-1`, `FR-2.3` |
| non-functional | `NFR` | `NFR-1` |
| cases | `CASE` | `CASE-1`, `CASE-bedroom-freeze` |
| engineering | `ENG` | `ENG-1`, `ENG-storage` |
| build | `BUILD` | `BUILD-1`, `BUILD-firmware` |

IDs are stable slugs — they come back in the user's feedback payload and they're what `refs:` arrays point at. Pick something that'll still make sense in a month.

### Phase 2 — Write the JSON tree, run the binary

Layout (one directory, multiple files):

```
plan/
├── index.json              ← top-level descriptor (mandatory)
├── overview.json           ← optional but conventional
├── functional.json
├── non-functional.json
├── cases.json
├── engineering.json
└── build.json
```

Then:

```bash
squiz-plan plan/index.json
```

The binary loads `index.json`, walks its `sections` list, loads each `<sectionId>.json` sibling, validates that every `refs:` ID actually exists in the plan, and emits `plan/index.html` next to `index.json`.

**Always hand the user a clickable `file://` URL** when telling them the path. `file:///C:/Users/.../plan/index.html` is clickable in modern terminals; the bare path is not.

### Phase 3 — User reviews, leaves notes at three scopes, proposes items, pastes back

The doc has six top tabs, badge cross-refs, a per-item **feedback widget** (✓ approve / ? question / ✗ reject + notes + optional inline edits), and a sticky `copy json` button.

The rendered HTML now exposes **three** places the user can leave notes:

1. **Per-item notes** — the textarea in each item's feedback widget (existing).
2. **Per-section notes** — sticky textarea at the top of each tab. For feedback that spans multiple items in a section (e.g. "these FRs are missing the export workflow").
3. **Plan-level notes** — single textarea inside the copy-json modal. For overall direction ("this scope is too ambitious for v1").

Each section's tab also has a **"+ add item"** button that lets the user *propose* brand-new items. Proposed items come back in the export as a typed `proposed_items[]` array — apply them as suggestions when regenerating.

> "Click `copy json` at the bottom of the plan, paste it back here, and I'll revise based on your feedback."

When they paste, parse the payload (shape below), apply edits and questions, regenerate the affected section files, re-render, and hand back the new URL. Resolved feedback removes the corresponding entries; unresolved items keep their feedback so the user knows what's pending.

## JSON Schema

### `plan/index.json` — the top-level descriptor

```jsonc
{
  "title":    "ThermoLog — home temperature logger",  // page H1
  "lede":     "Six load-bearing sections, traceable…", // one-line summary
  "theme":    "paper",     // OPTIONAL; same precedence as squiz (omit for auto)
  "density":  "compact",   // OPTIONAL; compact | comfortable
  "scanlines": false,      // OPTIONAL; CRT overlay
  "cursor":    true,       // OPTIONAL; blinking cursor in the wordmark
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
```

Canonical sections always render in canonical order regardless of position in `sections`. Custom sections (any ID not in the canonical six) append in the order declared.

### `plan/<sectionId>.json` — the per-section payload

```jsonc
{
  "items": [
    {
      "id":    "ENG-1",                              // stable; must start with section prefix
      "title": "Sensor driver (BLE)",
      "desc":  "Read from 4 cheap Bluetooth LE thermometers. Pi scans every 60 s.",
      "art":   "wf:dot-trend",                       // OPTIONAL; same forms as squiz
      "refs":  ["FR-1", "NFR-1"]                     // OPTIONAL; parent IDs (validated)
    }
  ]
}
```

Validator rejects: missing section file, wrong prefix, duplicate IDs, refs to nonexistent items.

## Art forms (same as squiz)

The `art` field uses the same forms — see `/squiz` SKILL.md for the full reference:

1. `"art": "wf:<name>"` — 50 named UI wireframes (calendar-grid, spark-rising, phone-card, …)
2. `"art": "arch:<name>"` — **NEW in v0.4.0** — ~30 system-design icons (`arch:server`, `arch:database`, `arch:queue`, `arch:load-balancer`, …). Distinct namespace from `wf:*`; pick `arch:*` for architecture diagrams and `wf:*` for UI sketches.
3. `"art": "<dsl-string>"` — parametric primitives. The originals (`grid:`, `spark:`, `bars:`, `swatches:`, `pills:`, `sample:`, `circle-pack:`) plus **NEW in v0.4.0**: `text:` (rich multi-line styled text), `flow:` (left-to-right pipeline of named boxes, optionally embedding `arch:*` icons), `box:` (single labeled box, optional icon), `arrow:` (standalone labeled arrow). See `/squiz` SKILL.md for full grammar.
4. `"art": "<raw svg>"` — escape hatch (use CSS vars for theme inheritance)
5. `"art": "none"` — explicitly hide the art slot for this item
6. `art` omitted — no art shown (plan items don't get a per-letter auto-pattern like squiz options; they just have no art)

Reach for `wf:` / `arch:` / DSL first; raw SVG only for bespoke metaphors. Reach for `"none"` instead of forcing irrelevant art.

## Cross-references (`refs`)

Each item carries optional `refs: ["OVR-1", "FR-3"]`. The renderer:
- Validates each ref ID exists somewhere in the plan
- Renders each ref as an inline badge with the section label: `[Functional · FR-3]`
- Click switches to the target tab + scrolls to + highlights the target item
- Browser back-button returns to the previous item

**Convention**: refs point "upward" in the spine — `build` items ref `engineering`, which refs `functional`/`non-functional`, which refs `overview`. You CAN sideways-ref (an `ENG-` item refs another `ENG-` item) but use it sparingly; it makes the plan feel knotty.

## Item options (decisions)

Any item can carry an `options: [...]` field with the same shape as squiz options (`id`, `label?`, `name`, `desc`, `art?`, `recommendation?`). When present, the rendered card shows a chooser; without it, the card is a flat statement. Use options when the item represents an *unsettled* decision (which database? which deploy strategy?) — leave them off for settled facts. The user's pick comes back in the export's `feedback[].chose` field.

Any option can carry an optional **`recommendation`** field — a one-or-two-sentence rationale for why this is the preferred choice given the spec. The renderer shows a `★ RECOMMENDED` chip next to the option's name and the explanation as an editorial callout under its desc. The user can still pick differently — recommendations are advisory, not constraints.

```jsonc
{
  "items": [
    {
      "id":    "ENG-2",
      "title": "Local storage engine",
      "desc":  "Pick the embedded store the Pi writes readings into.",
      "refs":  ["FR-1", "NFR-2"],
      "options": [
        {
          "id":   "sqlite",
          "name": "SQLite (single file)",
          "desc": "Boring, durable, queryable. ~2 MB binary footprint.",
          "art":  "wf:file-icons",
          "recommendation": "NFR-2 requires durability without an external dep; OVR-3 caps us at one VPS. SQLite hits both targets and is the option the audience can debug without learning anything new."
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
  ]
}
```

**Recommendation guidance:** at most one option per item. Use it when the spec/constraints/audience genuinely point at one choice. Don't recommend every option (defeats the purpose); don't fluff the explanation (the rationale should cite specific refs like `OVR-3` or `NFR-2` whenever possible).

**Validator note:** option `id`s must be unique *within* an item. Collisions across different items are fine (two items can both have an `id: "sqlite"` option) — the export disambiguates via the parent item's ID.

## Theme

Same auto-rotation as squiz. Omit `theme` from `index.json` unless overriding. Each repo gets a distinct theme on first render, persisted in `~/.squiz/themes.json`. CLI `--theme <name>` trumps both.

The 8 themes: `paper` / `phosphor` / `amber` / `beige` / `rose` / `ocean` / `forest` / `slate`.

## Export JSON shape (what the user pastes back)

```jsonc
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
      "chose":  "sqlite"                 // NEW: which option the user picked (null when no options)
    },
    {
      "id":     "FR-3",
      "status": "questioned",
      "anchor": "#item-FR-3",
      "note":   "20 minutes feels long — most pipe-freezing scenarios are faster than that.",
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
  "section_notes": {                     // NEW: keyed by section ID
    "functional":  "We're missing the rate-limit requirement.",
    "engineering": "ENG-3 could fold into ENG-2."
  },
  "plan_note": "Overall: too ambitious for v1, see section notes.",   // NEW
  "proposed_items": [                    // NEW: user-suggested additions
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
```

Items the user didn't touch don't appear in `feedback`. `status` is one of `"approved"`, `"questioned"`, `"rejected"`. `edits` is a sparse object of field overrides — apply them as suggestions, not authoritative truth. `chose` is the picked option `id` for items that carried an `options:` array (and `null` otherwise). `section_notes`, `plan_note`, and `proposed_items` are all NEW in v0.4.0 — see Rule 10.

## CLI reference

```bash
squiz-plan <plan/index.json>                # render index.html next to index.json
squiz-plan render <plan/index.json>         # explicit subcommand
squiz-plan <plan/index.json> --open         # also open in default browser
squiz-plan <plan/index.json> --theme phosphor
squiz-plan <plan/index.json> --out path.html
squiz-plan <plan/index.json> --stdout > x.html
squiz-plan version
squiz-plan help
```

Flags may appear before OR after the positional `index.json` (unlike Go's stdlib `flag.Parse` default).

## Accessibility (built into the renderer)

The rendered HTML ships with: a skip-to-tabs link, `tablist`/`tab`/`tabpanel` ARIA roles with arrow-key navigation, visible focus rings on every interactive, modal focus trap with return-focus on close, `aria-live` progress announcements on tab switch, proper `<label>` associations for notes textareas, `prefers-reduced-motion` support. No additional work for the agent.

## Rules

1. **One plan per invocation.** If feedback reveals deep ambiguity, switch back to `/squiz` or chat — don't render a second plan immediately.
2. **Refs upward only by default.** `build` → `engineering` → `func/non-func` → `overview`. Sideways refs are fine occasionally; never make refs into a spaghetti.
3. **Stable IDs with section prefix.** `FR-1`, not `req-1`. The prefix is what the validator checks against the section name.
4. **Omit `theme`** unless overriding. Auto-rotation does the right thing.
5. **Make `art` earn its slot.** Reach for `wf:` first, DSL second, raw SVG only when nothing else captures the idea. `"none"` is fine.
6. **Cases sell the plan.** A plan with no `cases.json` section feels abstract. Even 2-3 short cases make the rest feel concrete.
7. **Self-contained.** Title + lede must let a reader understand the plan in 10 seconds. The lede is the elevator pitch.
8. **Clickable links.** When you hand the user the rendered file, format as `file:///...` — bare paths aren't clickable in most terminals.
9. **Apply feedback as a follow-up.** When the user pastes back, restate what you understood, regenerate the affected section files, re-render, and hand back the new clickable URL.
10. **Round-trip the notes.** When the user pastes feedback back, treat `note` / `section_notes` / `plan_note` / `proposed_items` as the agent's instructions for the next round: rewrite affected items, restructure sections, or append the proposed items as fresh entries. Apply edits as suggestions, not authoritative changes (you may push back if they break the plan's spine).
11. **Recommend when you have a real preference.** Any item-option can carry `"recommendation": "<one or two sentences>"`. Use it only when the plan's overview/refs genuinely point at one option — cite refs in the explanation (`"OVR-3's $5-VPS constraint rules out k8s; that leaves systemd vs Docker; systemd is one fewer moving part"`). At most one per item. Don't fluff; if you can't justify in one sentence, don't recommend.
12. **Clarify at every level — `options[]` aren't just for engineering.** A plan exists *because* there's still ambiguity to resolve; if every overview/functional/case item is a flat assertion, the plan is pretending it has more certainty than it does. Audit each section for ambiguities (more than one defensible reading), contradictions (two items pull opposite directions), and redundancies (two items say the same thing in different words) — surface each as an `options:` chooser on the relevant item so the user resolves it the same way they'd resolve an engineering decision. Examples:
    - **Overview** — *"Mission"* with options *"v1 (shippable in 2 weeks)"* vs *"v1.5 (with analytics)"* — the choice scopes everything downstream.
    - **Functional** — *"Alert delivery"* with options *"real-time push"* vs *"daily digest email"* vs *"both, user-selectable"*.
    - **Non-functional** — *"Backup cadence"* with options *"continuous WAL"* vs *"hourly snapshot"* vs *"daily snapshot only"*.
    - **Cases** — *"Primary scenario when prioritising"* with options *"new-user onboarding"* vs *"power-user retention"* vs *"recovery from failure"*.
    - **Build** — *"Ship sequence"* with options *"backend first, then UI"* vs *"vertical slice (one user flow end-to-end)"* vs *"infra → backend → UI"*.

    Use the recommendation field on each (Rule 11) when the spec genuinely points at one. Items that aren't decisions stay as flat statements — don't fabricate ambiguity to look thorough.

13. **Every item carries visual weight by default.** When you omit `art` on an item in a canonical section, the renderer fills in a section-appropriate default so every card has an anchor — not just engineering and build. The defaults:

    | Section | Default art (when `art` is omitted) |
    |---|---|
    | overview | `wf:avatar-single` |
    | functional | `wf:phone-blank` |
    | non-functional | `wf:gauge` |
    | cases | `wf:phone-card` |
    | engineering | `arch:server` |
    | build | `wf:cmd-palette` |

    Override per item when something more specific fits (a `flow:` for a case that's really a pipeline; an `arch:queue` for an ENG item about messaging). Use `"art": "none"` to explicitly suppress for items that genuinely shouldn't have one. Custom sections get no default (omitted == no art).

## Files in this skill

- `SKILL.md` — this file.
- The `squiz-plan` binary, installed alongside `squiz` via the same install scripts.
