#!/usr/bin/env sh
# Install squiz: binary on PATH + SKILL.md files into ~/.claude/skills/<skill>/.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/ptetau/squiz/main/install.sh | sh
#
# Env overrides:
#   VERSION             pin a specific version (default: latest GitHub release)
#   SQUIZ_BIN_DIR       where to install the binary (default: ~/.local/bin)
#   SQUIZ_SKILLS_ROOT   where skill dirs land (default: ~/.claude/skills)
set -eu

OWNER="ptetau"
REPO="squiz"
BIN_DIR="${SQUIZ_BIN_DIR:-$HOME/.local/bin}"
SKILLS_ROOT="${SQUIZ_SKILLS_ROOT:-$HOME/.claude/skills}"

err() { printf '%s\n' "$*" >&2; exit 1; }

os=$(uname -s)
arch=$(uname -m)
case "$os" in
  Linux|Darwin) ;;
  *) err "unsupported OS: $os (this script handles Linux and Darwin; use install.ps1 for Windows)" ;;
esac
case "$arch" in
  x86_64|amd64)  arch=x86_64 ;;
  arm64|aarch64) arch=arm64 ;;
  *) err "unsupported arch: $arch" ;;
esac

version="${VERSION:-}"
if [ -z "$version" ]; then
  version=$(curl -fsSL "https://api.github.com/repos/$OWNER/$REPO/releases/latest" \
    | grep '"tag_name"' | head -1 | sed -E 's/.*"v?([^"]+)".*/\1/')
  [ -n "$version" ] || err "could not resolve latest version (set VERSION=… to pin)"
fi

archive="squiz_${version}_${os}_${arch}.tar.gz"
base_url="https://github.com/$OWNER/$REPO/releases/download/v${version}"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

printf '→ downloading %s\n' "$archive"
curl -fsSL "$base_url/$archive"       -o "$tmp/$archive"
curl -fsSL "$base_url/checksums.txt"  -o "$tmp/checksums.txt"

printf '→ verifying checksum\n'
if command -v sha256sum >/dev/null 2>&1; then
  (cd "$tmp" && grep " $archive\$" checksums.txt | sha256sum -c -)
elif command -v shasum >/dev/null 2>&1; then
  (cd "$tmp" && grep " $archive\$" checksums.txt | shasum -a 256 -c -)
else
  printf '⚠ no sha256sum/shasum found; skipping checksum verification\n' >&2
fi

printf '→ extracting\n'
tar -C "$tmp" -xzf "$tmp/$archive"

mkdir -p "$BIN_DIR"
# Install every binary the archive ships at its top level (squiz, squiz-plan, …).
for bin_src in "$tmp"/squiz "$tmp"/squiz-plan; do
  [ -e "$bin_src" ] || continue
  bin_name=$(basename "$bin_src")
  install -m 0755 "$bin_src" "$BIN_DIR/$bin_name"
  printf '  binary:  %s/%s\n' "$BIN_DIR" "$bin_name"
done

# Install every SKILL.md the archive ships under skills/<name>/.
for skill_src in "$tmp"/skills/*/SKILL.md; do
  [ -e "$skill_src" ] || continue
  skill_name=$(basename "$(dirname "$skill_src")")
  skill_dst="$SKILLS_ROOT/$skill_name"
  mkdir -p "$skill_dst"
  install -m 0644 "$skill_src" "$skill_dst/SKILL.md"
  printf '  skill:   %s/SKILL.md\n' "$skill_dst"
done

case ":$PATH:" in
  *":$BIN_DIR:"*) ;;
  *)
    printf '\n⚠ %s is not on PATH. Add to your shell rc:\n' "$BIN_DIR"
    printf '    export PATH="%s:$PATH"\n' "$BIN_DIR"
    ;;
esac

printf '✓ installed squiz %s\n' "$version"
"$BIN_DIR/squiz" version 2>/dev/null || true
if [ -x "$BIN_DIR/squiz-plan" ]; then
  "$BIN_DIR/squiz-plan" version 2>/dev/null || true
fi
