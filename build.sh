#!/bin/bash

set -e

# Build an image that does not depend on libc for network -- trying our
# best to make a truly static binary.
CGO_ENABLED=0 go build -tags netgo
