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

# Frontend directory
WEB_DIR := web
WEB_DIST_DEST := internal/web/dist

# ==================== Cross-compile toolchain ====================
# Auto-detected by OS, override via environment variables
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    CC_LINUX_AMD64 ?= x86_64-linux-gnu-gcc
    CC_LINUX_ARM64 ?= D:/tools/gcc-arm/bin/aarch64-linux-gnu-gcc.exe
    CC_LINUX_ARM32 ?= D:/tools/gcc-arm/bin/arm-linux-gnueabihf-gcc.exe
    CC_WINDOWS ?= gcc
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Darwin)
        # macOS (Homebrew: brew tap messense/macos-cross-toolchains)
        DETECTED_OS := macOS
        CC_LINUX_AMD64 ?= x86_64-unknown-linux-gnu-gcc
        CC_LINUX_ARM64 ?= aarch64-unknown-linux-gnu-gcc
        CC_LINUX_ARM32 ?= arm-unknown-linux-gnueabihf-gcc
        CC_WINDOWS ?= x86_64-w64-mingw32-gcc
    else
        # Linux
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

# Build frontend and copy to embed directory
.PHONY: web
web:
	@echo "Building frontend..."
	@cd $(WEB_DIR) && npm install && npm run build
	@$(MAKE) web-copy

# Copy frontend dist to embed directory
.PHONY: web-copy
web-copy:
	@echo "Copying frontend dist to $(WEB_DIST_DEST)..."
	@rm -rf $(WEB_DIST_DEST)
	@mkdir -p $(WEB_DIST_DEST)
	@cp -r $(WEB_DIR)/dist/* $(WEB_DIST_DEST)/
	@echo "Frontend dist copied"

# Clean frontend build
.PHONY: web-clean
web-clean:
	@echo "Cleaning frontend..."
	@rm -rf $(WEB_DIR)/dist
	@rm -rf $(WEB_DIR)/node_modules/.cache
	@rm -rf $(WEB_DIST_DEST)
	@echo "Frontend clean complete"

# ==================== Cross-platform build ====================

# Linux amd64
.PHONY: build-linux
build-linux:
	@echo "Building for Linux amd64..."
	@echo "CC: $(CC_LINUX_AMD64)"
	@mkdir -p $(BUILD_DIR)/linux-amd64
	CGO_ENABLED=1 CC=$(CC_LINUX_AMD64) GOOS=linux GOARCH=amd64 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/linux-amd64/$(APP_NAME) ./cmd/server/
	@cp -r configs $(BUILD_DIR)/linux-amd64/
	@echo "Build complete: $(BUILD_DIR)/linux-amd64/$(APP_NAME)"

# Linux arm64
.PHONY: build-arm64
build-arm64:
	@echo "Building for Linux arm64..."
	@echo "CC: $(CC_LINUX_ARM64)"
	@mkdir -p $(BUILD_DIR)/linux-arm64
	CGO_ENABLED=1 CC=$(CC_LINUX_ARM64) GOOS=linux GOARCH=arm64 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/linux-arm64/$(APP_NAME) ./cmd/server/
	@cp -r configs $(BUILD_DIR)/linux-arm64/
	@echo "Build complete: $(BUILD_DIR)/linux-arm64/$(APP_NAME)"

# Linux arm32 (ARMv7 hard-float)
.PHONY: build-arm32
build-arm32:
	@echo "Building for Linux arm32 (ARMv7 hard-float)..."
	@echo "CC: $(CC_LINUX_ARM32)"
	@mkdir -p $(BUILD_DIR)/linux-arm32
	CGO_ENABLED=1 CC=$(CC_LINUX_ARM32) GOOS=linux GOARCH=arm GOARM=7 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/linux-arm32/$(APP_NAME) ./cmd/server/
	@cp -r configs $(BUILD_DIR)/linux-arm32/
	@echo "Build complete: $(BUILD_DIR)/linux-arm32/$(APP_NAME)"

# Windows amd64
.PHONY: build-windows
build-windows:
	@echo "Building for Windows amd64..."
	@echo "CC: $(CC_WINDOWS)"
	@mkdir -p $(BUILD_DIR)/windows-amd64
	CGO_ENABLED=1 CC=$(CC_WINDOWS) GOOS=windows GOARCH=amd64 \
		$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/windows-amd64/$(APP_NAME).exe ./cmd/server/
	@cp -r configs $(BUILD_DIR)/windows-amd64/
	@echo "Build complete: $(BUILD_DIR)/windows-amd64/$(APP_NAME).exe"

# All platforms
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

# Lint
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Test
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -v -race -cover ./...

# Coverage
.PHONY: coverage
coverage:
	@echo "Running tests with coverage..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Format
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "Format complete"

# ==================== Dependencies & cleanup ====================

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies downloaded"

# Clean
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
	@echo "App Name:        $(APP_NAME)"
	@echo "Version:         $(VERSION)"
	@echo "Git Commit:      $(GIT_COMMIT)"
	@echo "Build Time:      $(BUILD_TIME)"
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
	@echo "  run             开发模式运行"
	@echo "  build           本地编译"
	@echo ""
	@echo "Frontend:"
	@echo "  web             编译前端并复制到嵌入目录"
	@echo "  web-copy        仅复制前端 dist 到嵌入目录"
	@echo "  web-clean       清理前端构建产物"
	@echo ""
	@echo "Cross Compile:"
	@echo "  build-linux     编译 Linux amd64"
	@echo "  build-arm64     编译 Linux arm64 (AArch64)"
	@echo "  build-arm32     编译 Linux arm32 (ARMv7 硬浮点)"
	@echo "  build-windows   编译 Windows"
	@echo "  build-all       全平台编译"
	@echo ""
	@echo "Code Generator:"
	@echo "  gen             生成模块代码 (make gen name=order cn=Order)"
	@echo ""
	@echo "Documentation:"
	@echo "  docs            生成 Swagger 文档"
	@echo "  protos          生成 Proto Buffer 代码"
	@echo "  swag-install    安装 swag 工具"
	@echo "  proto-install   安装 protoc-gen-go 工具"
	@echo ""
	@echo "Quality:"
	@echo "  lint            代码检查"
	@echo "  test            运行测试"
	@echo "  coverage        测试覆盖率"
	@echo "  fmt             格式化代码"
	@echo ""
	@echo "Misc:"
	@echo "  deps            下载依赖"
	@echo "  clean           清理构建产物"
	@echo "  info            显示当前编译配置"
	@echo "  help            显示帮助"
	@echo ""
	@echo "Cross-compile toolchain (auto-detected, can override via environment):"
	@echo "  - macOS:   brew tap messense/macos-cross-toolchains"
	@echo "             brew install aarch64-unknown-linux-gnu      # ARM64"
	@echo "             brew install arm-unknown-linux-gnueabihf    # ARM32"
	@echo "  - Windows: Set CC_LINUX_ARM64/CC_LINUX_ARM32 to your ARM GCC path"
	@echo "  - Linux:   apt install gcc-aarch64-linux-gnu           # ARM64"
	@echo "             apt install gcc-arm-linux-gnueabihf         # ARM32"
