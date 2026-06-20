#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

version="${YTDLP_VERSION:-2026.06.09}"
base="https://github.com/yt-dlp/yt-dlp/releases/download/${version}"

mkdir -p assets/ytdlp/windows-amd64 assets/ytdlp/linux-amd64 assets/ytdlp/darwin-arm64

curl -fsSL -o assets/ytdlp/windows-amd64/yt-dlp.exe "${base}/yt-dlp.exe"
curl -fsSL -o assets/ytdlp/linux-amd64/yt-dlp "${base}/yt-dlp_linux"
curl -fsSL -o assets/ytdlp/darwin-arm64/yt-dlp "${base}/yt-dlp_macos"
chmod +x assets/ytdlp/linux-amd64/yt-dlp assets/ytdlp/darwin-arm64/yt-dlp

echo "yt-dlp ${version} fetched into assets/ytdlp/"
