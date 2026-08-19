#!/usr/bin/env python3
"""Validate repository packaging without external Python dependencies."""

from __future__ import annotations

import json
import pathlib
import re
import sys


ROOT = pathlib.Path(__file__).resolve().parents[1]
PLUGIN = ROOT / "plugins" / "gh-pr-quality-gate"
SKILL = PLUGIN / "skills" / "gh-pr-quality-gate" / "SKILL.md"


def load_json(relative: str) -> dict:
    path = ROOT / relative
    try:
        value = json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        raise AssertionError(f"invalid JSON at {relative}: {exc}") from exc
    if not isinstance(value, dict):
        raise AssertionError(f"{relative} must contain a JSON object")
    return value


def main() -> int:
    codex = load_json("plugins/gh-pr-quality-gate/.codex-plugin/plugin.json")
    claude = load_json("plugins/gh-pr-quality-gate/.claude-plugin/plugin.json")
    codex_market = load_json(".agents/plugins/marketplace.json")
    claude_market = load_json(".claude-plugin/marketplace.json")

    versions = {
        codex.get("version"),
        claude.get("version"),
        claude_market["plugins"][0].get("version"),
    }
    if len(versions) != 1 or None in versions:
        raise AssertionError(f"plugin versions do not match: {versions}")
    for manifest in (codex, claude):
        if manifest.get("name") != "gh-pr-quality-gate":
            raise AssertionError("plugin name must be gh-pr-quality-gate")
        if manifest.get("license") != "MIT":
            raise AssertionError("plugin license must be MIT")

    if codex_market["plugins"][0]["source"]["path"] != "./plugins/gh-pr-quality-gate":
        raise AssertionError("Codex marketplace source path is invalid")
    if claude_market["plugins"][0]["source"] != "./plugins/gh-pr-quality-gate":
        raise AssertionError("Claude marketplace source path is invalid")

    skill_text = SKILL.read_text(encoding="utf-8")
    match = re.match(r"^---\nname: ([^\n]+)\ndescription: ([^\n]+)\n---\n", skill_text)
    if not match:
        raise AssertionError("SKILL.md requires name and description frontmatter")
    if match.group(1) != "gh-pr-quality-gate":
        raise AssertionError("SKILL.md name is invalid")
    if "TODO" in skill_text or "[TODO" in skill_text:
        raise AssertionError("SKILL.md contains a TODO placeholder")

    print(f"Package metadata is valid for version {versions.pop()}.")
    return 0


if __name__ == "__main__":
    try:
        raise SystemExit(main())
    except (AssertionError, KeyError, IndexError) as exc:
        print(f"validation error: {exc}", file=sys.stderr)
        raise SystemExit(1)
