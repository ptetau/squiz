---
name: squiz
description: >
  Visual document-style clarifier rendered by the Apple //e × IBM Plex "Squiz"
  Go binary. The agent writes a compact JSON spec (squizzes, options, optional
  spec narrative); the `squiz` CLI renders a self-contained, retro-styled
  interactive HTML doc with mini-wireframe option art and a sticky `copy json`
  status bar. The user fills it in at their own pace and pastes a single JSON
  payload back. The payload carries the rendered file path + per-decision
  anchors so the agent can navigate back to the exact section the user was
  responding to. Theme is auto-assigned per repo (sequential, persisted) so
  each project gets a distinct identity. Use whenever the user says "/squiz",
  "squiz this", or asks for a richer, more designed-feeling visual version of
  /clarify-with-docs. The visual twin of /quiz (inline chat). Prefer it when
  the task is high-stakes or visual decisions matter.
---

# Squiz: JSON-in, HTML-out, Apple //e × IBM Plex

## Purpose

Same intent as [[clarify-with-docs]] — gather every clarifying question up front and render them as one rich interactive document the user fills in at their own pace. Squiz is the **opinionated visual implementation**: numbered "squiz cards", optional inline spec text with `{{markers}}` linking to questions, mini-wireframe option art, and a sticky bottom status bar that exports as JSON. Eight retro-terminal × editorial themes (paper / phosphor / amber / beige / rose / ocean / forest / slate), auto-assigned per repo so every project feels distinct.

Use it for high-stakes clarification where the visual identity *is* part of the message: a UX/UI brief, a product spec, a design system call, anything where the user benefits from seeing wireframe previews of their options.

## When to use this vs. siblings

- **`/quiz`** (inline chat twin) — quick rounds, one question per turn in chat. Use when the user wants speed or is on mobile.
- **`/clarify-with-docs`** — same document-style pattern, but unstyled. Use when you don't want the retro flavor.
- **`/squiz`** (this skill) — when the visual identity matters, when options benefit from mini-wireframe previews, or when the user explicitly asks for "squiz".

## The flow (3 phases, 1 binary)

### Phase 1 — Gather

Review the conversation, memory, uploaded files. Identify every ambiguity that could materially change the output. Aim for **3-12 questions**. Fewer → use `/quiz`. More → the user will bail.

Decide whether to include a **spec narrative** at the top (paragraphs with `{{markers}}` linking down to cards) — only if you have real prose you can quote. Skip it for tighter squizzes.

### Phase 2 — Write JSON, run the binary

Write a `<name>.json` file (see Schema below), then invoke the renderer:

```bash
squiz <name>.json
```

That's it. The `squiz` binary (installed via `go install` of `<this-skill>/go-renderer`) reads the JSON, picks a theme for the repo (auto, sequential, persisted), and emits `<name>.html` next to the JSON. With no extra flag it does **not** auto-open.

**Always hand the user a clickable `file://` URL** when you tell them the path. Bare Windows paths like `C:\Users\…\foo.html` aren't clickable in most terminals; `file:///C:/Users/.../foo.html` is. POSIX form: `file:///home/u/foo.html`. Use `--open` only when the user has asked for automatic browser launch — and even then, still print the URL alongside.

**Default file location:** write `squiz.json` (or a named variant like `squiz-onboarding.json`) **next to the work the squiz is about** — typically the project root or the relevant subdirectory. Output `.html` lands next to the `.json` with the same basename. Both paths are deterministic.

### Phase 3 — User pastes back

The doc has a sticky **`copy json`** button at the bottom. Tell the user:

> "Click `copy json` at the bottom, paste it back here, and I'll continue."

When they paste, parse it (see "Export JSON shape" below) — the payload includes `source.file` (absolute path to the rendered HTML) and per-decision `anchor` (`#squiz-<id>`). Use these to navigate back to context if needed. Restate your updated understanding in plain prose, ask any *new* gaps as plain chat follow-ups (do NOT render a second squiz), then start the work.

## JSON Schema (input)

```jsonc
{
  // OPTIONAL — leave omitted for auto-rotation per repo. Set explicitly to override.
  "theme":     "paper",        // paper|phosphor|amber|beige|rose|ocean|forest|slate
  "density":   "compact",      // compact|comfortable  (default compact)
  "scanlines": false,          // CRT scanline overlay
  "cursor":    true,           // blinking cursor in the squiz wordmark

  "spec": {
    "path":  "/usr/specs/tide.md",          // shown in the topbar
    "title": "Tide — a habit tracker",      // page H1
    "lede":  "Eight decisions to lock in…", // one-line summary
    "paragraphs": [                          // OPTIONAL spec narrative
      { "text": "Users land on a {{onboarding}} the very first time…" }
    ]
  },

  "squizzes": [
    {
      "id":    "onboarding",                            // stable slug; appears in anchors
      "title": "First-launch experience",
      "desc":  "The first 60 seconds set expectations.",
      "quote": "the bar to 'doing the thing' should be embarrassingly low",
      "options": [
        {
          "id":    "jumpin",                              // stable slug
          "label": "Option A",                            // OPTIONAL — auto from index
          "name":  "Name one habit, go",                  // short display
          "desc":  "Single text field. No setup.",        // 1-2 sentence trade-off
          "art":   "wf:phone-input"                       // see Art forms below
        }
      ]
    }
  ]
}
```

## Art forms (the `art` field on each Option)

Five shapes, detected from the string content:

### 1. `"art": "wf:<name>"` — named library (preferred)

The binary ships ~50 curated wireframes baked in. Pick one by name. Theme-aware via CSS vars.

**Categories & names** (50 total):

- **calendars/dates** — `calendar-grid`, `calendar-week`, `streak-counter`, `day-strip`, `year-heatmap`, `time-of-day`, `clock`
- **charts** — `spark-rising`, `spark-flat`, `spark-noisy`, `bars-up`, `donut`, `gauge`, `dot-trend`
- **identities/avatars** — `avatar-single`, `avatar-pair`, `avatar-circle`, `avatar-feed`, `avatar-private`
- **phone screens** — `phone-blank`, `phone-list`, `phone-card`, `phone-input`, `phone-tabs`, `phone-onboard`, `phone-stats`
- **controls** — `toggle-on`, `toggle-off`, `button-accent`, `button-ghost`, `slider`, `dropdown`
- **status** — `badge-new`, `pill-row`, `snowflake`, `lock`, `check-large`
- **typography** — `serif-sample`, `sans-sample`, `mono-sample`
- **connections/graphs** — `graph-force`, `tree-hier`, `radial-burst`, `matrix-heatmap`
- **metaphors** — `plant-grow`, `garden`, `paper-fold`
- **misc** — `cmd-palette`, `text-cursor`, `file-icons`
- **system-design icons (arch:* namespace)** — server, database, cache, queue, load-balancer, cdn, gateway, api, worker, function, scheduler, user, browser, mobile, firewall, storage, blob, table, stream, log, metric, trace, container, pod, vpc, subnet, dns, secret, key-icon, topic (30 total)

> **Note:** `arch:*` and `wf:*` are distinct namespaces that coexist. `arch:*` icons are for system-design / architecture diagrams (servers, queues, databases); `wf:*` is for UI wireframes (phone screens, controls, charts). Pick the namespace that matches the *kind* of picture you're making.

### 2. `"art": "<dsl-string>"` — parametric DSL

Compact strings the binary parses into themed SVG. Primitives:

| Form | Example | Renders |
|---|---|---|
| `grid:NxM[@RATE]` | `"grid:7x7@0.55"` | N×M heatmap, RATE in [0,1] |
| `spark:[V,V,V,…]` | `"spark:[3,5,4,7,6,9,11]"` | sparkline from data |
| `bars:[V,V,V,…]` | `"bars:[3,5,4,7,6,9,11]"` | bar chart |
| `swatches:#A,#B,…` | `"swatches:#f1ebde,#1a1814,#b34a1a"` | palette swatches |
| `pills:A*\|B\|C*` | `"pills:morning*\|midday\|evening*"` | chip row, `*` = active |
| `sample:"text"[@FONT]` | `"sample:\"Quiet welcome back.\"@serif"` | styled sample text, FONT = `serif`/`sans`/`mono` |
| `circle-pack:N` | `"circle-pack:12"` | N organically-arranged circles |
| `text:"line 1\nline 2"[@FONT][?size=N&align=A&weight=W&color=C]` | `"text:\"Quiet\\nwelcome back.\"@mono?size=18&align=center&weight=700&color=accent"` | multi-line styled text (richer sibling of `sample:`). FONT = `mono`/`serif`/`sans` (default `sans`). `size` 6-36 (default 14). `align` = `left`/`center`/`right` (default `left`). `weight` 300-700 (default 400). `color` = `ink`/`ink-2`/`ink-3`/`accent`/`rule`/`rule-2` (default `ink`). Multi-line via `\n`. |
| `flow:[a,b,c]` or `flow:[a?icon=user,b?icon=api,c?icon=database]` | `"flow:[client?icon=user,api?icon=api,db?icon=database]"` | left-to-right pipeline of named boxes connected by arrows; optional `?icon=<arch>` embeds an arch icon in each box |
| `box:label[?icon=ARCH]` | `"box:web-tier?icon=server"` | single labeled box with optional arch icon |
| `arrow:"label"[?dir=DIR]` | `"arrow:\"async\"?dir=down"` | standalone labeled arrow glyph; `dir` = `right` (default) / `down` / `up` / `left` |

### 3. `"art": "<raw svg>"` — escape hatch

When library + DSL don't fit, inline raw SVG starting with `<svg`. Use CSS vars (`var(--accent)`, `var(--ink)`, `var(--ink-3)`, `var(--rule-2)`) so it inherits the active theme. Use `viewBox='0 0 100 60'` and `style='width:80%;height:auto'` to match the other forms visually.

### 4. `"art": "none"` — explicit hide

Drops the art slot entirely. Card collapses. Use when **no visual is appropriate** (e.g. a name/string question where art would be padding).

### 5. `art` omitted / empty — auto per-letter abstract

Subtle patterns based on option position: A = hatched, B = dotted, C = striped, D = grid, E = cross-hatch, F = waves. Looks intentional without authoring. Use as the default when you're moving fast and the visuals don't matter.

**Authoring order of preference:** `wf:` / `arch:` > DSL > `"none"` > raw SVG. Use `arch:*` for any system-architecture diagram (servers, queues, data stores, network topology); `wf:*` for UI/UX wireframes. Reach for raw SVG only when the option needs a bespoke metaphor (a "living garden" plant, a custom diagram) that nothing else captures. Reach for `"none"` instead of forcing art that doesn't help.

## Theme (auto by default — don't set unless overriding)

**Default behavior:** the binary auto-assigns one of 8 themes per repo, sequentially, persisted in `~/.squiz/themes.json`. Each new repo gets the next theme in rotation. Re-renders of the same repo always get the same theme. Repo key = `git remote get-url origin` if available, else the absolute working directory.

**Precedence:**
1. `--theme <name>` CLI flag (highest)
2. `"theme"` field in JSON
3. Auto-derived from repo cache

**Generally**: omit `theme` from the JSON. Let it auto. Only set it if you want a specific identity for one squiz (e.g. `phosphor` for a code/infra-flavored decision).

The 8 themes:

| Theme | Vibe | Mode |
|---|---|---|
| `paper` | Cream + ink + rust accent. Editorial, calm. | light |
| `phosphor` | Green-on-black CRT. | dark |
| `amber` | IBM 3279 amber on near-black. | dark |
| `beige` | PS/2 cream with IBM blue. | light |
| `rose` | Warm pink, plum ink, rose accent. | light |
| `ocean` | Pale blue-grey, deep teal, coral accent. | light |
| `forest` | Oat cream, moss, warm gold. | light |
| `slate` | Cool dark grey, electric blue accent. | dark |

Optional toggles: `data-scanlines="on"` for CRT character (best with `phosphor`/`amber`/`slate`), `data-density="comfortable"` for roomier layout.

## Export JSON shape (what the user pastes back)

```json
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
```

`source.file` is the absolute path of the rendered HTML; `anchor` is the `#squiz-<id>` you can append to it to navigate to a specific decision. `choice: null` means the user skipped that decision — treat as "you decide" unless `notes` say otherwise.

## CLI reference

```bash
squiz <input.json>                    # render <input>.html next to input
squiz render <input.json>             # same, explicit subcommand
squiz <input.json> --open             # also open in default browser
squiz <input.json> --theme phosphor   # force a specific theme
squiz <input.json> --out path.html    # explicit output path
squiz <input.json> --stdout > x.html  # write to stdout
squiz version
squiz help
```

## Installation (one-time)

```bash
cd ~/.claude/skills/squiz/go-renderer
go install .
```

Puts `squiz` on `$GOPATH/bin` (usually on `$PATH`). Verify with `squiz version`.

## Accessibility (built into the renderer)

The rendered HTML ships with: a skip-to-decisions link, `radiogroup` ARIA roles per squiz, arrow-key navigation within each group, Enter/Space to select, visible focus rings on all interactives, modal focus trap + return-focus on close, `aria-live` progress announcements, proper `<label>` associations for notes textareas, `prefers-reduced-motion` support. No additional work for the agent.

## Rules

1. **One squiz per invocation.** If answers reveal new ambiguity, follow up in plain chat — don't render a second squiz.
2. **3-12 questions** is the sweet spot. Fewer → `/quiz`. More → the user bails.
3. **Stable IDs.** Both squiz `id` and option `id` are stable — they come back in the JSON. Pick short kebab-or-camel slugs that will still make sense in a week.
4. **Omit `theme`** unless you have a reason to override. Auto-rotation does the right thing.
5. **Make `art` earn its slot.** Authoring preference is `wf:` / `arch:` > DSL > `"none"` > raw SVG. Use `arch:*` for system-architecture diagrams and `wf:*` for UI wireframes; reach for DSL primitives (`flow:`, `box:`, `arrow:`, `text:`, etc.) when you need to compose; raw SVG only for bespoke metaphors. Use `"none"` instead of forcing irrelevant art.
6. **Spec narrative is optional.** Include it only when you have real prose to quote with `{{markers}}` that map to squizzes.
7. **The `quote` field on a squiz is optional.** Use it when you can point to a specific spec line that motivates the question.
8. **Self-contained.** The doc should make sense to a user opening it cold. `SPEC_LEDE` is the one-liner that does this work.
9. **Clickable links.** When you hand the user the rendered file, format it as a `file://` URL (`file:///C:/Users/.../foo.html` on Windows). Bare paths aren't clickable; URLs are.
10. **Recommend when you have a real preference.** Any option can carry `"recommendation": "<one or two sentences explaining why>"`. The renderer shows a `★ RECOMMENDED` chip + the explanation as an editorial callout. Use it when the spec/constraints/audience genuinely point at one option — DON'T mark every option as recommended, and don't fluff the explanation. If you can't justify the pick in one sentence ("OVR-3's $5-VPS constraint rules out k8s; that leaves systemd vs Docker; systemd is one fewer moving part"), the recommendation isn't earned and should stay off. At most one recommendation per squiz/item.

## Files in this skill

- `SKILL.md` — this file.
- `go-renderer/` — Go module that becomes the `squiz` binary on install.
  - `main.go`, `schema.go`, `render.go`, `theme.go`, `art.go`, `dsl.go`, `wf.go`, `browser.go`, `strbuilder.go`
  - `templates/index.html.tmpl` + `templates/styles.css` — embedded via `//go:embed`
  - `testdata/smoke.json` — reference fixture exercising every art form
