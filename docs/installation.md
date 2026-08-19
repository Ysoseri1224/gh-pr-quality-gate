# Installation

## GitHub CLI extension

The precompiled extension is the recommended installation on Windows, Linux,
and macOS:

```shell
gh auth login
gh extension install Ysoseri1224/gh-pr-quality-gate
gh pr-quality-gate version
```

Update or remove it with:

```shell
gh extension upgrade pr-quality-gate
gh extension remove pr-quality-gate
```

For local development, clone the repository and run:

```shell
go build -o gh-pr-quality-gate ./cmd/gh-pr-quality-gate
gh extension install .
```

On Windows, build `gh-pr-quality-gate.exe` instead.

## Claude Code plugin

```shell
claude plugin marketplace add Ysoseri1224/gh-pr-quality-gate
claude plugin install gh-pr-quality-gate@ysoseri1224-quality-gates
```

Use `--scope project` on the marketplace command only when the team intends to
commit the marketplace declaration for that repository.

For local validation and installation:

```shell
claude plugin validate /absolute/path/to/gh-pr-quality-gate
claude plugin marketplace add /absolute/path/to/gh-pr-quality-gate
claude plugin install gh-pr-quality-gate@ysoseri1224-quality-gates
```

## Direct Agent Skill installation

Direct installation is useful when a platform supports Agent Skills but not the
plugin marketplace. Clone this repository, then run one of the bundled scripts.

Windows PowerShell:

```powershell
.\scripts\install-skill.ps1 -Target Codex
.\scripts\install-skill.ps1 -Target Claude
.\scripts\install-skill.ps1 -Target Gemini
.\scripts\install-skill.ps1 -Target All
```

Linux or macOS:

```shell
./scripts/install-skill.sh codex
./scripts/install-skill.sh claude
./scripts/install-skill.sh gemini
./scripts/install-skill.sh all
```

The installers refuse to replace an existing skill directory. Remove or back up
the old installation deliberately before installing a different version.

## Add the gate to a repository

Run a dry run first:

```shell
gh pr-quality-gate init --repo /path/to/repository
```

`conflict` means a target file already exists with different content. The tool
never overwrites it. Reconcile that file manually, rerun the preview, and apply
only when the plan is acceptable:

```shell
gh pr-quality-gate init --repo /path/to/repository --apply
```

Configure the actual local commands before enabling the generated workflows or
required branch checks. The generated pre-push hook is optional:

```powershell
.\scripts\install-pre-push.ps1
```

```shell
./scripts/install-pre-push.sh
```
