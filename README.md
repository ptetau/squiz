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
go install github.com/ptetau/squiz@latest
mkdir -p ~/.claude/skills/squiz
curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/SKILL.md \
  -o ~/.claude/skills/squiz/SKILL.md
```

The install scripts drop the `squiz` binary on PATH and copy `SKILL.md` into `~/.claude/skills/squiz/` so the Claude Code agent picks it up.

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

## How the skill works

See [SKILL.md](./SKILL.md) for the full agent-facing contract: the JSON schema, the 50 named wireframes, the 7 DSL primitives, the 8 themes, and the export payload shape.

## Build

```sh
go build -o squiz .
```

The templates and CSS are embedded via `//go:embed`, so the binary is fully self-contained.

## License

[MIT](./LICENSE)
