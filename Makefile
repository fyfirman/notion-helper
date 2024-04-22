# Makefile

# Include .env file and export its variables
-include .env
export

BINARY_NAME := notion-helper


# Build the Go application
build:
	go build -o $(BINARY_NAME) ./cmd

# Run the Go application
run:
	go run ./cmd/main.go

# Run test for all directories
lint:
	golangci-lint run ./...

# Run test for all directories
test:
	go test -v ./...

# Run test for CI mode
test-ci:
	go test -coverprofile=coverage.out ./... | grep -v "/mocks"

# Run test get coverage with HTML format
test-html:
	go test -coverprofile=coverage.out.tmp ./... ; cat coverage.out.tmp | grep -v "/mocks" > coverage.out ;  go tool cover -html=coverage.out