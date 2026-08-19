[CmdletBinding()]
param()

$ErrorActionPreference = "Stop"
$repoRoot = git rev-parse --show-toplevel
if (-not $repoRoot) { throw "Run this command inside a Git repository." }

$hookPath = Join-Path $repoRoot ".git/hooks/pre-push"
$hook = @'
#!/bin/sh
# Managed by gh-pr-quality-gate
exec gh pr-quality-gate validate --repo "$(git rev-parse --show-toplevel)" --run-local
'@

if (Test-Path -LiteralPath $hookPath) {
    $current = Get-Content -Raw -LiteralPath $hookPath
    if ($current -ne $hook) {
        throw "A different pre-push hook already exists at $hookPath. Integrate it manually."
    }
    Write-Host "Quality gate pre-push hook is already installed."
    exit 0
}

Set-Content -LiteralPath $hookPath -Value $hook -NoNewline
Write-Host "Installed $hookPath"
