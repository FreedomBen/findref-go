#!/usr/bin/env bash

set -e

VERSION=0.0.1

# See: https://stackoverflow.com/a/30068222/2062384 for list of valid targets
OSES=("linux" "windows" "darwin")
ARCHES=("amd64" "386")
#ARCHES=("amd64" "386" "arm64")

echo "Building findref version '${VERSION}' for OSes '${OSES[@]}', and arches '${ARCHES[@]}'..."

root_dir=releases/${VERSION}
for os in ${OSES[@]}; do
    for arch in ${ARCHES[@]}; do
        echo "Building version ${VERSION} for OS ${os}, arch ${arch}"
        docker run \
          --rm \
          --volume "$PWD":/usr/src/myapp \
          --workdir /usr/src/myapp \
          --env GOOS=${os} \
          --env GOARCH=${arch} \
          golang:1.9-alpine go build -v
        mkdir -p ${root_dir}/${os}/${arch}
        mv --force findref ${root_dir}/${os}/${arch}/
    done
done

echo "Done!"
