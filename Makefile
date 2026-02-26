# Date: 2026-02-25
# Copyright (c) 2026. All rights reserved.

BINARY_NAME=skyclerk
VERSION ?= dev
LDFLAGS=-ldflags "-X github.com/cloudmanic/skyclerk-cli/cmd.Version=$(VERSION)"

## build: Build the binary for the current platform
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

## install: Build and install to GOPATH/bin
install:
	go install $(LDFLAGS) .

## test: Run all tests
test:
	go test ./... -v

## test-short: Run tests without verbose output
test-short:
	go test ./...

## coverage: Generate test coverage report
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

## lint: Run go fmt and go vet
lint:
	go fmt ./...
	go vet ./...

## clean: Remove build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

## cross-build: Build for multiple platforms
cross-build:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe .

.PHONY: build install test test-short coverage lint clean cross-build
