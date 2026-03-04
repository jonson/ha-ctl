#!/usr/bin/env bash
set -euo pipefail

REPO="jonson/ha-ctl"
INSTALL_DIR="${HOME}/.openclaw/skills/ha-ctl"

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

if [[ "$OS" != "linux" ]]; then
    echo "Unsupported OS: $OS (only linux is supported)"
    exit 1
fi

BINARY="ha-ctl-${OS}-${ARCH}"

echo "Installing ha-ctl for ${OS}/${ARCH}..."

# Create install directory
mkdir -p "$INSTALL_DIR"

# Download binary (from GitHub releases or local path)
if [[ -n "${HA_CTL_LOCAL_PATH:-}" ]]; then
    cp "$HA_CTL_LOCAL_PATH" "$INSTALL_DIR/ha-ctl"
else
    LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep tag_name | cut -d'"' -f4)
    curl -fsSL "https://github.com/${REPO}/releases/download/${LATEST}/${BINARY}" -o "$INSTALL_DIR/ha-ctl"
fi

chmod +x "$INSTALL_DIR/ha-ctl"

# Copy SKILL.md
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
if [[ -f "$SCRIPT_DIR/SKILL.md" ]]; then
    cp "$SCRIPT_DIR/SKILL.md" "$INSTALL_DIR/SKILL.md"
fi

echo "ha-ctl installed to $INSTALL_DIR/ha-ctl"
echo "Make sure $INSTALL_DIR is in your PATH or reference the full path."
