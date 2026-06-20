#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")/../.."

VERSION="${1:-$(cat VERSION 2>/dev/null || echo dev)}"
STAGING="installer/staging"
OUT_DIR=".build"
APP="booptube-gui.app"
DMG="booptube-${VERSION}-macos-arm64-setup.dmg"
WORK="installer/macos/work"

[[ -d "$STAGING/booptube-gui" ]] || { echo "run make stage first" >&2; exit 1; }

rm -rf "$WORK"
mkdir -p "$WORK/$APP/Contents/MacOS" "$WORK/$APP/Contents/Resources/tools" "$OUT_DIR"

cp -f "$STAGING/booptube-gui" "$WORK/$APP/Contents/MacOS/booptube-gui"
cp -f "$STAGING/booptube" "$WORK/$APP/Contents/Resources/booptube"
cp -f "$STAGING/tools/"* "$WORK/$APP/Contents/Resources/tools/"
chmod +x "$WORK/$APP/Contents/MacOS/booptube-gui" "$WORK/$APP/Contents/Resources/booptube" "$WORK/$APP/Contents/Resources/tools/"*

sed "s/VERSION_PLACEHOLDER/$VERSION/g" installer/macos/Info.plist > "$WORK/$APP/Contents/Info.plist"
mkdir -p "$WORK/$APP/Contents/Resources"
if [[ -f "$STAGING/README.md" ]]; then
	cp -f "$STAGING/README.md" "$WORK/$APP/Contents/Resources/README.txt"
fi

ln -sf /Applications "$WORK/Applications"
cp -f doc/usuario.md "$WORK/README.txt" 2>/dev/null || true

rm -f "$OUT_DIR/$DMG"
hdiutil create -volname "booptube $VERSION" -srcfolder "$WORK" -ov -format UDZO "$OUT_DIR/$DMG"

echo "Created $OUT_DIR/$DMG"
echo "If Gatekeeper blocks: xattr -cr /Applications/booptube-gui.app"

rm -rf "$WORK"
