#!/bin/bash
# macOS Native Installation Script for Packemon
# This script helps set up and build Packemon natively on macOS

echo "Packemon macOS Installation"
echo "=========================="
echo

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go first."
    echo "Visit: https://golang.org/doc/install"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "Detected Go version: $GO_VERSION"

# Check if libpcap is installed
if ! brew list libpcap &> /dev/null; then
    echo "Installing libpcap using Homebrew..."
    brew install libpcap
    if [ $? -ne 0 ]; then
        echo "Failed to install libpcap. Please install it manually."
        echo "Run: brew install libpcap"
        exit 1
    fi
else
    echo "libpcap is already installed."
fi

# Install required Go dependencies
echo "Installing Go dependencies..."
go mod tidy

# Build Packemon for macOS
echo "Building Packemon for macOS..."
go build -o packemon cmd/packemon/*.go

if [ $? -eq 0 ]; then
    echo
    echo "Build successful! Packemon has been built for macOS."
    echo
    echo "To run Packemon in monitor mode, use:"
    echo "sudo ./packemon"
    echo
    echo "To run Packemon in generator mode, use:"
    echo "sudo ./packemon --send"
    echo
    echo "Note: Packemon requires administrator privileges on macOS to capture and send packets."
else
    echo "Build failed. Please check the error messages above."
    exit 1
fi
