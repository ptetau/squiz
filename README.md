# Squiz

Two visual document tools for [Claude Code](https://github.com/anthropics/claude-code), sharing the Apple //e × IBM Plex aesthetic:

- **`squiz`** — clarifying questions. The agent writes a JSON spec; the CLI renders it as an interactive document with mini-wireframe option art and a `copy json` payload the user pastes back.
- **`squiz-plan`** — structured plans (overview → functional → non-functional → cases → engineering → build). The agent writes a multi-file plan tree; the CLI renders one tabbed HTML doc where every item carries clickable `[FR-3]`-style badges back to the parent items that motivated it.

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

Drops both `squiz` and `squiz-plan` on PATH; lays down `~/.claude/skills/squiz/SKILL.md` and `~/.claude/skills/squiz-plan/SKILL.md` so Claude Code picks them up.

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

```sh
# clarifier
squiz testdata/smoke.json --open

# structured plan
squiz-plan testdata/plan-example/index.json --open

# verify
squiz version
squiz-plan version
```

Then in Claude Code:

> /squiz let's design a personal habit tracker

> /squiz-plan turn the resolved decisions into a structured build plan

## CLI

```
squiz       <input.json>   [--out path] [--stdout] [--open] [--theme name]
squiz-plan  <index.json>   [--out path] [--stdout] [--open] [--theme name]
```

Both accept flags before or after the positional argument. Both support `version` and `help` subcommands.

## Skills

- **[skills/squiz/SKILL.md](./skills/squiz/SKILL.md)** — full agent contract: JSON schema, 50 named wireframes, 7 DSL primitives, 8 themes, export payload shape.
- **[skills/squiz-plan/SKILL.md](./skills/squiz-plan/SKILL.md)** — agent contract for structured plans: section model, ID conventions, refs, feedback shape.

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
