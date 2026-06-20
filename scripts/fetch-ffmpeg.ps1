$ErrorActionPreference = "Stop"
Set-Location (Resolve-Path (Join-Path $PSScriptRoot ".."))

$Version = if ($env:FFMPEG_VERSION) { $env:FFMPEG_VERSION } else { "8.1.1" }
$Minor = ($Version -split '\.')[0,1] -join '.'

$destDirs = @(
    "assets/ffmpeg/windows-amd64",
    "assets/ffmpeg/linux-amd64",
    "assets/ffmpeg/darwin-arm64"
)
foreach ($d in $destDirs) {
    New-Item -ItemType Directory -Force -Path $d | Out-Null
}

function Get-FetchPlatforms {
    if ($env:FETCH_FFMPEG_PLATFORMS -eq "all") { return @("windows", "linux", "macos") }
    if ($env:FETCH_FFMPEG_PLATFORMS) { return $env:FETCH_FFMPEG_PLATFORMS.Split(",") }
    if ($IsWindows -or $env:OS -match "Windows") { return @("windows") }
    if ($IsMacOS) { return @("macos") }
    if ($IsLinux) { return @("linux") }
    return @("windows", "linux", "macos")
}

function Copy-Binaries($files, $destDir) {
    foreach ($f in $files) {
        $name = Split-Path $f -Leaf
        Copy-Item -Force $f (Join-Path $destDir $name)
    }
}

function Invoke-Download($urls, $outFile) {
    foreach ($url in $urls) {
        try {
            Invoke-WebRequest -Uri $url -OutFile $outFile -TimeoutSec 600
            if ((Get-Item $outFile).Length -gt 0) { return }
        } catch {
            if (Test-Path $outFile) { Remove-Item -Force $outFile }
        }
    }
    throw "failed to download: $($urls -join ', ')"
}

$platforms = Get-FetchPlatforms

if ($platforms -contains "windows") {
    $winZip = Join-Path $env:TEMP "booptube-ffmpeg-win.zip"
    $winDir = Join-Path $env:TEMP "booptube-ffmpeg-win"
    $winUrls = @(
        "https://www.gyan.dev/ffmpeg/builds/packages/ffmpeg-$Version-essentials_build.zip",
        "https://github.com/GyanD/codexffmpeg/releases/download/$Version/ffmpeg-$Version-essentials_build.zip"
    )
    Invoke-Download $winUrls $winZip
    if (Test-Path $winDir) { Remove-Item -Recurse -Force $winDir }
    Expand-Archive -Path $winZip -DestinationPath $winDir -Force
    $winBin = Get-ChildItem -Path $winDir -Recurse -File | Where-Object {
        $_.Name -in @("ffmpeg.exe", "ffprobe.exe")
    }
    Copy-Binaries $winBin.FullName "assets/ffmpeg/windows-amd64"
    Remove-Item -Force $winZip
    Remove-Item -Recurse -Force $winDir
}

if ($platforms -contains "linux") {
    $linuxTar = Join-Path $env:TEMP "booptube-ffmpeg-linux.tar.xz"
    $linuxDir = Join-Path $env:TEMP "booptube-ffmpeg-linux"
    $linuxUrls = @(
        "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz",
        "https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-n$Minor-latest-linux64-gpl-$Minor.tar.xz"
    )
    Invoke-Download $linuxUrls $linuxTar
    if (Test-Path $linuxDir) { Remove-Item -Recurse -Force $linuxDir }
    New-Item -ItemType Directory -Force -Path $linuxDir | Out-Null
    tar -xf $linuxTar -C $linuxDir
    $linuxBin = Get-ChildItem -Path $linuxDir -Recurse -File | Where-Object {
        $_.Name -in @("ffmpeg", "ffprobe") -and $_.DirectoryName -notmatch "model"
    }
    Copy-Binaries $linuxBin.FullName "assets/ffmpeg/linux-amd64"
    Remove-Item -Force $linuxTar
    Remove-Item -Recurse -Force $linuxDir
}

if ($platforms -contains "macos") {
    $macFfmpegZip = Join-Path $env:TEMP "booptube-ffmpeg-mac-ffmpeg.zip"
    $macFfprobeZip = Join-Path $env:TEMP "booptube-ffmpeg-mac-ffprobe.zip"
    $macDir = Join-Path $env:TEMP "booptube-ffmpeg-mac"
    for ($i = 1; $i -le 3; $i++) {
        try {
            Invoke-WebRequest -Uri "https://evermeet.cx/ffmpeg/getrelease/ffmpeg/zip" -OutFile $macFfmpegZip -TimeoutSec 300
            Invoke-WebRequest -Uri "https://evermeet.cx/ffmpeg/getrelease/ffprobe/zip" -OutFile $macFfprobeZip -TimeoutSec 300
            break
        } catch {
            if ($i -eq 3) { throw }
            Start-Sleep -Seconds (5 * $i)
        }
    }
    if (Test-Path $macDir) { Remove-Item -Recurse -Force $macDir }
    New-Item -ItemType Directory -Force -Path $macDir | Out-Null
    Expand-Archive -Path $macFfmpegZip -DestinationPath $macDir -Force
    Expand-Archive -Path $macFfprobeZip -DestinationPath $macDir -Force
    $macBin = Get-ChildItem -Path $macDir -Recurse -File | Where-Object {
        $_.Name -in @("ffmpeg", "ffprobe")
    }
    Copy-Binaries $macBin.FullName "assets/ffmpeg/darwin-arm64"
    Remove-Item -Force $macFfmpegZip, $macFfprobeZip
    Remove-Item -Recurse -Force $macDir
}

Write-Host "ffmpeg $Version fetched into assets/ffmpeg/"
