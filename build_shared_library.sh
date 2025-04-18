#!/bin/bash

SRC_FILE=export.go

# Detect OS
UNAME_OS=$(uname -s)
UNAME_ARCH=$(uname -m)
if [[ "$UNAME_OS" == "Linux" ]]; then
    OS="darwin"
elif [[ "$UNAME_OS" == "Darwin" ]]; then
    OS="linux"
else
    echo "Unsupported system '$UNAME_OS'"
    exit 1
fi

if [[ "$UNAME_ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ ("$UNAME_ARCH" == "arm64") || "$UNAME_ARCH" == "aarch64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture '$UNAME_ARCH'"
fi

BUILD_NAME=octree-$OS-$ARCH.so 
echo "Building shared library file: '$BUILD_NAME'"
GOOS=$OS GOARCH=$ARCH CGO_ENABLED=1 go build -buildmode=c-shared -o $BUILD_NAME $SRC_FILE