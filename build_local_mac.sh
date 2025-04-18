#!/bin/bash

SRC_FILE=export.go

OS=""
BUILD_NAME=""

# Build mac
OS=darwin
for ARCH in "arm64" "amd64"; do
    BUILD_NAME=octree-$OS-$ARCH.so 
    GOOS=$OS GOARCH=$ARCH CGO_ENABLED=1 go build -buildmode=c-shared -o $BUILD_NAME $SRC_FILE
done