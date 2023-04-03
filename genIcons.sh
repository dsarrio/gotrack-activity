#!/bin/sh

set -euo pipefail
#set -x

rm -rf .tmp/app.iconset
mkdir -p .tmp/app.iconset

for SIZE in 16 32 64 128 256 512; do
    sips -z $SIZE $SIZE icon.png --out .tmp/app.iconset/icon_${SIZE}x${SIZE}.png
    sips -z $SIZE $SIZE icon.png --out .tmp/app.iconset/icon_${SIZE}x${SIZE}@2x.png
done
sips -z 1024 1024 icon.png --out .tmp/app.iconset/icon_1024x1024.png

iconutil -c icns -o gotrack.app.tpl/Contents/Resources/icon.icns .tmp/app.iconset
