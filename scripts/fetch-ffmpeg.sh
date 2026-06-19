#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

version="${FFMPEG_VERSION:-8.1.1}"
tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

mkdir -p assets/ffmpeg/windows-amd64 assets/ffmpeg/linux-amd64 assets/ffmpeg/darwin-arm64

win_zip="${tmpdir}/ffmpeg-win.zip"
win_dir="${tmpdir}/ffmpeg-win"
curl -fsSL -o "$win_zip" "https://github.com/GyanD/codexffmpeg/releases/download/${version}/ffmpeg-${version}-essentials_build.zip"
mkdir -p "$win_dir"
unzip -q "$win_zip" -d "$win_dir"
find "$win_dir" -type f \( -name ffmpeg.exe -o -name ffprobe.exe \) -exec cp -f {} assets/ffmpeg/windows-amd64/ \;

linux_tar="${tmpdir}/ffmpeg-linux.tar.xz"
linux_dir="${tmpdir}/ffmpeg-linux"
curl -fsSL -o "$linux_tar" "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
mkdir -p "$linux_dir"
tar -xf "$linux_tar" -C "$linux_dir"
find "$linux_dir" -maxdepth 2 -type f \( -name ffmpeg -o -name ffprobe \) -exec cp -f {} assets/ffmpeg/linux-amd64/ \;
chmod +x assets/ffmpeg/linux-amd64/ffmpeg assets/ffmpeg/linux-amd64/ffprobe

mac_ffmpeg_zip="${tmpdir}/ffmpeg-mac.zip"
mac_ffprobe_zip="${tmpdir}/ffprobe-mac.zip"
mac_dir="${tmpdir}/ffmpeg-mac"
for i in 1 2 3; do
	if curl -fsSL -o "$mac_ffmpeg_zip" "https://evermeet.cx/ffmpeg/getrelease/ffmpeg/zip" \
		&& curl -fsSL -o "$mac_ffprobe_zip" "https://evermeet.cx/ffmpeg/getrelease/ffprobe/zip"; then
		break
	fi
	if [[ "$i" -eq 3 ]]; then
		echo "failed to download macOS ffmpeg after 3 attempts" >&2
		exit 1
	fi
	sleep $((5 * i))
done
mkdir -p "$mac_dir"
unzip -q -o "$mac_ffmpeg_zip" -d "$mac_dir"
unzip -q -o "$mac_ffprobe_zip" -d "$mac_dir"
find "$mac_dir" -type f \( -name ffmpeg -o -name ffprobe \) -exec cp -f {} assets/ffmpeg/darwin-arm64/ \;
chmod +x assets/ffmpeg/darwin-arm64/ffmpeg assets/ffmpeg/darwin-arm64/ffprobe

echo "ffmpeg ${version} fetched into assets/ffmpeg/"
