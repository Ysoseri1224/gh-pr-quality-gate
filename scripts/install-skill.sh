#!/bin/sh
set -eu

target=${1:-all}
case "$target" in
  codex|claude|gemini|all) ;;
  *) echo "usage: $0 [codex|claude|gemini|all]" >&2; exit 2 ;;
esac

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
source_dir="$script_dir/../plugins/gh-pr-quality-gate/skills/gh-pr-quality-gate"

install_one() {
  platform=$1
  destination=$2
  mkdir -p "$(dirname "$destination")"
  cp -R "$source_dir" "$destination"
  echo "Installed gh-pr-quality-gate for $platform at $destination"
}

check_destination() {
  destination=$1
  if [ -e "$destination" ]; then
    echo "Refusing to replace existing installation: $destination" >&2
    exit 1
  fi
}

if [ "$target" = "codex" ] || [ "$target" = "all" ]; then
  check_destination "$HOME/.codex/skills/gh-pr-quality-gate"
fi
if [ "$target" = "claude" ] || [ "$target" = "all" ]; then
  check_destination "$HOME/.claude/skills/gh-pr-quality-gate"
fi
if [ "$target" = "gemini" ] || [ "$target" = "all" ]; then
  check_destination "$HOME/.gemini/skills/gh-pr-quality-gate"
fi

if [ "$target" = "codex" ] || [ "$target" = "all" ]; then
  install_one Codex "$HOME/.codex/skills/gh-pr-quality-gate"
fi
if [ "$target" = "claude" ] || [ "$target" = "all" ]; then
  install_one Claude "$HOME/.claude/skills/gh-pr-quality-gate"
fi
if [ "$target" = "gemini" ] || [ "$target" = "all" ]; then
  install_one Gemini "$HOME/.gemini/skills/gh-pr-quality-gate"
fi
