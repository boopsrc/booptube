#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

MODE="${1:-bundled}"
STAGING="installer/staging"
BUILD=".build"

version="$(cat VERSION 2>/dev/null || echo dev)"

rm -rf "$STAGING"
mkdir -p "$STAGING/tools"

copy_tools() {
	local os="$1"
	case "$os" in
		windows)
			cp -f assets/ytdlp/windows-amd64/yt-dlp.exe "$STAGING/tools/"
			cp -f assets/ffmpeg/windows-amd64/ffmpeg.exe "$STAGING/tools/"
			cp -f assets/ffmpeg/windows-amd64/ffprobe.exe "$STAGING/tools/"
			cp -f "$BUILD/booptube.exe" "$STAGING/booptube.exe"
			cp -f "$BUILD/booptube-gui.exe" "$STAGING/booptube-gui.exe"
			;;
		linux)
			cp -f assets/ytdlp/linux-amd64/yt-dlp "$STAGING/tools/"
			cp -f assets/ffmpeg/linux-amd64/ffmpeg "$STAGING/tools/"
			cp -f assets/ffmpeg/linux-amd64/ffprobe "$STAGING/tools/"
			chmod +x "$STAGING/tools/"*
			cp -f "$BUILD/booptube" "$STAGING/booptube"
			cp -f "$BUILD/booptube-gui" "$STAGING/booptube-gui"
			chmod +x "$STAGING/booptube" "$STAGING/booptube-gui"
			;;
		macos)
			cp -f assets/ytdlp/darwin-arm64/yt-dlp "$STAGING/tools/"
			cp -f assets/ffmpeg/darwin-arm64/ffmpeg "$STAGING/tools/"
			cp -f assets/ffmpeg/darwin-arm64/ffprobe "$STAGING/tools/"
			chmod +x "$STAGING/tools/"*
			cp -f "$BUILD/booptube" "$STAGING/booptube"
			cp -f "$BUILD/booptube-gui" "$STAGING/booptube-gui"
			chmod +x "$STAGING/booptube" "$STAGING/booptube-gui"
			;;
	esac
}

detect_os() {
	case "$(uname -s)" in
		MINGW*|MSYS*|CYGWIN*|Windows*) echo windows ;;
		Darwin) echo macos ;;
		Linux) echo linux ;;
		*) echo "unsupported"; exit 1 ;;
	esac
}

OS="$(detect_os)"

if [[ "$MODE" == "bundled" ]]; then
	copy_tools "$OS"
	echo "Staged bundled release v$version in $STAGING/"
elif [[ "$MODE" == "portable" ]]; then
	rm -rf "$STAGING/tools"
	case "$OS" in
		windows)
			cp -f "$BUILD/booptube.exe" "$STAGING/booptube.exe"
			cp -f "$BUILD/booptube-gui.exe" "$STAGING/booptube-gui.exe"
			;;
		*)
			cp -f "$BUILD/booptube" "$STAGING/booptube"
			cp -f "$BUILD/booptube-gui" "$STAGING/booptube-gui"
			;;
	esac
	echo "Staged portable exes v$version in $STAGING/"
else
	echo "usage: stage.sh bundled|portable" >&2
	exit 1
fi

if [[ -f doc/usuario.md ]]; then
	cp -f doc/usuario.md "$STAGING/README.md"
fi

echo "$version" > "$STAGING/VERSION.txt"
