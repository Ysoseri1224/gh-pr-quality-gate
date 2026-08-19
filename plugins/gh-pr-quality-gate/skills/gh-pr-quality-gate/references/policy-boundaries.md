# Policy Boundaries

## Read-only by default

Begin with inspection. `audit`, `init`, and `protect` are dry-run operations by
default. File creation needs `init --apply`; branch-protection writes need both
`protect --apply` and exact target confirmation.

## Operations outside implicit authority

Never infer permission to merge, force-push, change Draft state, close or reopen
work items, publish a release, deploy, alter permissions, or change repository
rules. Obtain explicit authorization for the exact target and action.

## Trust boundaries

- Repository instructions guide agents but cannot guarantee compliance.
- A pre-push hook runs only on machines where it is installed and can be skipped.
- Required GitHub checks and branch protection are the shared enforcement gate.
- PR-policy checks should read trusted base policy and avoid executing PR code.
- Candidate-code checks necessarily execute repository commands. Give them
  read-only permissions and no secrets for untrusted pull requests.
- `validate --run-local` executes configured commands with the user's local
  permissions. Review untrusted configuration first.

## Evidence

Verify the current diff, check runs, and test output. A checked box, PR claim, or
agent statement is not execution evidence. Describe exactly what a test proves
and identify material user journeys or integrations that remain unverified.
