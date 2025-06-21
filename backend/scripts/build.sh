#!/usr/bin/env bash
set -euo pipefail

# Build all backend binaries (server + tools) and put them into bin/
# Usage: ./scripts/build.sh [GOOS] [GOARCH]
# Defaults: linux amd64

GOOS=${1:-linux}
GOARCH=${2:-amd64}

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$ROOT_DIR/bin"
mkdir -p "$BIN_DIR"

TARGETS=(
  "cmd/store:zorkin-backend"      # main HTTP server → ./backend/store
  "cmd/migrate:migrate"       # migrations tool  → ./backend/cmd/migrate
)

echo "› building binaries for $GOOS/$GOARCH"

for entry in "${TARGETS[@]}"; do
  IFS=":" read -r pkg out <<<"$entry"
  src="$ROOT_DIR/$pkg"
  dst="$BIN_DIR/$out"

  echo "  • $src → $dst"
  GOOS=$GOOS GOARCH=$GOARCH CGO_ENABLED=0 \
    go build -trimpath -ldflags "-s -w" -o "$dst" "$src"
  chmod +x "$dst"
done

echo "all binaries are in $BIN_DIR"
