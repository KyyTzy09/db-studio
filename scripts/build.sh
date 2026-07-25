#!/usr/bin/env bash
set -e

echo "=========================================="
echo " 🚀 Starting DBStudio Full Build Pipeline "
echo "=========================================="

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"

echo ""
echo "[1/3] Building SvelteKit Frontend Static Assets..."
cd "$REPO_ROOT/web"
bun run build
echo "✅ Frontend build completed successfully."

echo ""
echo "[2/3] Compiling Go Single Binaries (CGO_ENABLED=0)..."
cd "$REPO_ROOT"
mkdir -p bin

export CGO_ENABLED=0

echo " -> Building Windows x64 (bin/dbstudio-win-x64.exe)..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/dbstudio-win-x64.exe main.go

echo " -> Building Linux x64 (bin/dbstudio-linux-x64)..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/dbstudio-linux-x64 main.go

echo " -> Building macOS ARM64 (bin/dbstudio-darwin-arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/dbstudio-darwin-arm64 main.go

echo " -> Building macOS x64 (bin/dbstudio-darwin-x64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/dbstudio-darwin-x64 main.go

chmod +x bin/dbstudio-linux-x64 bin/dbstudio-darwin-arm64 bin/dbstudio-darwin-x64

echo ""
echo "[3/3] Build Summary:"
ls -lh bin/

echo ""
echo "🎉 All DBStudio binaries compiled successfully in bin/!"
