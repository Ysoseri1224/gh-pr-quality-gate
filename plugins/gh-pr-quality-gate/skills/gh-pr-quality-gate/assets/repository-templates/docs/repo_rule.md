# Repository Operation Rules

Read the repository's current product, API, architecture, development, security,
and testing sources before changing code. Preserve unrelated work and secrets.

Every pull request must reference an issue. Use `Closes #N` only when the pull
request fully completes the issue; otherwise use `Refs #N`. Run the configured
local quality gate before pushing and record exact test evidence. Update the
authoritative documentation in the same pull request when a public contract
changes.

Required CI and branch protection are the enforcement boundary. A pre-push hook
and agent instructions provide earlier feedback but are not sufficient controls.

Agents must not merge, force-push, change Draft state, close or reopen work
items, publish releases, deploy, or alter repository rules without explicit
authorization for that exact action.
