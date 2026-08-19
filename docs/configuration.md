# Configuration

Repository policy lives in `.github/pr-quality-gate.yml`.

```yaml
version: 1
local_gate:
  windows:
    - ./scripts/check.ps1 -SkipInstall
  posix:
    - ./scripts/check.sh
required_checks:
  - Repository quality gate
  - PR policy
required_files:
  - AGENTS.md
  - docs/repo_rule.md
  - .github/pull_request_template.md
issue_reference:
  required: true
  allow_partial_refs: true
agent_authority:
  merge: explicit-authorisation
  draft_transition: explicit-authorisation
  force_push: prohibited
branch_protection:
  branch: main
  required_approvals: 1
  dismiss_stale_reviews: true
  require_last_push_approval: true
  require_conversation_resolution: true
  enforce_admins: true
```

## Local gate

List commands in their required order. Each command runs from the repository
root and stops the gate on its first non-zero exit. Keep installation separate
from verification where practical so a pre-push check is deterministic and
does not silently change dependency locks.

Review configuration from an untrusted branch before using `--run-local`; these
commands execute with the current user's permissions.

## Required checks

Names must match the job names reported by GitHub Actions. Do not apply branch
protection until each named check has run successfully on the target repository.
Renaming a workflow or job can leave a protected branch waiting for a check that
will never report.

## Issue references

When `required` is true, the pull request body must contain an issue reference.
Closing forms include `Closes #123`, `Fixes #123`, and `Resolves #123`. If
`allow_partial_refs` is true, non-closing forms such as `Refs #123` are accepted.

Use a closing form only when the pull request fully satisfies the issue.

## Agent authority

The three authority values are intentionally fixed. Configuration cannot grant
an agent implicit merge or Draft-transition authority, and force-pushing remains
prohibited. This prevents a repository branch from weakening the tool's own
control boundary.

## Branch protection

`protect` maps the configured checks and review settings to GitHub's branch
protection API. It is a dry run unless both `--apply` and the exact
`--confirm owner/repo:branch` value are supplied. Repository administrator
permission may be required. Version 0.1 creates protection only on an
unprotected branch. It refuses to replace existing protection because doing so
could remove restrictions or settings outside this configuration model; merge
those settings manually in GitHub instead.
