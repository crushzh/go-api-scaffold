# Go Scaffold - Simple Template Makefile
# Cross-platform build, code generator, Swagger docs, Proto Buffer

# ==================== Variables ====================
APP_NAME := myapp
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
BUILD_TIME := $(shell date +%Y-%m-%dT%H:%M:%S)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go build flags
GO := go
GOFLAGS := -trimpath
LDFLAGS := -s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.GitCommit=$(GIT_COMMIT)'

# Output directory
BUILD_DIR := build

# ==================== Cross-compile toolchain ====================
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    CC_LINUX_AMD64 ?= x86_64-linux-gnu-gcc
    CC_LINUX_ARM64 ?= aarch64-linux-gnu-gcc
    CC_LINUX_ARM32 ?= arm-linux-gnueabihf-gcc
    CC_WINDOWS ?= gcc
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Darwin)
        DETECTED_OS := macOS
        CC_LINUX_AMD64 ?= x86_64-unknown-linux-gnu-gcc
        CC_LINUX_ARM64 ?= aarch64-unknown-linux-gnu-gcc
        CC_LINUX_ARM32 ?= arm-unknown-linux-gnueabihf-gcc
        CC_WINDOWS ?= x86_64-w64-mingw32-gcc
    else
        DETECTED_OS := Linux
        CC_LINUX_AMD64 ?= gcc
        CC_LINUX_ARM64 ?= aarch64-linux-gnu-gcc
        CC_LINUX_ARM32 ?= arm-linux-gnueabihf-gcc
        CC_WINDOWS ?= x86_64-w64-mingw32-gcc
    endif
endif

# ==================== Default target ====================
.PHONY: all
all: build

# ==================== Development ====================

# Run in dev mode
.PHONY: run
run:
	$(GO) run ./cmd/server/ -c configs/config.yaml

# Build locally
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server/
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

# ==================== Frontend ====================

# Build frontend (configure frontend project path as needed)
.PHONY: web
web:
	@echo "Building frontend..."
	@if [ -d "../$(APP_NAME)-web" ]; then \
		cd ../$(APP_NAME)-web && npm run build; \
		rm -rf internal/web/dist; \
		cp -r ../$(APP_NAME)-web/dist internal/web/dist; \
		echo "Frontend build complete"; \
	else \
		echo "Warning: Frontend project not found at ../$(APP_NAME)-web"; \
	fi

# ==================== Cross-platform build ====================

.PHONY: build-linux
build-linux:
	@echo "Building for Linux amd64..."
	@mkdir -p $(BUILD_DIR)/linux-amd64
	CGO_ENABLED=1 CC=$(CC_LINUX_AMD64) GOOS=linux GOARCH=amd64 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/linux-amd64/$(APP_NAME) ./cmd/server/
	@cp -r configs $(BUILD_DIR)/linux-amd64/
	@echo "Build complete: $(BUILD_DIR)/linux-amd64/$(APP_NAME)"

.PHONY: build-arm64
build-arm64:
	@echo "Building for Linux arm64..."
	@mkdir -p $(BUILD_DIR)/linux-arm64
	CGO_ENABLED=1 CC=$(CC_LINUX_ARM64) GOOS=linux GOARCH=arm64 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/linux-arm64/$(APP_NAME) ./cmd/server/
	@cp -r configs $(BUILD_DIR)/linux-arm64/
	@echo "Build complete: $(BUILD_DIR)/linux-arm64/$(APP_NAME)"

.PHONY: build-arm32
build-arm32:
	@echo "Building for Linux arm32 (ARMv7)..."
	@mkdir -p $(BUILD_DIR)/linux-arm32
	CGO_ENABLED=1 CC=$(CC_LINUX_ARM32) GOOS=linux GOARCH=arm GOARM=7 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/linux-arm32/$(APP_NAME) ./cmd/server/
	@cp -r configs $(BUILD_DIR)/linux-arm32/
	@echo "Build complete: $(BUILD_DIR)/linux-arm32/$(APP_NAME)"

.PHONY: build-windows
build-windows:
	@echo "Building for Windows amd64..."
	@mkdir -p $(BUILD_DIR)/windows-amd64
	CGO_ENABLED=1 CC=$(CC_WINDOWS) GOOS=windows GOARCH=amd64 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/windows-amd64/$(APP_NAME).exe ./cmd/server/
	@cp -r configs $(BUILD_DIR)/windows-amd64/
	@echo "Build complete: $(BUILD_DIR)/windows-amd64/$(APP_NAME).exe"

.PHONY: build-all
build-all: build-linux build-arm64 build-arm32 build-windows
	@echo "All platforms build complete!"

# ==================== Code generation ====================

# Generate module code
# Usage: make gen name=order cn=Order
.PHONY: gen
gen:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make gen name=order cn=Order"; \
		exit 1; \
	fi
	$(GO) run ./cmd/gen/ -name $(name) -cn "$(cn)"

# ==================== Docs & Proto ====================

# Install swag tool
.PHONY: swag-install
swag-install:
	@echo "Installing swag..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Swag installed"

# Generate Swagger docs
.PHONY: docs
docs:
	@echo "Generating swagger docs..."
	swag init -d ./cmd/server,./internal -g main.go -o docs/swagger --parseDependency --parseInternal
	@echo "Swagger docs generated at docs/swagger/"

# Install protoc-gen-go tools
.PHONY: proto-install
proto-install:
	@echo "Installing protoc-gen-go..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Proto tools installed"

# Generate Proto Buffer code
.PHONY: protos
protos:
	@echo "Generating proto files..."
	@mkdir -p api/proto/gen
	protoc --go_out=api/proto/gen --go_opt=paths=source_relative \
		--go-grpc_out=api/proto/gen --go-grpc_opt=paths=source_relative \
		api/proto/*.proto
	@echo "Proto generation complete"

# ==================== Code quality ====================

.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...

.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -v -race -cover ./...

.PHONY: coverage
coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# ==================== Dependencies & cleanup ====================

.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# ==================== Info ====================

.PHONY: info
info:
	@echo "=========================================="
	@echo "Build Configuration"
	@echo "=========================================="
	@echo "App:             $(APP_NAME)"
	@echo "Version:         $(VERSION)"
	@echo "Detected OS:     $(DETECTED_OS)"
	@echo "CC_LINUX_AMD64:  $(CC_LINUX_AMD64)"
	@echo "CC_LINUX_ARM64:  $(CC_LINUX_ARM64)"
	@echo "CC_LINUX_ARM32:  $(CC_LINUX_ARM32)"
	@echo "CC_WINDOWS:      $(CC_WINDOWS)"
	@echo "=========================================="

.PHONY: help
help:
	@echo "$(APP_NAME) Build System"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Development:"
	@echo "  run             Run in dev mode"
	@echo "  build           Build locally"
	@echo "  web             Build frontend and copy to dist/"
	@echo ""
	@echo "Cross Compile:"
	@echo "  build-linux     Build Linux amd64"
	@echo "  build-arm64     Build Linux arm64"
	@echo "  build-arm32     Build Linux arm32 (ARMv7)"
	@echo "  build-windows   Build Windows"
	@echo "  build-all       Build all platforms"
	@echo ""
	@echo "Code Generator:"
	@echo "  gen             Generate module code (make gen name=order cn=Order)"
	@echo ""
	@echo "Documentation:"
	@echo "  docs            Generate Swagger docs"
	@echo "  protos          Generate Proto Buffer code"
	@echo "  swag-install    Install swag tool"
	@echo "  proto-install   Install protoc-gen-go tools"
	@echo ""
	@echo "Quality:"
	@echo "  lint            Run linter"
	@echo "  test            Run tests"
	@echo "  coverage        Test coverage"
	@echo "  fmt             Format code"
	@echo ""
	@echo "Misc:"
	@echo "  deps            Download dependencies"
	@echo "  clean           Clean build artifacts"
	@echo "  info            Show build configuration"
	@echo "  help            Show this help"
	@echo ""
	@echo "Cross-compile toolchain:"
	@echo "  macOS:   brew tap messense/macos-cross-toolchains"
	@echo "           brew install aarch64-unknown-linux-gnu"
	@echo "           brew install arm-unknown-linux-gnueabihf"
	@echo "  Linux:   apt install gcc-aarch64-linux-gnu"
	@echo "           apt install gcc-arm-linux-gnueabihf"
