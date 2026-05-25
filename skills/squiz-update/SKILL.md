---
name: squiz-update
description: >
  Updates the squiz + squiz-plan binaries AND their SKILL.md files to the
  latest GitHub release (or a pinned --version). Detects where the
  binaries and skills are currently installed and replaces them in place
  — both global (~/.claude/skills/) and project-local (./.claude/skills/)
  are honored. Triggered by /squiz-update. Use when the user wants the
  toolchain refreshed, when "squiz version" looks behind, or when rolling
  back from a bad release.
---

# squiz-update — keep the toolchain current

## When the user types /squiz-update

Run one command and report what changed. Do NOT prompt them to confirm again — typing `/squiz-update` IS the confirmation, so pass `--yes` to suppress the script's prompt.

**macOS / Linux:**
```sh
curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/update.sh | sh -s -- --yes
```

**Windows (PowerShell):**
```powershell
$u = irm https://raw.githubusercontent.com/ptetau/squiz/main/update.ps1; iex "$u; Update-Squiz -Yes"
```

After it runs, parse the script's stderr/stdout and tell the user **what actually changed** — at minimum: old version → new version, and which paths got touched. If the script reported "already at vX.Y.Z — nothing to do", say so verbatim; don't pretend something happened.

## When the user pins a version

`/squiz-update 0.5.0` → invoke with `--version 0.5.0` (or PowerShell `-Version 0.5.0`). This is both for forward-pinning and for rollbacks when a release ships a regression.

```sh
curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/update.sh | sh -s -- --yes --version 0.5.0
```

## What the script does (so you can describe it accurately)

1. **Detects installed binaries** — `squiz` and/or `squiz-plan` via `command -v` (`Get-Command` on Windows). If neither is found, errors and points at `install.sh`/`install.ps1`. It doesn't try to install fresh — that's a different command.
2. **Reads current versions** by running `squiz version` / `squiz-plan version`.
3. **Resolves target version** — `--version <X.Y.Z>` if passed, otherwise the GitHub API's "latest release" tag.
4. **Detects SKILL.md locations** — checks both `~/.claude/skills/<name>/SKILL.md` (global) AND `./.claude/skills/<name>/SKILL.md` (project-local) for `squiz`, `squiz-plan`, and `squiz-update`. Updates whichever exist; never creates new ones (so an install-script install only refreshes global, a project-local install only refreshes project).
5. **Short-circuits** when versions match AND no SKILL.md files are due to refresh — prints "already at vX.Y.Z — nothing to do" and exits 0.
6. **Prompts to confirm** (unless `--yes` / `-Yes`); shows the plan first (e.g. `squiz 0.6.0 → 0.7.0 (path)`, `skill: ~/.claude/skills/squiz/SKILL.md`).
7. **Downloads + verifies SHA256** from the release archive on GitHub.
8. **Replaces each binary at its current absolute path.** On Windows, renames the old `.exe` to `.exe.old` first because Windows can't overwrite a running `.exe`. The `.old` file is locked until the process exits; not a problem in practice.
9. **Replaces each existing SKILL.md** with the archive's copy.
10. **Reports the result** to stderr.

## Other flags worth knowing

| Flag | sh | ps1 | What it does |
|---|---|---|---|
| Pin version | `--version 0.5.0` | `-Version 0.5.0` | Update to / roll back to a specific tag |
| Skip confirm | `--yes` / `-y` | `-Yes` | Use when invoked by /squiz-update |
| Preview only | `--dry-run` | `-DryRun` | Show the plan, change nothing |
| Help | `--help` / `-h` | `Get-Help Update-Squiz` | Usage |

## Rules

1. **The slash command IS the consent.** Always pass `--yes` (or `-Yes`). Never make the user confirm twice.
2. **Always report what changed.** Quote the script's last lines so the user sees the version transition + the files touched. If nothing changed, say "already at vX.Y.Z".
3. **`--version` is a rollback too.** "Roll back to 0.5" → `--version 0.5.0`. Useful when a release ships a regression.
4. **Don't install fresh from update.** If the script reports "no squiz binaries on PATH", tell the user to run `install.sh` / `install.ps1` first — don't try to download the binary another way.
5. **Project-local installs are real.** If `./.claude/skills/squiz/SKILL.md` exists, the script updates it AND the global one if global exists too. Both stay in sync after a single update.
6. **Use `--dry-run` if the user asks "what would change?"** — runs all the detection, prints the plan, makes no changes.
7. **No need to call /squiz or /squiz-plan after** — the binary and skill are atomically aligned post-update.
