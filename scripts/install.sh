#!/bin/sh
# Build Hermes from this repo and put it on your PATH.
# Default: ~/.local/bin/hermes  (override with HERMES_INSTALL_DIR)
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
DIR="${HERMES_INSTALL_DIR:-$HOME/.local/bin}"

if ! command -v go >/dev/null 2>&1; then
  echo "hermes: Go is required to build from source" >&2
  echo "install Go, then re-run, or: go build -o bin/hermes ./cmd/hermes" >&2
  exit 1
fi

mkdir -p "$DIR"
echo "building hermes → $DIR/hermes"
(cd "$ROOT" && go build -o "$DIR/hermes" ./cmd/hermes)
chmod 0755 "$DIR/hermes"

echo "installed $DIR/hermes"
"$DIR/hermes" version

case ":$PATH:" in
  *":$DIR:"*) ;;
  *) echo "note: $DIR is not on your PATH"
     echo "      export PATH=\"$DIR:\$PATH\"" ;;
esac

echo "try:  hermes doctor"
echo "      hermes handoff --from none -m \"try hermes\""
