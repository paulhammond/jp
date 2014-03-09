#!/bin/bash

# To get this script working you need go set up go to do cross compilation.
#  . for mac/homebrew, run "brew install go --cross-compile-common"
#  . on linux install from source then run something like:
#     GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 ./make.bash --no-clean
#     GOOS=linux GOARCH=386 CGO_ENABLED=0 ./make.bash --no-clean

set -eu

VERSION=${1:-"dev"}

for DEST in linux-386 linux-amd64 darwin-amd64; do
	OS=${DEST%-*}
	ARCH=${DEST#*-}
	DIR=pkg/build/$DEST/jp-$VERSION
	mkdir -p $DIR
	cp README.md $DIR
	cp LICENSE.txt $DIR
	GOOS=$OS GOARCH=$ARCH go build -o $DIR/jp github.com/paulhammond/jp/jp
	cd pkg/build/$DEST
	tar -czf ../../jp-${VERSION}-${OS}-${ARCH}.tgz jp-$VERSION
	cd ../../..
done
