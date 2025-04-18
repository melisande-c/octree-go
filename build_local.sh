#!/bin/bash

# !!! This script was written by chat GPT but I think it works !!!
set -e

SRC_FILE=export.go

# Detect OS
UNAME_OS=$(uname -s)
UNAME_ARCH=$(uname -m)

# Build for macOS (keep as .so)
if [[ "$UNAME_OS" == "Darwin" ]]; then
  echo "🛠️ macOS detected: building for platform..."

  GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 \
    go build -buildmode=c-shared -o libgo_arm64.so $SRC_FILE

  GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 \
    go build -buildmode=c-shared -o libgo_amd64.so $SRC_FILE

  # combine arm64 and amd64
  lipo -create -output libgo.so libgo_arm64.so libgo_amd64.so

  # remove temp builds
  rm libgo_arm64.so
  rm libgo_amd64.so
  rm libgo_arm64.h
  rm libgo_amd64.h

  echo "✅ Built macOS shared lib: libgo.so"

# Build for Windows (output .dll)
elif [[ "$UNAME_OS" == MINGW* || "$UNAME_OS" == MSYS* || "$UNAME_OS" == CYGWIN* ]]; then
  echo "🪟 Windows detected: building for host arch..."

  case "$UNAME_ARCH" in
    x86_64)
      GOARCH=amd64
      ;;
    arm64 | aarch64)
      GOARCH=arm64
      ;;
    *)
      echo "❌ Unsupported architecture: $UNAME_ARCH"
      exit 1
      ;;
  esac

  GOOS=windows GOARCH=$GOARCH CGO_ENABLED=1 \
    go build -buildmode=c-shared -o libgo.dll $SRC_FILE

  echo "✅ Built Windows shared lib: libgo.dll"

# Build for Linux (output .so)
else
  echo "🛠️ Linux detected: building for host arch..."

  case "$UNAME_ARCH" in
    x86_64)
      GOARCH=amd64
      ;;
    arm64 | aarch64)
      GOARCH=arm64
      ;;
    *)
      echo "❌ Unsupported architecture: $UNAME_ARCH"
      exit 1
      ;;
  esac

  GOOS=linux GOARCH=$GOARCH CGO_ENABLED=1 \
    go build -buildmode=c-shared -o libgo.so $SRC_FILE

  echo "✅ Built Linux shared lib: libgo.so"
fi