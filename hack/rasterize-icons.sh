#!/usr/bin/env bash
set -euo pipefail

ICONSET=GoNES.iconset

cd "$(git rev-parse --show-toplevel)/assets"

rm -rf "$ICONSET"
mkdir -p "$ICONSET"

for SIZE in 16 32 64 128 256 512; do (
    DEST="$ICONSET/icon_${SIZE}x${SIZE}.png"
    basename "$DEST"

    inkscape icon.svg \
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
