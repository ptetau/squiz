# Squiz

Two visual document tools for [Claude Code](https://github.com/anthropics/claude-code), sharing the Apple //e × IBM Plex aesthetic:

- **`squiz`** — clarifying questions. The agent writes a JSON spec; the CLI renders it as an interactive document with mini-wireframe option art and a `copy json` payload the user pastes back.
- **`squiz-plan`** — structured plans that can include decisions (overview → functional → non-functional → cases → engineering → build). The agent writes a multi-file plan tree; the CLI renders one tabbed HTML doc where every item carries clickable `[FR-3]`-style badges back to the parent items that motivated it, and any item can carry an `options:` chooser when it represents an unsettled call.

Eight themes (paper / phosphor / amber / beige / rose / ocean / forest / slate), auto-rotated per repo so every project gets a distinct identity.

## Install

**macOS / Linux**
```sh
curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/install.sh | sh
```

**Windows (PowerShell)**
```powershell
irm https://raw.githubusercontent.com/ptetau/squiz/main/install.ps1 | iex
```

Drops both `squiz` and `squiz-plan` on PATH; lays down `~/.claude/skills/{squiz,squiz-plan,squiz-update}/SKILL.md` so Claude Code picks them up.

## Update

After install, refresh to the latest release without re-running install:

```sh
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/update.sh | sh

# Windows (PowerShell)
irm https://raw.githubusercontent.com/ptetau/squiz/main/update.ps1 | iex
```

Detects binaries on PATH + SKILL.md files under `~/.claude/skills/` AND `./.claude/skills/` (project-local) and updates whichever it finds. Prompts before replacing — pass `--yes` (sh) / `-Yes` (ps) to skip. Pin a specific tag with `--version 0.5.0` / `-Version 0.5.0` for rollbacks. `--dry-run` / `-DryRun` previews without changes.

In Claude Code:

> /squiz-update

**From source** (Go ≥ 1.22)
```sh
go install github.com/ptetau/squiz/cmd/squiz@latest
go install github.com/ptetau/squiz/cmd/squiz-plan@latest
mkdir -p ~/.claude/skills/squiz ~/.claude/skills/squiz-plan
curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/skills/squiz/SKILL.md \
  -o ~/.claude/skills/squiz/SKILL.md
curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/skills/squiz-plan/SKILL.md \
  -o ~/.claude/skills/squiz-plan/SKILL.md
```

## Quick start

Both binaries can scaffold their own canonical sample so you can see what a real input looks like immediately after install — no source-tree paths to remember.

```sh
# clarifier — scaffold a sample, render + open it
squiz example                            # writes ./squiz-example.json
squiz squiz-example.json --open

# structured plan — scaffold a 7-file sample tree, render + open it
squiz-plan example                       # writes ./squiz-plan-example/
squiz-plan squiz-plan-example/index.json --open

# verify
squiz version
squiz-plan version
```

`squiz example --stdout` and `squiz-plan example --out my-plan` give you finer control. The scaffolded samples are real fixtures that exercise every art form and feature — copy and adapt them.

Then in Claude Code:

> /squiz let's design a personal habit tracker

> /squiz-plan turn the resolved decisions into a structured build plan

**New in v0.6.0 (agent-native CLI surface):** both binaries now expose introspection verbs so agents can author + validate input without reading external docs.

```sh
squiz schema                       # JSON Schema for the input format
squiz validate <input.json>        # parse + business-rule check; exit 0/1
squiz catalog wf                   # list all wireframes with descriptions
squiz catalog wf --previews        # gallery HTML of every wireframe
squiz catalog arch                 # 30 system-design icons
squiz catalog dsl                  # 11 DSL primitives with grammar
squiz catalog themes               # 8 themes with vibe descriptions
squiz preview wf:calendar-grid     # render one art form to a standalone page
squiz help                         # list topical help (art, themes, dsl, …)
squiz help art                     # deep reference on art forms
squiz skill                        # dump the embedded SKILL.md to stdout
```

All catalog/validate verbs accept `--json` for machine output. The same verbs exist on `squiz-plan` (with extra topics: `sections`, `refs`, `notes`, `proposed-items`).

**New in v0.8.0 (composition mechanism):** raw SVG can now embed library + DSL primitives via `<use href="wf:phone-card"/>` / `<use href="arch:database"/>` / `<use href="callout:..."/>`. The library + DSL items are PARTS the agent remixes into bespoke illustrations — not finished pictures to pick from. Seven new annotation primitives (`callout:`, `brace:`, `divider:vs`, `badge:tick/cross/warn/star/dot`, `range:LO-HI`, `baseline:N`, `times:N`) give you the labels, arrows, and baselines that turn nouns into statements. `squiz catalog wf --json` now emits `naturalBox` per entry so `<use>` sizing is precise. `squiz validate` adds composition-health warnings (single-token-heavy sections, sibling-art-collision, missing viewBox, unknown wf/arch refs).

**Earlier highlights:** v0.5.0 added per-option `recommendation` (with explanation); v0.4.0 added plan-item `options:` choosers, three note channels, `arch:*` icons, and the `text:` / `flow:` / `box:` / `arrow:` DSL primitives.

## CLI

```
squiz       <input.json>          [--out path] [--stdout] [--open] [--theme name]
squiz       render <input.json>   [--out path] [--stdout] [--open] [--theme name]
squiz       example               [--out path] [--stdout]
squiz       schema                [--out path]
squiz       validate <input.json> [--json]
squiz       catalog [wf|arch|dsl|themes]  [--json] [--previews [--out path] [--theme name]]
squiz       preview <art-spec>    [--out path] [--stdout] [--theme name]
squiz       help [topic]
squiz       skill                 [--out path]
squiz       version

squiz-plan  <index.json>          [--out path] [--stdout] [--open] [--theme name]
squiz-plan  example               [--out dir]
squiz-plan  schema | validate | catalog | preview | help | skill | version
            (same shapes as squiz)
```

Flags may appear before or after the positional argument. `--help` on any subcommand prints flag-level help; `squiz help <topic>` gives the long-form reference.

## Skills

- **[skills/squiz/SKILL.md](./skills/squiz/SKILL.md)** — full agent contract: JSON schema, 50 named UI wireframes (`wf:*`), 30 system-design icons (`arch:*`, new in v0.4.0), the parametric DSL (including the new `text:` / `flow:` / `box:` / `arrow:` primitives), 8 themes, export payload shape.
- **[skills/squiz-plan/SKILL.md](./skills/squiz-plan/SKILL.md)** — agent contract for structured plans: section model, ID conventions, refs, item `options:` choosers, the three note channels (item / section / plan), proposed-items, feedback shape.

## Layout

```
cmd/squiz/             # squiz CLI entry point
cmd/squiz-plan/        # squiz-plan CLI entry point
pkg/renderer/          # exported library: themes, art, DSL, base templates
internal/planview/     # squiz-plan-specific: parser, render, template
skills/<name>/         # SKILL.md per skill (installed to ~/.claude/skills/<name>/)
testdata/              # smoke + plan fixtures used by tests
```

## Build

```sh
go build -o squiz       ./cmd/squiz
go build -o squiz-plan  ./cmd/squiz-plan
go test ./... -count=1
```

Templates and CSS are embedded via `//go:embed`; both binaries are fully self-contained.

## License

[MIT](./LICENSE)
