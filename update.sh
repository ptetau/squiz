#!/usr/bin/env sh
# Update squiz + squiz-plan to the latest GitHub release (or a pinned
# --version). Detects where the binaries + SKILL.md files actually live
# and replaces them in place — global (~/.claude/skills/) AND project-
# local (./.claude/skills/) are both honored.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/update.sh | sh
#   curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/update.sh | sh -s -- --yes
#   curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/update.sh | sh -s -- --version 0.5.0
#   curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/update.sh | sh -s -- --dry-run
#
# Env overrides (alternatives to flags):
#   VERSION   pin a specific version
#   YES=1     skip the confirmation prompt
#   DRY_RUN=1 print the plan but don't change anything
set -eu

OWNER="ptetau"
REPO="squiz"
TARGET="${VERSION:-}"
ASSUME_YES="${YES:-0}"
DRY_RUN="${DRY_RUN:-0}"

err() { printf '%s\n' "$*" >&2; exit 1; }
info() { printf '%s\n' "$*"; }

while [ $# -gt 0 ]; do
  case "$1" in
    --version)   TARGET="$2"; shift 2 ;;
    --version=*) TARGET="${1#*=}"; shift ;;
    --yes|-y)    ASSUME_YES=1; shift ;;
    --dry-run)   DRY_RUN=1; shift ;;
    --help|-h)
      sed -n '2,17p' "$0" | sed 's|^# \{0,1\}||'
      exit 0 ;;
    *) err "unknown flag: $1 (try --help)" ;;
  esac
done

# 1. Detect installed binaries.
SQUIZ_BIN=$(command -v squiz 2>/dev/null || true)
SQP_BIN=$(command -v squiz-plan 2>/dev/null || true)
if [ -z "$SQUIZ_BIN" ] && [ -z "$SQP_BIN" ]; then
  err "no squiz binaries on PATH — run install.sh first:
  curl -fsSL https://raw.githubusercontent.com/$OWNER/$REPO/main/install.sh | sh"
fi

current_version() {
  # `squiz version` prints "squiz 0.6.0" — grab the second field.
  "$1" version 2>/dev/null | awk '{print $2}' || echo "?"
}

SQUIZ_VER=""
SQP_VER=""
[ -n "$SQUIZ_BIN" ] && SQUIZ_VER=$(current_version "$SQUIZ_BIN")
[ -n "$SQP_BIN" ]   && SQP_VER=$(current_version "$SQP_BIN")

# 2. Resolve target version.
if [ -z "$TARGET" ]; then
  TARGET=$(curl -fsSL "https://api.github.com/repos/$OWNER/$REPO/releases/latest" \
    | grep '"tag_name"' | head -1 | sed -E 's/.*"v?([^"]+)".*/\1/')
  [ -n "$TARGET" ] || err "could not resolve latest version from GitHub API (set --version to pin)"
fi

# 3. Detect skill directories. Update only what already exists.
SKILL_ROOTS=""
add_root() {
  if [ -d "$1" ]; then
    case " $SKILL_ROOTS " in
      *" $1 "*) ;;
      *) SKILL_ROOTS="$SKILL_ROOTS $1" ;;
    esac
  fi
}
add_root "$HOME/.claude/skills"
add_root "$(pwd)/.claude/skills"

# Find the existing SKILL.md files under those roots that we ship.
EXISTING_SKILLS=""
for root in $SKILL_ROOTS; do
  for name in squiz squiz-plan squiz-update; do
    path="$root/$name/SKILL.md"
    if [ -f "$path" ]; then
      EXISTING_SKILLS="$EXISTING_SKILLS $path"
    fi
  done
done

# 4. Compare and short-circuit if already current. We key on binary
# version only: SKILL.md files travel WITH the binary release, so if
# the binary version matches the target there's no useful refresh to
# do (and re-downloading the archive just to copy the same SKILL.md
# over itself wastes the user's time + network).
need_binary_update=0
[ -n "$SQUIZ_BIN" ] && [ "$SQUIZ_VER" != "$TARGET" ] && need_binary_update=1
[ -n "$SQP_BIN" ]   && [ "$SQP_VER" != "$TARGET" ]   && need_binary_update=1

if [ "$need_binary_update" = "0" ]; then
  info "already at v$TARGET — nothing to do"
  exit 0
fi

# 5. Show plan.
info "update plan:"
[ -n "$SQUIZ_BIN" ] && info "  squiz       $SQUIZ_VER → $TARGET   ($SQUIZ_BIN)"
[ -n "$SQP_BIN" ]   && info "  squiz-plan  $SQP_VER → $TARGET   ($SQP_BIN)"
for s in $EXISTING_SKILLS; do
  info "  skill       (refresh)         ($s)"
done

if [ "$DRY_RUN" = "1" ]; then
  info "[dry-run] no changes made"
  exit 0
fi

# 6. Confirm.
if [ "$ASSUME_YES" != "1" ]; then
  # Use /dev/tty so confirmation works even when piped through `sh`.
  if [ -t 0 ] || [ -r /dev/tty ]; then
    printf "proceed? [y/N] " > /dev/tty 2>/dev/null || printf "proceed? [y/N] "
    if [ -r /dev/tty ]; then
      read -r answer < /dev/tty
    else
      read -r answer
    fi
    case "$answer" in
      y|Y|yes|YES|Yes) ;;
      *) info "cancelled"; exit 1 ;;
    esac
  else
    err "no terminal for confirmation — re-run with --yes (or set YES=1)"
  fi
fi

# 7. Detect OS + arch for archive selection.
os=$(uname -s)
arch=$(uname -m)
case "$os" in
  Linux|Darwin) ;;
  *) err "unsupported OS: $os (this script handles Linux and Darwin; use update.ps1 for Windows)" ;;
esac
case "$arch" in
  x86_64|amd64)  arch=x86_64 ;;
  arm64|aarch64) arch=arm64 ;;
  *) err "unsupported arch: $arch" ;;
esac

archive="squiz_${TARGET}_${os}_${arch}.tar.gz"
base_url="https://github.com/$OWNER/$REPO/releases/download/v${TARGET}"

# 8. Download + verify.
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

info "→ downloading $archive"
curl -fsSL "$base_url/$archive"      -o "$tmp/$archive" || err "download failed (does v$TARGET exist as a release?)"
curl -fsSL "$base_url/checksums.txt" -o "$tmp/checksums.txt" || err "could not fetch checksums.txt"

info "→ verifying checksum"
if command -v sha256sum >/dev/null 2>&1; then
  (cd "$tmp" && grep " $archive\$" checksums.txt | sha256sum -c -)
elif command -v shasum >/dev/null 2>&1; then
  (cd "$tmp" && grep " $archive\$" checksums.txt | shasum -a 256 -c -)
else
  info "⚠ no sha256sum/shasum found; skipping checksum verification"
fi

info "→ extracting"
tar -C "$tmp" -xzf "$tmp/$archive"

# 9. Replace binaries (in their current absolute paths).
replace_bin() {
  local current="$1" name="$2"
  [ -n "$current" ] || return 0
  if [ ! -e "$tmp/$name" ]; then
    info "⚠ archive missing $name; skipping"
    return 0
  fi
  install -m 0755 "$tmp/$name" "$current"
  info "  binary: $current"
}
replace_bin "$SQUIZ_BIN" squiz
replace_bin "$SQP_BIN"   squiz-plan

# 10. Replace SKILL.mds — only at locations where they already exist.
for path in $EXISTING_SKILLS; do
  name=$(basename "$(dirname "$path")")
  src="$tmp/skills/$name/SKILL.md"
  if [ -f "$src" ]; then
    install -m 0644 "$src" "$path"
    info "  skill:  $path"
  else
    info "  ⚠ archive missing skills/$name/SKILL.md; left $path untouched"
  fi
done

info "✓ updated to v$TARGET"
[ -n "$SQUIZ_BIN" ] && "$SQUIZ_BIN" version 2>/dev/null || true
[ -n "$SQP_BIN" ]   && "$SQP_BIN"   version 2>/dev/null || true
