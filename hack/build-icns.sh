#!/usr/bin/env bash

SRC='../assets/icon.svg'
DIST='../assets'

set -euo pipefail

SCRIPT_DIR="$(dirname "$(realpath "$0")")"

ICONSET="$SCRIPT_DIR/$DIST/GoNES.iconset"
ICNS="$SCRIPT_DIR/$DIST/GoNES.icns"

rm -rf "$ICONSET" "$ICNS"
mkdir -p "$ICONSET"

for SIZE in 16 32 64 128 256 512; do (
    DEST="$ICONSET/icon_${SIZE}x${SIZE}.png"
    basename "$DEST"

    inkscape "$SCRIPT_DIR/$SRC" \
      --export-width="$SIZE" \
      --export-type=png \
      --export-filename=- \
    | convert - \
      -strip \
      -background transparent \
      -gravity center \
      -extent "${SIZE}x${SIZE}" \
      "$DEST"

    if [[ "$SIZE" != 16 ]]; then (
      HALF="$(bc <<<"$SIZE/2")"
      HALF_DEST="$ICONSET/icon_${HALF}x${HALF}@2x.png"
      basename "$HALF_DEST"
      cp "$DEST" "$HALF_DEST"
    ) fi
) done

echo Generate "$(basename "$ICNS")"
iconutil \
  --convert icns \
  --output "$ICNS" \
  "$ICONSET"
