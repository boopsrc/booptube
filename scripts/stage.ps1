param(
    [ValidateSet("bundled", "portable")]
    [string]$Mode = "bundled"
)

$ErrorActionPreference = "Stop"
Set-Location (Resolve-Path (Join-Path $PSScriptRoot ".."))

$Staging = "installer/staging"
$Build = ".build"
$Version = if (Test-Path VERSION) { (Get-Content VERSION -Raw).Trim() } else { "dev" }

Remove-Item -Recurse -Force $Staging -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path "$Staging/tools" | Out-Null

function Copy-Tools {
    Copy-Item -Force "assets/ytdlp/windows-amd64/yt-dlp.exe" "$Staging/tools/"
    Copy-Item -Force "assets/ffmpeg/windows-amd64/ffmpeg.exe" "$Staging/tools/"
    Copy-Item -Force "assets/ffmpeg/windows-amd64/ffprobe.exe" "$Staging/tools/"
}

if ($Mode -eq "bundled") {
    Copy-Tools
    Copy-Item -Force "$Build/booptube.exe" "$Staging/booptube.exe"
    Copy-Item -Force "$Build/booptube-gui.exe" "$Staging/booptube-gui.exe"
    Write-Host "Staged bundled release v$Version in $Staging/"
} else {
    Remove-Item -Recurse -Force "$Staging/tools" -ErrorAction SilentlyContinue
    Copy-Item -Force "$Build/booptube.exe" "$Staging/booptube.exe"
    Copy-Item -Force "$Build/booptube-gui.exe" "$Staging/booptube-gui.exe"
    Write-Host "Staged portable exes v$Version in $Staging/"
}

if (Test-Path "doc/usuario.md") {
    Copy-Item -Force "doc/usuario.md" "$Staging/README.md"
}
Set-Content -Path "$Staging/VERSION.txt" -Value $Version -NoNewline
