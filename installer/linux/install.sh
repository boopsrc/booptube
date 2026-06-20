#!/usr/bin/env bash
set -euo pipefail

# Install booptube from bundled tarball (no nfpm).
# Usage: sudo ./install.sh [install-dir]

INSTALL_DIR="${1:-/opt/booptube}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

if [[ "$(id -u)" -ne 0 ]]; then
	echo "Run as root: sudo $0" >&2
	exit 1
fi

SRC="$SCRIPT_DIR"
if [[ -f "$SCRIPT_DIR/../booptube" ]]; then
	SRC="$(cd "$SCRIPT_DIR/.." && pwd)"
fi

mkdir -p "$INSTALL_DIR/tools"
install -m 755 "$SRC/booptube" "$INSTALL_DIR/booptube"
install -m 755 "$SRC/booptube-gui" "$INSTALL_DIR/booptube-gui"
install -m 755 "$SRC/tools/"* "$INSTALL_DIR/tools/"

ln -sf "$INSTALL_DIR/booptube" /usr/local/bin/booptube
ln -sf "$INSTALL_DIR/booptube-gui" /usr/local/bin/booptube-gui

if [[ -f "$SCRIPT_DIR/booptube-gui.desktop" ]]; then
	install -m 644 "$SCRIPT_DIR/booptube-gui.desktop" /usr/share/applications/booptube-gui.desktop
	sed -i "s|Exec=.*|Exec=$INSTALL_DIR/booptube-gui|" /usr/share/applications/booptube-gui.desktop
fi

echo "Installed to $INSTALL_DIR"
