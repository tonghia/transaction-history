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
MAIN=cmd/main.go

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
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN)

# Run the application with sample arguments
# You can modify the PERIOD and FILE_PATH as needed
run: build
	@echo "🚀 Running the application..."
	./$(BUILD_DIR)/$(BINARY_NAME) -interactive

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
	$(GO) test -cover $(shell go list ./... | grep -Fv -e /cmd -e /script ) -v

# Lint the code (requires golint to be installed)
lint:
	@echo "🔍 Linting the code..."
	golint ./...

# Cross-compile for multiple platforms
# Builds binaries for Linux, Windows, and macOS
cross-compile:
	@echo "🌐 Cross-compiling for different platforms..."
	env GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 -ldflags="$(LDFLAGS)" $(MAIN)
	env GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe -ldflags="$(LDFLAGS)" $(MAIN)
	env GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 -ldflags="$(LDFLAGS)" $(MAIN)
	@echo "✅ Cross-compilation completed."


# ============================
#         PHONY
# ============================

.PHONY: all build run clean format test lint cross-compile install
