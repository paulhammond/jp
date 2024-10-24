#!/bin/sh

set -ex

BUILD="build/jp-${REF}"
TGZ="jp-${OS}-${ARCH}-${REF}.tgz"
mkdir -p "$BUILD"

cp README.md "$BUILD"
cp LICENSE "$BUILD"

GOOS=${OS}
[ "$GOOS" = "macos" ] && GOOS=darwin
GOARCH="${ARCH}" GOOS="${GOOS}" go build -o "$BUILD/jp" ./cmd/jp

cd build
tar -czf "${TGZ}" "jp-${REF}"

echo "::notice title=sha256::$(shasum -a 256 "${TGZ}")"
