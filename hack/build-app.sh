#!/usr/bin/env bash

BINARY_NAME='gones'
APP_NAME='GoNES'

set -euo pipefail

SCRIPT_DIR="$(dirname "$0")"
DIST="$SCRIPT_DIR/../dist"
ASSETS="$SCRIPT_DIR/../assets"

rm -rf "$DIST"
mkdir -p "$DIST"

export GOOS=darwin CGO_ENABLED=1
for ARCH in amd64 arm64; do
  echo Build "$BINARY_NAME-$ARCH"
  GOARCH="$ARCH" go build -ldflags='-w -s' -trimpath -o "$DIST/$BINARY_NAME-$ARCH" "$(git rev-parse --show-toplevel)"
done
lipo -create -output "$DIST/$BINARY_NAME" "$DIST/$BINARY_NAME-amd64" "$DIST/$BINARY_NAME-arm64"
rm "$DIST/$BINARY_NAME-amd64" "$DIST/$BINARY_NAME-arm64"
echo ...done

echo Generate "$APP_NAME.app"
APP_CONTENTS="$DIST/$APP_NAME.app/Contents"
mkdir -p "$APP_CONTENTS"
cp "$ASSETS/info.plist" "$APP_CONTENTS"
mkdir "$APP_CONTENTS/Resources"
cp "$ASSETS/GoNES.icns" "$APP_CONTENTS/Resources"
mkdir "$APP_CONTENTS/MacOS"
cp -p "$DIST/$BINARY_NAME" "$APP_CONTENTS/MacOS"
rm "$DIST/$BINARY_NAME"
echo ...done
