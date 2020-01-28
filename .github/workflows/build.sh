#!/bin/sh

set -ex

DIR=jp-${OS}-${REF}
rm -rf "$DIR"
mkdir "$DIR"

cp README.md "$DIR"

cp LICENSE.txt "$DIR"

GOOS=${OS}
[ "$GOOS" = "mac" ] && GOOS=darwin
GOARCH=amd64 GOOS="$GOOS" go build -o "${DIR}/jp" ./jp

tar -czf "${DIR}.tgz" "$DIR"
