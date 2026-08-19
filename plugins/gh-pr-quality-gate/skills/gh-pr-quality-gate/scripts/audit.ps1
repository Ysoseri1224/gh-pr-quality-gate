[CmdletBinding()]
param(
    [string]$Repository = "."
)

$ErrorActionPreference = "Stop"
gh pr-quality-gate audit --repo $Repository
