[CmdletBinding()]
param(
    [ValidateSet("Codex", "Claude", "Gemini", "All")]
    [string]$Target = "All"
)

$ErrorActionPreference = "Stop"
$repoRoot = Split-Path -Parent $PSScriptRoot
$source = Join-Path $repoRoot "plugins/gh-pr-quality-gate/skills/gh-pr-quality-gate"
$userProfilePath = [Environment]::GetFolderPath("UserProfile")

$destinations = @{
    Codex  = Join-Path $userProfilePath ".codex/skills/gh-pr-quality-gate"
    Claude = Join-Path $userProfilePath ".claude/skills/gh-pr-quality-gate"
    Gemini = Join-Path $userProfilePath ".gemini/skills/gh-pr-quality-gate"
}

$selected = if ($Target -eq "All") { @("Codex", "Claude", "Gemini") } else { @($Target) }
foreach ($platform in $selected) {
    $destination = $destinations[$platform]
    if (Test-Path -LiteralPath $destination) {
        throw "Refusing to replace existing installation: $destination"
    }
}

foreach ($platform in $selected) {
    $destination = $destinations[$platform]
    New-Item -ItemType Directory -Path (Split-Path -Parent $destination) -Force | Out-Null
    Copy-Item -Recurse -LiteralPath $source -Destination $destination
    Write-Host "Installed gh-pr-quality-gate for $platform at $destination"
}
