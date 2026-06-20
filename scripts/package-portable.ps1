param(
    [ValidateSet("auto", "windows")]
    [string]$Target = "auto"
)

$ErrorActionPreference = "Stop"
Set-Location (Resolve-Path (Join-Path $PSScriptRoot ".."))

$Build = ".build"
$Version = if (Test-Path VERSION) { (Get-Content VERSION -Raw).Trim() } else { "dev" }

New-Item -ItemType Directory -Force -Path $Build | Out-Null

if (-not (Test-Path "$Build/booptube.exe")) {
    throw "Run scripts/build.ps1 first"
}

$name = "booptube-$Version-windows-amd64-portable.zip"
$zipPath = Join-Path $Build $name

if (Test-Path $zipPath) { Remove-Item -Force $zipPath }

Compress-Archive -Path "$Build/booptube.exe", "$Build/booptube-gui.exe" -DestinationPath $zipPath -Force
Write-Host "Created $zipPath"
