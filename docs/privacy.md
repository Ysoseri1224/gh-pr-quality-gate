# Privacy Policy

Effective date: 2026-08-19

`gh-pr-quality-gate` is local, open-source developer tooling. The Skill and CLI
do not operate a hosted service, use analytics, place cookies, or send repository
content to the maintainer.

Commands may invoke software already configured by the user:

- `gh` communicates with GitHub under the user's existing authentication;
- repository-defined local gate commands may use their own tools and services;
- Codex, ChatGPT, Claude, Gemini, or another host processes prompts and files
  under that host's terms and privacy policy.

The CLI reads repository policy, pull request body files supplied by the user,
and local file existence. `protect --apply` sends the displayed branch
protection payload to GitHub. No credentials are stored by this project.

Do not place credentials or private repository content in public issues. Report
security concerns using GitHub private vulnerability reporting.
