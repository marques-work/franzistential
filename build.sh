#!/bin/bash

set -e

if [ $# -gt 0 ]; then
  if [ "$1" = "--dist" ]; then
    # build for all platforms
    (cd "$(pwd -P)" && go get -d)
    gox -tags netgo -osarch="darwin/amd64 linux/amd64 linux/386 linux/arm linux/arm64" -output="dists/{{.OS}}/{{.Arch}}/franz"
    exit
  fi
fi

# Build an image that does not depend on libc for network -- trying our
# best to make a truly static binary.
CGO_ENABLED=0 go build -tags netgo -o franz
