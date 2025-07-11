#!/bin/bash

# Simple test script for Hardcover API

echo "Hardcover Book Embed - API Test Script"
echo "======================================"

# Check if username is provided
if [ $# -eq 0 ]; then
    echo "Error: Username is required"
    echo "Usage: ./test.sh <username> [book-type]"
    echo "  book-type: currently-reading (default) or last-read"
    exit 1
fi

USERNAME=$1

# Check if API token is set
if [ -z "$HARDCOVER_API_TOKEN" ]; then
    echo "Error: HARDCOVER_API_TOKEN environment variable is not set"
    echo ""
    echo "To run this test:"
    echo "1. Get your API token from https://hardcover.app/account/api"
    echo "2. Export it: export HARDCOVER_API_TOKEN=your_token_here"
    echo "3. Run this script again: ./test.sh <username>"
    exit 1
fi

echo "Building test application..."
cd "$(dirname "$0")"
go build -o test_hardcover test_hardcover.go

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

echo "Running Hardcover API test for user: $USERNAME"
echo ""

# Test with optional book type parameter (defaults to currently-reading)
BOOK_TYPE=${2:-"currently-reading"}
./test_hardcover "$USERNAME" "$BOOK_TYPE"

# Clean up
rm -f test_hardcover

echo ""
echo "Test completed!"