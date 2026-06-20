$ErrorActionPreference = "Stop"
Set-Location (Resolve-Path (Join-Path $PSScriptRoot ".."))

$Version = if ($env:YTDLP_VERSION) { $env:YTDLP_VERSION } else { "2026.06.09" }
$Base = "https://github.com/yt-dlp/yt-dlp/releases/download/$Version"

$dirs = @(
    "assets/ytdlp/windows-amd64",
    "assets/ytdlp/linux-amd64",
    "assets/ytdlp/darwin-arm64"
)
foreach ($d in $dirs) {
    New-Item -ItemType Directory -Force -Path $d | Out-Null
}

Invoke-WebRequest -Uri "$Base/yt-dlp.exe" -OutFile "assets/ytdlp/windows-amd64/yt-dlp.exe"
Invoke-WebRequest -Uri "$Base/yt-dlp_linux" -OutFile "assets/ytdlp/linux-amd64/yt-dlp"
Invoke-WebRequest -Uri "$Base/yt-dlp_macos" -OutFile "assets/ytdlp/darwin-arm64/yt-dlp"

Write-Host "yt-dlp $Version fetched into assets/ytdlp/"
