# Releasing

1. Confirm `main` is clean and CI has passed.
2. Update matching versions in both plugin manifests and both marketplace files.
3. Run the full local validation documented in `CONTRIBUTING.md`.
4. Create and push an annotated `vX.Y.Z` tag only with maintainer authorization.
5. The release workflow builds and attests precompiled extension binaries.
6. Verify the GitHub release, then install it in a clean environment with
   `gh extension install Ysoseri1224/gh-pr-quality-gate`.
7. Validate Codex and Claude marketplace installation against the release.

Publishing the GitHub release does not publish the plugin to OpenAI's public
directory. Follow `docs/openai-submission.md` for that separate process.
