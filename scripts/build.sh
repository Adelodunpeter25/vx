#!/bin/bash

# Build script for vx editor - builds for all platforms locally
# Usage: ./scripts/build.sh

set -e

echo "Building vx for all platforms..."
echo ""

# Create build directory
mkdir -p build

# Linux AMD64
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/vx-linux-amd64 cmd/vx/*.go

# Linux ARM64
echo "Building for Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o build/vx-linux-arm64 cmd/vx/*.go

# macOS Intel
echo "Building for macOS Intel..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/vx-darwin-amd64 cmd/vx/*.go

# macOS Apple Silicon
echo "Building for macOS Apple Silicon..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/vx-darwin-arm64 cmd/vx/*.go

echo ""
echo "✓ Build complete! Binaries are in the build/ directory:"
ls -lh build/
