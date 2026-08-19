# Configuration Contract

Use `.github/pr-quality-gate.yml` with `version: 1`.

- `local_gate.windows` and `local_gate.posix`: ordered shell commands executed
  from the repository root.
- `required_checks`: exact GitHub check job names to require.
- `required_files`: repository-relative policy and workflow files.
- `issue_reference.required`: require a PR body issue reference.
- `issue_reference.allow_partial_refs`: accept `Refs #N` in addition to closing
  keywords.
- `agent_authority.merge`: must be `explicit-authorisation`.
- `agent_authority.draft_transition`: must be `explicit-authorisation`.
- `agent_authority.force_push`: must be `prohibited`.
- `branch_protection`: target branch, approval count, stale-review handling,
  last-push approval, conversation resolution, and administrator enforcement.

Do not apply branch protection until all configured required checks have
reported successfully on the target repository. Preserve stable check names.
The initial release creates protection only on an unprotected branch and
refuses to replace existing protection that may contain unmanaged restrictions.
