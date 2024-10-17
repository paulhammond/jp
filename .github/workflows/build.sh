#!/bin/sh

set -ex

BUILD="build/jp-${REF}"
mkdir -p "$BUILD"

cp README.md "$BUILD"
cp LICENSE "$BUILD"

GOOS=${OS}
[ "$GOOS" = "macos" ] && GOOS=darwin
GOARCH="${ARCH}" GOOS="${GOOS}" go build -o "$BUILD/jp" ./jp

cd build
tar -czf "jp-${OS}-${ARCH}-${REF}.tgz" "jp-${REF}"
