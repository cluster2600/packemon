#!/bin/bash
# macOS Docker Setup Script for Packemon
# This script helps set up and run Packemon in Docker on macOS

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed. Please install Docker Desktop for Mac first."
    echo "Visit: https://docs.docker.com/desktop/install/mac/"
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo "Docker is not running. Please start Docker Desktop and try again."
    exit 1
fi

echo "Building Packemon Docker image..."
docker build -t packemon .

echo "Docker image built successfully!"
echo ""
echo "To run Packemon in monitor mode, use:"
echo "docker run --rm -it --privileged packemon"
echo ""
echo "To run Packemon in generator mode, use:"
echo "docker run --rm -it --privileged packemon --send"
echo ""
echo "Note: Network interface access may be limited in Docker on macOS."
echo "For full functionality, consider using a Linux VM or WSL2 on Windows."
