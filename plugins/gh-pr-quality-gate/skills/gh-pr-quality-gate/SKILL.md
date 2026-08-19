---
name: gh-pr-quality-gate
description: Audit, install, validate, and review GitHub pull request quality gates across coding agents, local hooks, CI, issue linkage, documentation contracts, and branch protection. Use when starting repository governance, preparing or reviewing a pull request, checking whether local and required CI passed, adding AGENTS.md or agent-specific instruction entry points, configuring a safe pre-push hook, or proposing branch-protection changes without granting an agent merge, Draft-transition, force-push, or work-item closure authority.
---

# GitHub PR Quality Gate

Apply one evidence-based delivery policy across local development, coding agents,
pull requests, and GitHub enforcement. Treat repository documentation and CI as
the authority; never treat model compliance alone as a quality or safety gate.

## Establish context

1. Read every applicable `AGENTS.md` or `AGENT.md`, then the repository's current
   product, architecture, API, development, security, and testing sources.
2. Inspect `git status`, the current branch, remotes, open pull request metadata,
   and existing workflows before proposing changes.
3. Identify the issue, acceptance criteria, dependencies, and exact authority
   granted by the user. Preserve unrelated work.
4. Read [policy-boundaries.md](references/policy-boundaries.md) before any
   operation that could alter GitHub state or repository policy.

## Choose the operation

- **Audit:** run `gh pr-quality-gate audit`. Report missing policy, required
  files, check names, issue linkage evidence, and authority boundaries. Do not
  edit files.
- **Bootstrap:** run `gh pr-quality-gate init` first. Review every `add`,
  `unchanged`, and `conflict`. Run again with `--apply` only after authorization.
  Never replace conflicts automatically.
- **Validate:** run `gh pr-quality-gate validate --run-local` before pushing.
  When validating a pull request, supply its body with `--pr-body-file` and
  verify that closing keywords match actual completion.
- **Review:** compare the issue, current contracts, diff, tests, and actual CI.
  Approve only when ready; request changes for concrete blockers. Do not merge.
- **Protect:** run `gh pr-quality-gate protect` as a dry run. Apply only after
  the user explicitly authorizes the exact repository and branch, and supply
  the required `--confirm owner/repo:branch` value.

Read [configuration.md](references/configuration.md) when adding or changing the
configuration. Use the templates under `assets/repository-templates/` only when
the CLI is unavailable, and preserve existing files.

## Enforce evidence boundaries

- Require every pull request to use `Closes #N` only for fully completed work or
  `Refs #N` for partial contribution.
- Run the repository-defined local commands before push. Report exact commands
  and meaningful outcomes; do not claim end-user behavior from unit/API checks.
- Require authoritative API, schema, command, migration, or workflow docs to
  change in the same pull request as their contract.
- Treat a pre-push hook as early feedback. Treat required CI and branch
  protection as the shared enforcement boundary.
- Inspect current GitHub state before reviews or protection changes. Do not rely
  on a pull request description as proof that checks passed.

## Preserve human authority

Do not merge, force-push, close or reopen issues or pull requests, publish a
release, deploy, change Draft state, or alter repository rules unless the user
explicitly authorizes that exact operation. Configuration cannot waive this
boundary. Stop and explain the target, impact, evidence, and required authority
when permission is missing.

## Report the result

State:

- repository and branch inspected;
- issue and pull request linkage;
- files added, unchanged, or conflicting;
- local commands and CI checks actually observed;
- documentation or contract changes required;
- review readiness and remaining blockers;
- any external action intentionally not performed.

Keep the result factual and distinguish verified state from recommendations.

## Bundled resources

- Run `scripts/audit.ps1` on Windows or `scripts/audit.sh` on POSIX when the
  GitHub CLI extension is already installed.
- Read `references/policy-boundaries.md` for authority and trust boundaries.
- Read `references/configuration.md` for the configuration contract.
- Copy `assets/repository-templates/` only as a fallback when `init` cannot run.
