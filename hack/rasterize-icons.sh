#!/usr/bin/env bash
set -euo pipefail

cd "$(git rev-parse --show-toplevel)/assets"

rm -rf png
mkdir -p png

for SIZE in 16 32 48 64 128 256 512; do (
    DEST="png/icon_${SIZE}x${SIZE}.png"
    basename "$DEST"

    inkscape icon.svg \
      --export-width="$SIZE" \
      --export-type=png \
      --export-filename=- \
    | magick - \
      -strip \
      -background transparent \
      -gravity center \
      -extent "${SIZE}x${SIZE}" \
      "$DEST"
) done
