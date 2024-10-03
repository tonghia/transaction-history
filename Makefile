# Makefile for transaction-history

# ============================
#         VARIABLES
# ============================

# Binary name
BINARY_NAME=transaction-history

# Build directory
BUILD_DIR=build

# Go settings
GO=go
PKG=.

# Versioning
VERSION=$(shell git describe --tags --always --dirty)
# Build flags
LDFLAGS=-X main.version=$(VERSION)

# ============================
#         TARGETS
# ============================

# Default target: build the application
all: build

# Build the binary
build:
	@echo "🏗️  Building the application..."
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(LDFLAGS) $(PKG)

# Run the application with sample arguments
# You can modify the PERIOD and FILE_PATH as needed
run: build
	@echo "🚀 Running the application..."
	./$(BUILD_DIR)/$(BINARY_NAME) -period=202201 -file=transactions.csv

# Clean build artifacts
clean:
	@echo "🧹 Cleaning up..."
	$(GO) clean
	rm -f $(BUILD_DIR)/$(BINARY_NAME)

# Format the code using gofmt
format:
	@echo "🛠️  Formatting the code..."
	$(GO) fmt ./...

# Run tests (assuming you have tests)
test:
	@echo "🧪 Running tests..."
	$(GO) test ./...

# Lint the code (requires golint to be installed)
lint:
	@echo "🔍 Linting the code..."
	golint ./...

# Cross-compile for multiple platforms
# Builds binaries for Linux, Windows, and macOS
cross-compile:
	@echo "🌐 Cross-compiling for different platforms..."
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 -ldflags="$(LDFLAGS)" GOOS=linux GOARCH=amd64 $(PKG)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe -ldflags="$(LDFLAGS)" GOOS=windows GOARCH=amd64 $(PKG)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 -ldflags="$(LDFLAGS)" GOOS=darwin GOARCH=amd64 $(PKG)
	@echo "✅ Cross-compilation completed."

# Install the binary to GOPATH/bin
install:
	@echo "📦 Installing the binary..."
	$(GO) install $(PKG)

# ============================
#         PHONY
# ============================

.PHONY: all build run clean format test lint cross-compile install
