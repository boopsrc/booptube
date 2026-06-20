#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

TARGET="${1:-auto}"
VERSION="$(cat VERSION 2>/dev/null || echo dev)"
BUILD=".build"

mkdir -p "$BUILD"

detect_os() {
	case "$(uname -s)" in
		MINGW*|MSYS*|CYGWIN*|Windows*) echo windows ;;
		Darwin) echo macos ;;
		Linux) echo linux ;;
		*) echo "unsupported"; exit 1 ;;
	esac
}

if [[ "$TARGET" == "auto" ]]; then
	TARGET="$(detect_os)"
fi

pack_portable() {
	local os="$1"
	local name="booptube-${VERSION}-${os}-amd64-portable"
	case "$os" in
		windows)
			[[ -f "$BUILD/booptube.exe" ]] || { echo "run build first" >&2; exit 1; }
			(
				cd "$BUILD"
				zip -q "${name}.zip" booptube.exe booptube-gui.exe
			)
			;;
		linux|macos)
			[[ -f "$BUILD/booptube" ]] || { echo "run build first" >&2; exit 1; }
			tar -czf "$BUILD/${name}.tar.gz" -C "$BUILD" booptube booptube-gui
			;;
	esac
	echo "Created $BUILD/${name}.*"
}

pack_portable "$TARGET"
