# Security Policy

## Supported versions

Security fixes are applied to the latest release. Pin a reviewed release tag in
automated consumers and update deliberately.

## Report a vulnerability

Use GitHub private vulnerability reporting for this repository. Do not include
tokens, credentials, private repository content, or exploit details in a public
issue.

## Trust boundaries

The extension executes local commands explicitly configured by the repository.
Review `.github/pr-quality-gate.yml` before running `validate --run-local`,
especially on an untrusted branch.

The supplied PR-policy workflow checks out the trusted base commit and does not
execute code from the pull request. The supplied repository-quality workflow
does execute candidate repository commands, so it uses read-only permissions
and must not receive secrets from untrusted pull requests.

Branch-protection writes require `--apply` and an exact confirmation value. The
tool cannot merge, force-push, change Draft state, or close issues and pull
requests.
