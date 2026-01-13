#!/bin/bash

echo "Building Zuon for Linux..."

# 1. Build Server
echo "Building server..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/zuon-server ./cmd/server

# 2. Build UI Client
# Check if fyne-cross is installed
if ! command -v fyne-cross &> /dev/null
then
    echo "fyne-cross could not be found. Installing..."
    go install github.com/fyne-io/fyne-cross@latest
fi

echo "Building UI client using fyne-cross..."
fyne-cross linux -arch amd64 ./cmd/zuon

# Move to build directory
mkdir -p build
cp fyne-cross/bin/linux-amd64/zuon build/zuon-linux

echo "Build complete! Files are in the 'build' directory."
