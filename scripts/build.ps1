param(
    [ValidateSet("cli", "gui", "all")]
    [string]$Target = "all"
)

$ErrorActionPreference = "Stop"
Set-Location (Resolve-Path (Join-Path $PSScriptRoot ".."))

$BuildDir = ".build"
New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null

$VersionFile = Join-Path $PWD "VERSION"
$Version = if (Test-Path $VersionFile) { (Get-Content $VersionFile -Raw).Trim() } else { "dev" }

$Commit = "none"
try {
    $Commit = (git rev-parse --short HEAD 2>$null)
    if (-not $Commit) { $Commit = "none" }
} catch {
    $Commit = "none"
}

$BuildDate = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

$LdFlags = @(
    "-s", "-w",
    "-X", "booptube/buildinfo.Version=$Version",
    "-X", "booptube/buildinfo.Commit=$Commit",
    "-X", "booptube/buildinfo.BuildDate=$BuildDate"
)

function Build-Cli {
    Write-Host "Building booptube (CLI) v$Version..."
    go build -trimpath -ldflags ($LdFlags -join " ") -o (Join-Path $BuildDir "booptube.exe") ./cmd/cli
}

function Build-Gui {
    Write-Host "Building booptube-gui v$Version..."
    $GuiLdFlags = $LdFlags + @("-H=windowsgui")
    $env:CGO_ENABLED = "1"
    go build -trimpath -ldflags ($GuiLdFlags -join " ") -o (Join-Path $BuildDir "booptube-gui.exe") ./cmd/gui
}

switch ($Target) {
    "cli" { Build-Cli }
    "gui" { Build-Gui }
    "all" {
        Build-Cli
        Build-Gui
    }
}

Write-Host "Done. Binaries in $BuildDir/"
