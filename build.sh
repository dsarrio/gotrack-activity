#!/bin/sh

set -euo pipefail
#set -x

rm -rf gotrack.app

cp -R gotrack.app.tpl gotrack.app

go build -o gotrack.app/Contents/MacOS/
