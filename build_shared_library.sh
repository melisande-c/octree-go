#!/bin/bash

REPO_DIR=$(pwd)
SRC_FILE=export.go

# Detect OS
UNAME_OS=$(uname -s)
UNAME_ARCH=$(uname -m)
if [[ "$UNAME_OS" == "Linux" ]]; then
    OS="linux"
elif [[ "$UNAME_OS" == "Darwin" ]]; then
    OS="darwin"
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

if [[ "$UNAME_OS" == "Linux" ]]; then
    echo Getting glibc version 2.28

    mkdir -p $HOME/glibc-2.28/build
    cd $HOME/glibc-2.28

    wget http://ftp.gnu.org/gnu/libc/glibc-2.28.tar.gz
    tar -xvzf glibc-2.28.tar.gz
    cd glibc-2.28
    mkdir build && cd build

    ../configure --prefix=$HOME/glibc-2.28/install
    make -j$(nproc)
    make install
fi
cd $REPO_DIR

BUILD_NAME=octree-$OS-$ARCH.so 
echo "Building shared library file: '$BUILD_NAME'"

# with linux use glibc version 2.28 that we just downloaded
if [[ "$UNAME_OS" == "Linux" ]]; then
    GOOS=$OS GOARCH=$ARCH CGO_ENABLED=1 \
        CC="gcc -nostdlib -Wl,--rpath=$HOME/glibc-2.28/install/lib" \
        LD_LIBRARY_PATH="$HOME/glibc-2.28/install/lib" \
        LIBRARY_PATH="$HOME/glibc-2.28/install/lib" \
        C_INCLUDE_PATH="$HOME/glibc-2.28/install/include"\
        go build -buildmode=c-shared -o $BUILD_NAME $SRC_FILE
else
    GOOS=$OS GOARCH=$ARCH CGO_ENABLED=1 \
        go build -buildmode=c-shared -o $BUILD_NAME $SRC_FILE
fi