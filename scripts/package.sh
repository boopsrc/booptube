#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/.."

TARGET="${1:-auto}"
OUT_DIR=".build"
STAGING="installer/staging"
VERSION="$(cat VERSION 2>/dev/null || echo dev)"

mkdir -p "$OUT_DIR"

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

[[ -d "$STAGING" ]] || { echo "run make stage first" >&2; exit 1; }

case "$TARGET" in
	windows)
		if command -v ISCC.exe >/dev/null 2>&1; then
			ISCC.exe installer/windows/booptube.iss "/DAppVersion=$VERSION"
		elif command -v iscc >/dev/null 2>&1; then
			iscc installer/windows/booptube.iss "/DAppVersion=$VERSION"
		else
			echo "Inno Setup not found. Install: winget install JRSoftware.InnoSetup" >&2
			exit 1
		fi
		;;
	linux)
		if command -v nfpm >/dev/null 2>&1; then
			VERSION="$VERSION" nfpm pkg --config installer/linux/nfpm.yaml --packager deb --target "$OUT_DIR"
			VERSION="$VERSION" nfpm pkg --config installer/linux/nfpm.yaml --packager rpm --target "$OUT_DIR"
		else
			echo "nfpm not found; creating tarball installer instead"
			tar -czf "$OUT_DIR/booptube-${VERSION}-linux-amd64-bundled.tar.gz" -C "$STAGING" .
			echo "Run: sudo bash installer/linux/install.sh from extracted dir"
		fi
		;;
	macos)
		bash installer/macos/build-dmg.sh "$VERSION"
		;;
	*)
		echo "unknown target: $TARGET" >&2
		exit 1
		;;
esac

echo "Installer artifacts in $OUT_DIR/"
