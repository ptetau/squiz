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

### Phase 3 — User reviews + pastes back

The doc has six top tabs, badge cross-refs, a per-item **feedback widget** (✓ approve / ? question / ✗ reject + notes + optional inline edits), and a sticky `copy json` button.

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

The `art` field uses the same five forms — see `/squiz` SKILL.md for the full reference:

1. `"art": "wf:<name>"` — 50 named wireframes (calendar-grid, spark-rising, phone-card, …)
2. `"art": "<dsl-string>"` — 7 parametric primitives (`grid:`, `spark:`, `bars:`, `swatches:`, `pills:`, `sample:`, `circle-pack:`)
3. `"art": "<raw svg>"` — escape hatch (use CSS vars for theme inheritance)
4. `"art": "none"` — explicitly hide the art slot for this item
5. `art` omitted — no art shown (plan items don't get a per-letter auto-pattern like squiz options; they just have no art)

Reach for `wf:` / DSL first; raw SVG only for bespoke metaphors. Reach for `"none"` instead of forcing irrelevant art.

## Cross-references (`refs`)

Each item carries optional `refs: ["OVR-1", "FR-3"]`. The renderer:
- Validates each ref ID exists somewhere in the plan
- Renders each ref as an inline badge with the section label: `[Functional · FR-3]`
- Click switches to the target tab + scrolls to + highlights the target item
- Browser back-button returns to the previous item

**Convention**: refs point "upward" in the spine — `build` items ref `engineering`, which refs `functional`/`non-functional`, which refs `overview`. You CAN sideways-ref (an `ENG-` item refs another `ENG-` item) but use it sparingly; it makes the plan feel knotty.

## Theme

Same auto-rotation as squiz. Omit `theme` from `index.json` unless overriding. Each repo gets a distinct theme on first render, persisted in `~/.squiz/themes.json`. CLI `--theme <name>` trumps both.

The 8 themes: `paper` / `phosphor` / `amber` / `beige` / `rose` / `ocean` / `forest` / `slate`.

## Export JSON shape (what the user pastes back)

```json
{
  "plan": "ThermoLog — home temperature logger",
  "source": {
    "file":     "C:\\dev\\thermolog\\plan\\index.html",
    "basename": "index.html"
  },
  "generatedAt": "2026-05-25T12:34:56Z",
  "feedback": [
    {
      "id":     "FR-3",
      "status": "questioned",
      "anchor": "#item-FR-3",
      "note":   "20 minutes feels long — most pipe-freezing scenarios are faster than that.",
      "edits":  null
    },
    {
      "id":     "BUILD-2",
      "status": "approved",
      "anchor": "#item-BUILD-2",
      "note":   null,
      "edits":  { "title": "Pi firmware (Go single binary)" }
    }
  ],
  "summary": { "total": 21, "approved": 18, "questioned": 2, "rejected": 0, "withNotes": 4, "withEdits": 1 }
}
```

Items the user didn't touch don't appear in `feedback`. `status` is one of `"approved"`, `"questioned"`, `"rejected"`. `edits` is a sparse object of field overrides — apply them as suggestions, not authoritative truth.

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

## Files in this skill

- `SKILL.md` — this file.
- The `squiz-plan` binary, installed alongside `squiz` via the same install scripts.
