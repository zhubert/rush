# Rush Programming Language Makefile

# Configuration
BINARY_NAME = rush
BUILD_CMD = cmd/rush/main.go
INSTALL_PREFIX = /usr/local
INSTALL_DIR = $(INSTALL_PREFIX)/lib/rush
BIN_DIR = $(INSTALL_PREFIX)/bin

# Default target
.PHONY: help
help:
	@echo "Rush Programming Language Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  build     - Build the Rush binary"
	@echo "  install   - Install Rush system-wide (requires sudo)"
	@echo "  uninstall - Remove Rush from system (requires sudo)"
	@echo "  clean     - Remove build artifacts"
	@echo "  test      - Run all tests"
	@echo "  dev       - Run Rush from source (development mode)"
	@echo "  repl      - Start Rush REPL from source"
	@echo "  help      - Show this help message"

# Build the Rush binary
.PHONY: build
build:
	@echo "Building Rush..."
	go build -tags llvm20 -o $(BINARY_NAME) $(BUILD_CMD)
	@echo "Build complete: $(BINARY_NAME)"

# Install Rush system-wide
.PHONY: install
install: build
	@echo "Installing Rush to $(INSTALL_PREFIX)..."
	sudo mkdir -p $(INSTALL_DIR)
	sudo cp $(BINARY_NAME) $(INSTALL_DIR)/
	sudo cp -r std/ $(INSTALL_DIR)/std/
	sudo ln -sf $(INSTALL_DIR)/$(BINARY_NAME) $(BIN_DIR)/$(BINARY_NAME)
	@echo ""
	@echo "Rush installed successfully!"
	@echo "- Binary: $(BIN_DIR)/$(BINARY_NAME)"
	@echo "- Standard Library: $(INSTALL_DIR)/std/"
	@echo ""
	@echo "You can now run 'rush' from anywhere on your system."
	@echo "Try: rush --help"

# Uninstall Rush from system
.PHONY: uninstall
uninstall:
	@echo "Uninstalling Rush..."
	sudo rm -f $(BIN_DIR)/$(BINARY_NAME)
	sudo rm -rf $(INSTALL_DIR)
	@echo "Rush uninstalled successfully."

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	@echo "Clean complete."

# Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Development mode - run from source
.PHONY: dev
dev:
	@if [ -z "$(FILE)" ]; then \
		echo "Usage: make dev FILE=path/to/file.rush"; \
		echo "Example: make dev FILE=examples/comprehensive_demo.rush"; \
	else \
		go run -tags llvm20 $(BUILD_CMD) $(FILE); \
	fi

# Start REPL from source
.PHONY: repl
repl:
	@echo "Starting Rush REPL..."
	go run -tags llvm20 $(BUILD_CMD)