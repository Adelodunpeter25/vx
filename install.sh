#!/bin/bash

set -e

# VX Editor Installation Script

REPO="Adelodunpeter25/vx"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="vx"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$OS" in
    linux)
        OS="linux"
        ;;
    darwin)
        OS="darwin"
        ;;
    *)
        echo "Unsupported operating system: $OS"
        exit 1
        ;;
esac

case "$ARCH" in
    x86_64|amd64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

BINARY="vx-${OS}-${ARCH}"

echo "Installing VX Editor..."
echo "OS: $OS"
echo "Architecture: $ARCH"
echo ""

# Get latest release URL
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep "browser_download_url.*$BINARY\"" | cut -d '"' -f 4)

if [ -z "$LATEST_RELEASE" ]; then
    echo "Error: Could not find release for $BINARY"
    echo "Please visit https://github.com/$REPO/releases to download manually"
    exit 1
fi

echo "Downloading $BINARY..."
curl -L -o "/tmp/$BINARY" "$LATEST_RELEASE"

echo "Installing to $INSTALL_DIR/$BINARY_NAME..."
chmod +x "/tmp/$BINARY"

# Check if we need sudo
if [ -w "$INSTALL_DIR" ]; then
    mv "/tmp/$BINARY" "$INSTALL_DIR/$BINARY_NAME"
else
    sudo mv "/tmp/$BINARY" "$INSTALL_DIR/$BINARY_NAME"
fi

echo ""
echo "✓ VX Editor installed successfully!"
echo ""

# Check for clipboard support on Linux
if [ "$OS" = "linux" ]; then
    if ! command -v xclip &> /dev/null && ! command -v xsel &> /dev/null; then
        echo "⚠ Clipboard support requires xclip or xsel"
        echo ""
        read -p "Install xclip for clipboard support? [Y/n] " -n 1 -r
        echo ""
        
        if [[ ! $REPLY =~ ^[Nn]$ ]]; then
            # Detect package manager and install
            if command -v apt-get &> /dev/null; then
                echo "Installing xclip via apt..."
                sudo apt-get update && sudo apt-get install -y xclip
            elif command -v yum &> /dev/null; then
                echo "Installing xclip via yum..."
                sudo yum install -y xclip
            elif command -v dnf &> /dev/null; then
                echo "Installing xclip via dnf..."
                sudo dnf install -y xclip
            elif command -v pacman &> /dev/null; then
                echo "Installing xclip via pacman..."
                sudo pacman -S --noconfirm xclip
            elif command -v zypper &> /dev/null; then
                echo "Installing xclip via zypper..."
                sudo zypper install -y xclip
            else
                echo "Could not detect package manager. Please install xclip manually:"
                echo "  Ubuntu/Debian: sudo apt-get install xclip"
                echo "  Fedora/RHEL:   sudo dnf install xclip"
                echo "  Arch:          sudo pacman -S xclip"
            fi
            echo ""
        fi
    fi
fi

echo "Run 'vx --help' to get started"
echo "Run 'vx filename' to edit a file"
