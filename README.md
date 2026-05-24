# Squiz

A visual document-style clarifier for [Claude Code](https://github.com/anthropics/claude-code). The agent writes a compact JSON spec; the `squiz` CLI renders it as a self-contained, retro-styled interactive HTML doc with mini-wireframe option art and a sticky `copy json` status bar. The user fills it in at their own pace and pastes the JSON payload back.

Apple //e × IBM Plex aesthetic. Eight themes (paper / phosphor / amber / beige / rose / ocean / forest / slate), auto-rotated per repo so every project gets a distinct identity.

## Install

**macOS / Linux**
```sh
curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/install.sh | sh
```

**Windows (PowerShell)**
```powershell
irm https://raw.githubusercontent.com/ptetau/squiz/main/install.ps1 | iex
```

**From source** (Go ≥ 1.22)
```sh
go install github.com/ptetau/squiz/cmd/squiz@latest
mkdir -p ~/.claude/skills/squiz
curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/skills/squiz/SKILL.md \
  -o ~/.claude/skills/squiz/SKILL.md
```

The install scripts drop the `squiz` binary on PATH and copy each `skills/<name>/SKILL.md` into `~/.claude/skills/<name>/` so the Claude Code agent picks them up. Today that's `squiz` plus a placeholder `squiz-plan` stub (full implementation in v0.3.0).

## Quick start

```sh
# render a sample
squiz testdata/smoke.json --open

# verify install
squiz version
```

Then ask Claude Code:

> /squiz let's design a personal habit tracker

The agent will write a `.json` spec, run the binary, and hand you back the rendered HTML.

## CLI

```
squiz <input.json>                    # render <input>.html next to input
squiz <input.json> --open             # also open in default browser
squiz <input.json> --theme phosphor   # force a specific theme
squiz <input.json> --out path.html    # explicit output path
squiz <input.json> --stdout > x.html  # write to stdout
squiz version
squiz help
```

## How the skills work

- **[skills/squiz/SKILL.md](./skills/squiz/SKILL.md)** — the active skill: JSON schema, 50 named wireframes, 7 DSL primitives, 8 themes, export payload shape.
- **[skills/squiz-plan/SKILL.md](./skills/squiz-plan/SKILL.md)** — sibling skill (placeholder; v0.3.0).

## Layout

```
cmd/squiz/        # CLI entry point
pkg/renderer/     # exported library: themes, art, DSL, templates
skills/<name>/    # SKILL.md per skill
```

## Build

```sh
go build -o squiz ./cmd/squiz
```

The templates and CSS are embedded via `//go:embed` in `pkg/renderer`, so the binary is fully self-contained.

## License

[MIT](./LICENSE)
