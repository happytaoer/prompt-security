.PHONY: all build clean linux macos windows help

VERSION ?= latest
OUTPUT_DIR = dist
BINARY_NAME = prompt-security
LDFLAGS = -s -w

all: clean build

help:
	@echo "Prompt Security - Build System"
	@echo ""
	@echo "Usage:"
	@echo "  make build          - Build for all platforms"
	@echo "  make linux          - Build for Linux only"
	@echo "  make macos          - Build for macOS only"
	@echo "  make windows        - Build for Windows only"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "Options:"
	@echo "  VERSION=x.x.x       - Set version (default: latest)"

build: linux macos windows
	@echo ""
	@echo "✅ Build complete! Binaries are in the '$(OUTPUT_DIR)' directory"
	@ls -lh $(OUTPUT_DIR)/

linux:
	@echo "📦 Building for Linux..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-amd64 -ldflags "$(LDFLAGS)" .
	GOOS=linux GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-linux-arm64 -ldflags "$(LDFLAGS)" .

macos:
	@echo "📦 Building for macOS..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 -ldflags "$(LDFLAGS)" .
	GOOS=darwin GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 -ldflags "$(LDFLAGS)" .

windows:
	@echo "📦 Building for Windows..."
	@mkdir -p $(OUTPUT_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-amd64.exe -ldflags "$(LDFLAGS)" .
	GOOS=windows GOARCH=arm64 go build -o $(OUTPUT_DIR)/$(BINARY_NAME)-windows-arm64.exe -ldflags "$(LDFLAGS)" .

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(OUTPUT_DIR)
	@echo "✅ Clean complete!"
