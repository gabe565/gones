#!/usr/bin/env bash

BINARY_NAME='gones'
APP_NAME='GoNES'

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

rm -rf dist assets/GoNES.iconset
mkdir -p dist

# Generate icns
cp -a assets/{png,GoNES.iconset}
cp -a assets/GoNES.iconset/icon_{32x32,16x16@2x}.png
rm assets/GoNES.iconset/icon_48x48.png
cp -a assets/GoNES.iconset/icon_{64x64,32x32@2x}.png
cp -a assets/GoNES.iconset/icon_{128x128,64x64@2x}.png
cp -a assets/GoNES.iconset/icon_{256x256,128x128@2x}.png
cp -a assets/GoNES.iconset/icon_{512x512,256x256@2x}.png
iconutil --convert icns --output assets/GoNES.icns assets/GoNES.iconset

# Build binaries
export GOOS=darwin CGO_ENABLED=1
for ARCH in amd64 arm64; do
  echo Build "$BINARY_NAME-$ARCH"
  GOARCH="$ARCH" go build -ldflags='-w -s' -trimpath -o "dist/$BINARY_NAME-$ARCH" "$(git rev-parse --show-toplevel)"
done

# Merge binaries
lipo -create -output "dist/$BINARY_NAME" "dist/$BINARY_NAME-amd64" "dist/$BINARY_NAME-arm64"
rm "dist/$BINARY_NAME-amd64" "dist/$BINARY_NAME-arm64"
echo ...done

# Generate app
echo Generate "$APP_NAME.app"
APP_CONTENTS="dist/$APP_NAME.app/Contents"
mkdir -p "$APP_CONTENTS"
cp "assets/info.plist" "$APP_CONTENTS"
mkdir "$APP_CONTENTS/Resources"
cp "assets/GoNES.icns" "$APP_CONTENTS/Resources"
mkdir "$APP_CONTENTS/MacOS"
cp -p "dist/$BINARY_NAME" "$APP_CONTENTS/MacOS"
rm "dist/$BINARY_NAME"
echo ...done
