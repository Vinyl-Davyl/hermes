#!/bin/sh
# Install the hermes binary to ~/.local/bin (override with HERMES_INSTALL_DIR).
# Usage:  curl -fsSL https://tryhermes.pages.dev/install | sh
set -eu

REPO="Vinyl-Davyl/hermes"
DEST="${HERMES_INSTALL_DIR:-$HOME/.local/bin}"
BIN="$DEST/hermes"

os=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
case "$arch" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) arch="" ;;
esac

mkdir -p "$DEST"

script_dir=""
if [ -n "${0:-}" ] && [ "$0" != "sh" ] && [ "$0" != "bash" ] && [ "$0" != "-" ]; then
  script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd) || script_dir=""
fi

if [ -n "$script_dir" ] && [ -f "$script_dir/../cmd/hermes/main.go" ]; then
  if ! command -v go >/dev/null 2>&1; then
    echo "hermes: Go is required to build from this checkout" >&2
    echo "https://go.dev/dl/" >&2
    exit 1
  fi
  echo "building hermes → $BIN"
  (cd "$script_dir/.." && go build -o "$BIN" ./cmd/hermes)
else
  url="https://github.com/${REPO}/releases/latest/download/hermes_${os}_${arch}"
  tmp=$(mktemp)
  if [ -n "$arch" ] && [ "$os" = "darwin" -o "$os" = "linux" ] && curl -fsSL "$url" -o "$tmp"; then
    mv "$tmp" "$BIN"
  else
    rm -f "$tmp"
    if ! command -v go >/dev/null 2>&1; then
      echo "hermes: no release for ${os}/${arch:-unknown}, and Go is not installed" >&2
      echo "install Go from https://go.dev/dl/ then re-run" >&2
      exit 1
    fi
    echo "installing hermes"
    GOBIN="$DEST" go install "github.com/${REPO}/cmd/hermes@latest"
  fi
fi

chmod 0755 "$BIN"
echo "installed $BIN"
"$BIN" version

case ":$PATH:" in
  *":$DEST:"*) ;;
  *) echo "note: $DEST is not on your PATH"
     echo "      export PATH=\"$DEST:\$PATH\"" ;;
esac

echo "try:  hermes doctor"
