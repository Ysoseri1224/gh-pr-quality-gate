# Platform Integrations

## Codex and ChatGPT Work

The repository contains a Codex plugin manifest at
`plugins/gh-pr-quality-gate/.codex-plugin/plugin.json` and a repo marketplace at
`.agents/plugins/marketplace.json`. The plugin packages the shared Agent Skill;
it does not add an MCP server or request access to a connected service.

## Claude Code

The same plugin directory contains a Claude manifest, and
`.claude-plugin/marketplace.json` exposes it through Claude's Git-backed
marketplace. Claude discovers the shared `skills/gh-pr-quality-gate/SKILL.md`
without a duplicate workflow definition.

## GitHub Copilot

`gh pr-quality-gate init --apply` creates
`.github/copilot-instructions.md`. It directs Copilot to the same `AGENTS.md`
and `docs/repo_rule.md` used by other agents. GitHub Actions enforce the
machine-checkable subset independently of whether Copilot read the file.

## Gemini

Repository initialization creates `GEMINI.md` as an entry point to the shared
repository policy. The direct installer can also place the Agent Skill under
`~/.gemini/skills/gh-pr-quality-gate` for Gemini environments that support Agent
Skills. Repository instructions and CI remain the enforcement boundary when a
specific Gemini client does not load that directory.

## Other coding agents

Agents that recognize `AGENTS.md` receive the repository-wide entry point.
Agents that recognize neither the file nor Agent Skills can still use the CLI,
pre-push hook, PR template, and required GitHub Actions. Do not treat model
compliance as a substitute for CI or branch protection.

## GitHub Actions

The generated `PR policy` workflow checks trusted base-branch policy and does
not execute pull-request code. The generated `Repository quality gate` workflow
runs repository-defined candidate commands with read-only permissions and no
additional secrets.

This repository also publishes reusable workflows:

```yaml
jobs:
  pr-policy:
    uses: Ysoseri1224/gh-pr-quality-gate/.github/workflows/pr-policy-reusable.yml@v0.1.0
```

Pin an immutable release tag or commit. Review updates before changing the pin.
