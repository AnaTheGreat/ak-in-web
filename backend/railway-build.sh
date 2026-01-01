#!/bin/bash
# Clean any leftover test files
rm -f *_test.go *_debug.go *_function.go test_*.go debug_*.go

# Build Go app
go mod tidy
go build -ldflags="-w -s" -o out
