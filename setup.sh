#!/bin/bash

# Sentinel Setup Script
# This script helps you get started with Sentinel development

set -e

echo "================================"
echo "  ⌀ Sentinel Setup Script"
echo "================================"
echo ""

# Colors for output aesthetics
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # no Color

# Check prerequisites
echo "Checking prerequisites..."

# Check Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go is not installed${NC}"
    echo "Please install Go 1.21+ from https://golang.org/dl/"
    exit 1
else
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}Go found: $GO_VERSION${NC}"
fi

# Check Git
if ! command -v git &> /dev/null; then
    echo -e "${RED}Git is not installed${NC}"
    echo "Please install Git from https://git-scm.com/downloads"
    exit 1
else
    GIT_VERSION=$(git --version | awk '{print $3}')
    echo -e "${GREEN}Git found: $GIT_VERSION${NC}"
fi

# Check Make (optional)
if command -v make &> /dev/null; then
    MAKE_VERSION=$(make --version | head -n1)
    echo -e "${GREEN}Make found: $MAKE_VERSION${NC}"
    USE_MAKE=true
else
    echo -e "${YELLOW}Make not found (optional)${NC}"
    USE_MAKE=false
fi

echo ""
echo "Downloading dependencies..."
echo "This may take a minute on first run..."
go mod tidy
go mod download

echo ""
echo "Building Sentinel..."

if [ "$USE_MAKE" = true ]; then
    make build
    BINARY_PATH="bin/sentinel"
else
    go build -o sentinel
    BINARY_PATH="./sentinel"
fi

echo -e "${GREEN}Build successful!${NC}"
echo ""

# Verify build
if [ -f "$BINARY_PATH" ]; then
    echo "Testing binary..."
    $BINARY_PATH --version || echo "Version: 0.1.0"
    echo ""

    echo -e "${GREEN}Sentinel is ready!${NC}"
    echo ""
    echo "================================"
    echo "  Quick Start Commands"
    echo "================================"
    echo ""
    echo "Run a test scan:"
    echo "  $BINARY_PATH scan --path ./examples --verbose"
    echo ""
    echo "Scan current directory:"
    echo "  $BINARY_PATH scan --path ."
    echo ""
    echo "Run tests:"
    if [ "$USE_MAKE" = true ]; then
        echo "  make test"
    else
        echo "  go test ./..."
    fi
    echo ""
    echo "For more information, see:"
    echo "  - QUICKSTART.md    : User guide"
    echo "  - README.md        : Just a read me"
    echo "  - PROJECT_STRUCTURE.md : Architecture overview"
    echo ""

    # Offer to run example scan
    echo -e "${YELLOW}Would you like to run a demo scan on the example files? (y/n)${NC}"
    read -r response
    if [[ "$response" =~ ^[Yy]$ ]]; then
        echo ""
        echo "Running demo scan..."
        echo "================================"
        $BINARY_PATH scan --path ./examples --verbose
        echo "================================"
        echo ""
        echo -e "${GREEN}Demo complete!${NC}"
    fi

else
    echo -e "${RED}Build failed - binary not found${NC}"
    exit 1
fi

echo ""
echo "Hell Yeah! Setup has been completed!"
