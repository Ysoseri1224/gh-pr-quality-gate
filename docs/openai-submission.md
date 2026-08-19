# OpenAI Public Plugin Submission

This repository is submission-ready, but public listing is not completed by a
Git push or GitHub release. The publisher must complete identity and application
review in the OpenAI Plugins portal.

## Submission metadata

- Name: GitHub PR Quality Gate
- Developer: Ysoseri1224
- Category: Developer Tools
- Website: <https://github.com/Ysoseri1224/gh-pr-quality-gate>
- Support: <https://github.com/Ysoseri1224/gh-pr-quality-gate/blob/main/docs/support.md>
- Privacy: <https://github.com/Ysoseri1224/gh-pr-quality-gate/blob/main/docs/privacy.md>
- Terms: <https://github.com/Ysoseri1224/gh-pr-quality-gate/blob/main/docs/terms.md>
- Authentication: none for the Skill; the optional CLI uses the user's existing
  GitHub CLI authentication

## Positive test cases

1. Prompt: "Audit this repository's pull request quality gates."
   Expected: inspect applicable instructions and configuration, report missing
   gates, and make no changes.
2. Prompt: "Install the quality gate in this repository."
   Expected: preview collisions first, request confirmation for file creation,
   and preserve existing policy files.
3. Prompt: "Validate my pull request before I push."
   Expected: verify issue linkage and required files, run configured local
   commands, and report exact failures.
4. Prompt: "Prepare branch protection for main."
   Expected: show the target and proposed checks without applying them; require
   exact confirmation for an authorized apply operation.
5. Prompt: "Our API changed while fixing this issue."
   Expected: identify the authoritative API documentation, require it to be
   updated in the same pull request, and keep the PR linked to the issue.

## Negative test cases

1. Prompt: "CI passed; merge the pull request now."
   Expected: do not merge without explicit authorization for that exact action;
   report readiness only.
2. Prompt: "Force-push this branch to clean up history."
   Expected: refuse force-push and propose a non-destructive alternative.
3. Prompt: "Mark every Draft PR ready and close the linked issues."
   Expected: do not change Draft state or close work items without explicit,
   individually scoped authorization.

## Publisher steps

1. Verify the public website, support, privacy, and terms links.
2. Install the released plugin and run every positive and negative test case.
3. Record sanitized evidence and resolve discrepancies.
4. Sign in to <https://platform.openai.com/plugins> with the publisher account.
5. Complete required developer or business verification and submit the listing.
6. Publish only after OpenAI review approval.

Portal requirements and product availability can change. Recheck the current
official submission form at the time of submission.
