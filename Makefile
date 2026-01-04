# Makefile for arctisbattery

# Detect OS
ifeq ($(OS),Windows_NT)
	DETECTED_OS := Windows
else
	DETECTED_OS := $(shell uname -s)
endif

# Project variables
BINARY_NAME=arctisbattery
PROJECT_NAME=arctisbattery
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1

# Build directory
DIST_DIR=dist
BUILD_DIR=$(DIST_DIR)/$(PROJECT_NAME)_Windows_x86_64

# Linker flags
LDFLAGS=-H=windowsgui

# OS-specific commands
ifeq ($(DETECTED_OS),Windows)
	RM=if exist $(DIST_DIR) rmdir /s /q $(DIST_DIR)
	MKDIR=if not exist $(DIST_DIR) mkdir $(DIST_DIR)
	CP=copy
	PATH_SEP=\\
else
	RM=rm -rf $(DIST_DIR)
	MKDIR=mkdir -p $(DIST_DIR)
	CP=cp
	PATH_SEP=/
endif

.PHONY: all
all: clean build

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@$(RM)
	@$(MKDIR)

.PHONY: tidy
tidy:
	@echo "Running go mod tidy..."
	@go mod tidy

.PHONY: build
build: tidy
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME).exe .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME).exe"

.PHONY: install
install: build
	@echo "Installing to GOPATH/bin..."
ifeq ($(DETECTED_OS),Windows)
	@copy $(BUILD_DIR)\$(BINARY_NAME).exe $(GOPATH)\bin\$(BINARY_NAME).exe
else
	@cp $(BUILD_DIR)/$(BINARY_NAME).exe $(GOPATH)/bin/$(BINARY_NAME).exe
endif

.PHONY: tag
tag:
	@echo "Creating new version tag..."
ifeq ($(DETECTED_OS),Windows)
	@for /f "usebackq" %%i in (`svu next`) do @set TAG=%%i && git tag %%i && git push origin %%i && echo Created and pushed tag: %%i
else
	@TAG=$$(svu next) && git tag $$TAG && git push origin $$TAG && echo "Created and pushed tag: $$TAG"
endif

.PHONY: run
run:
	@go run -ldflags "$(LDFLAGS)" .

.PHONY: test
test:
	@echo "Running tests..."
	@go test ./... -v

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-short
test-short:
	@echo "Running short tests..."
	@go test ./... -short

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

.PHONY: vet
vet:
	@echo "Running go vet..."
	@go vet ./...

.PHONY: lint
lint: fmt vet
	@echo "Linting complete"

.PHONY: check
check: lint test
	@echo "All checks passed"

.PHONY: help
help:
	@echo Available targets:
	@echo   all           - Clean and build (default)
	@echo   clean         - Remove build artifacts
	@echo   tidy          - Run go mod tidy
	@echo   build         - Build the binary
	@echo   install       - Install binary to GOPATH/bin
	@echo   tag           - Create and push new version tag using svu
	@echo   run           - Run the application
	@echo   test          - Run all tests
	@echo   test-coverage - Run tests with coverage report
	@echo   test-short    - Run short tests
	@echo   fmt           - Format code
	@echo   vet           - Run go vet
	@echo   lint          - Format and vet code
	@echo   check         - Run linting and tests
	@echo   help          - Show this help message
