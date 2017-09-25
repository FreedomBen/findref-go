#!/usr/bin/env sh

set -e

VERSION=0.0.1

# See: https://stackoverflow.com/a/30068222/2062384 for list of valid targets
OSES=("linux" "windows" "darwin")
ARCHES=("amd64" "386" "arm" "arm64")

echo "Building findref version ${VERSION} for OSes ${OSES[@]}, and arches ${ARCHES[@]}..."

mkdir -p releases/${VERSION}
cd releases/${VERSION}

orig_dir=$(pwd)
for os in ${OSES[@]}; do
    mkdir -p $os
    for arch in ${ARCHES[@]}; do
        cd $os
        mkdir -p $arch && cd $arch
        echo "Building version ${VERSION} for OS ${os}, arch ${arch}"
        docker run \
          --rm \
          --volume "$PWD":/usr/src/myapp \
          --workdir /usr/src/myapp \
          --env GOOS=${os} \
          --env GOARCH=${arch} \
          golang:1.8 go build -o findref
    done
    cd $orig_dir
done

echo "Done!"
