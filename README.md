# gh-pr-quality-gate

`gh-pr-quality-gate` is a cross-agent repository quality gate. It gives teams a
single policy for issue-linked pull requests, local verification, required CI,
contract documentation, review evidence, and explicit authority boundaries.

The project includes:

- a precompiled, cross-platform GitHub CLI extension;
- an Agent Skill shared by Codex, ChatGPT, and Claude Code plugins;
- Codex and Claude marketplace manifests;
- repository instruction entry points for Codex, Claude, GitHub Copilot,
  Gemini, and other agents;
- collision-safe repository templates, pre-push hooks, and GitHub Actions;
- dry-run branch-protection planning with exact apply confirmation.

It deliberately does not merge pull requests, force-push, change Draft state,
or close work items. Those operations remain under explicit human authority.

## Install the GitHub CLI extension

Prerequisites: [GitHub CLI](https://cli.github.com/) and an authenticated
`gh auth login` session.

```shell
gh extension install Ysoseri1224/gh-pr-quality-gate
gh pr-quality-gate version
```

The release contains binaries for Windows, Linux, and macOS. To build locally:

```shell
go install github.com/Ysoseri1224/gh-pr-quality-gate/cmd/gh-pr-quality-gate@latest
gh-pr-quality-gate version
```

## Install the Agent plugin

Codex and ChatGPT Work:

```shell
codex plugin marketplace add Ysoseri1224/gh-pr-quality-gate
```

Restart the ChatGPT desktop app, open the Plugins Directory, select
`Ysoseri1224 Quality Gates`, and install `gh-pr-quality-gate`.

Claude Code:

```shell
claude plugin marketplace add Ysoseri1224/gh-pr-quality-gate
claude plugin install gh-pr-quality-gate@ysoseri1224-quality-gates
```

Direct Skill installation and repository-level integrations are documented in
[`docs/installation.md`](docs/installation.md) and
[`docs/platform-integrations.md`](docs/platform-integrations.md).

## Use it

Audit an existing repository without changing it:

```shell
gh pr-quality-gate audit --repo /path/to/repository
```

Preview and then install repository files:

```shell
gh pr-quality-gate init --repo /path/to/repository
gh pr-quality-gate init --repo /path/to/repository --apply
```

Edit `.github/pr-quality-gate.yml` to contain the repository's real local test,
lint, type-check, and build commands. Then validate before pushing:

```shell
gh pr-quality-gate validate --repo /path/to/repository --run-local
```

Preview branch protection. Applying requires the exact repository and branch as
a second confirmation:

```shell
gh pr-quality-gate protect --repository owner/repo
gh pr-quality-gate protect --repository owner/repo --apply --confirm owner/repo:main
```

Do not enable required checks until their workflows have completed successfully
on the target repository. See [`docs/configuration.md`](docs/configuration.md).

## Safety model

- `audit`, `init`, and `protect` are non-mutating by default.
- `init --apply` creates only missing files and refuses content collisions.
- `protect --apply` requires an exact `owner/repo:branch` confirmation.
- `protect --apply` refuses to replace existing branch protection.
- PR policy runs against the trusted base commit and does not execute PR code.
- Candidate-code CI has read-only permissions and receives no repository secrets.
- Local hooks provide early feedback; required CI is the shared enforcement gate.

Security details and threat boundaries are in [`SECURITY.md`](SECURITY.md).

## License

[MIT](LICENSE)
