#!/usr/bin/env bash

BINARY_NAME='GoNES'

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

rm -rf -- dist/*.exe assets/winres/*.png *.syso
mkdir -p dist

command -v go-winres &>/dev/null || go install github.com/tc-hib/go-winres@latest

# Generate metadata
cp -a assets/{png/icon_16x16,winres/icon16}.png
cp -a assets/{png/icon_32x32,winres/icon32}.png
cp -a assets/{png/icon_48x48,winres/icon48}.png
cp -a assets/{png/icon_64x64,winres/icon64}.png
cp -a assets/{png/icon_128x128,winres/icon128}.png
cp -a assets/{png/icon_256x256,winres/icon256}.png
go-winres make --arch=amd64,arm64 --in=assets/winres/winres.json

# Build binary
export GOOS=windows CGO_ENABLED=1
for ARCH in amd64 arm64; do
  echo Build "$BINARY_NAME-$ARCH.exe"
  GOARCH="$ARCH" go build -ldflags='-w -s' -trimpath -o "dist/$BINARY_NAME-$ARCH.exe" "$(git rev-parse --show-toplevel)"
done
