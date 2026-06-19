$ErrorActionPreference = "Stop"
Set-Location (Resolve-Path (Join-Path $PSScriptRoot ".."))

$Version = if ($env:FFMPEG_VERSION) { $env:FFMPEG_VERSION } else { "8.1.1" }

$destDirs = @(
    "assets/ffmpeg/windows-amd64",
    "assets/ffmpeg/linux-amd64",
    "assets/ffmpeg/darwin-arm64"
)
foreach ($d in $destDirs) {
    New-Item -ItemType Directory -Force -Path $d | Out-Null
}

function Copy-Binaries($files, $destDir) {
    foreach ($f in $files) {
        $name = Split-Path $f -Leaf
        Copy-Item -Force $f (Join-Path $destDir $name)
    }
}

$winZip = Join-Path $env:TEMP "booptube-ffmpeg-win.zip"
$winDir = Join-Path $env:TEMP "booptube-ffmpeg-win"
$winUrl = "https://github.com/GyanD/codexffmpeg/releases/download/$Version/ffmpeg-$Version-essentials_build.zip"
Invoke-WebRequest -Uri $winUrl -OutFile $winZip
if (Test-Path $winDir) { Remove-Item -Recurse -Force $winDir }
Expand-Archive -Path $winZip -DestinationPath $winDir -Force
$winBin = Get-ChildItem -Path $winDir -Recurse -File | Where-Object {
    $_.Name -in @("ffmpeg.exe", "ffprobe.exe")
}
Copy-Binaries $winBin.FullName "assets/ffmpeg/windows-amd64"
Remove-Item -Force $winZip
Remove-Item -Recurse -Force $winDir

$linuxTar = Join-Path $env:TEMP "booptube-ffmpeg-linux.tar.xz"
$linuxDir = Join-Path $env:TEMP "booptube-ffmpeg-linux"
$linuxUrl = "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
Invoke-WebRequest -Uri $linuxUrl -OutFile $linuxTar
if (Test-Path $linuxDir) { Remove-Item -Recurse -Force $linuxDir }
New-Item -ItemType Directory -Force -Path $linuxDir | Out-Null
tar -xf $linuxTar -C $linuxDir
$linuxBin = Get-ChildItem -Path $linuxDir -Recurse -File | Where-Object {
    $_.Name -in @("ffmpeg", "ffprobe") -and $_.DirectoryName -notmatch "model"
}
Copy-Binaries $linuxBin.FullName "assets/ffmpeg/linux-amd64"
Remove-Item -Force $linuxTar
Remove-Item -Recurse -Force $linuxDir

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

Write-Host "ffmpeg $Version fetched into assets/ffmpeg/"
