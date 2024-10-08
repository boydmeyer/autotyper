#!/bin/bash

# Set up common variables
NAME="autotyper"
PROJECT_PATH=$(pwd)
BIN_DIR="${PROJECT_PATH}/bin"
DIST_DIR="${PROJECT_PATH}/dist"
VERSION=$(git describe --tags --always --dirty)

# Create necessary directories
mkdir -p "$BIN_DIR"
mkdir -p "$DIST_DIR/windows"
mkdir -p "$DIST_DIR/linux"
mkdir -p "$DIST_DIR/macos"

# Build for Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o "${BIN_DIR}/${NAME}-win.exe" -ldflags="-X main.Version=${VERSION}" .
cp "${BIN_DIR}/${NAME}-win.exe" "$DIST_DIR/windows/${NAME}-win.exe"

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o "${BIN_DIR}/${NAME}-linux" -ldflags="-X main.Version=${VERSION}" .
cp "${BIN_DIR}/${NAME}-linux" "$DIST_DIR/linux/${NAME}-linux"

# Build for macOS
echo "Building for macOS..."
export PATH=/mnt/c/osxcross/target/bin:$PATH
export CC=o64-clang
export CXX=o64-clang++
export CGO_ENABLED=1
export GOOS=darwin
export GOARCH=amd64  # or arm64 for Apple Silicon
export MACOSX_DEPLOYMENT_TARGET=10.14

go build -o "${BIN_DIR}/${NAME}-mac" -ldflags="-X main.Version=${VERSION}" .
cp "${BIN_DIR}/${NAME}-mac" "$DIST_DIR/macos/${NAME}-mac"

echo "Build and packaging complete."