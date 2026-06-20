param(
    [ValidateSet("auto", "windows")]
    [string]$Target = "auto"
)

$ErrorActionPreference = "Stop"
Set-Location (Resolve-Path (Join-Path $PSScriptRoot ".."))

$Staging = "installer/staging"
$OutDir = ".build"
$Version = if (Test-Path VERSION) { (Get-Content VERSION -Raw).Trim() } else { "dev" }

if (-not (Test-Path $Staging)) {
    throw "Run scripts/stage.ps1 after build-bundled first"
}

New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

$Iscc = $null
foreach ($cmd in @("ISCC.exe", "iscc")) {
    $found = Get-Command $cmd -ErrorAction SilentlyContinue
    if ($found) { $Iscc = $found.Source; break }
}
if (-not $Iscc) {
    $candidates = @(
        "$env:LOCALAPPDATA\Programs\Inno Setup 6\ISCC.exe",
        "${env:ProgramFiles(x86)}\Inno Setup 6\ISCC.exe",
        "$env:ProgramFiles\Inno Setup 6\ISCC.exe"
    )
    foreach ($path in $candidates) {
        if (Test-Path $path) { $Iscc = $path; break }
    }
}

if (-not $Iscc) {
    Write-Error "Inno Setup not found. Install: winget install JRSoftware.InnoSetup"
}

& $Iscc "installer/windows/booptube.iss" "/DAppVersion=$Version"
Write-Host "Installer artifacts in $OutDir/"
