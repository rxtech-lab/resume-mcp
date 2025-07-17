.PHONY: build test run clean help

BINARY_NAME=resume-mcp
BUILD_DIR=./bin

# Default target
all: build

# Build the project
build:
	go build -o bin/resume-mcp ./cmd/main.go
	go build -o bin/config-updater ./cmd/config-updater/main.go

# Run tests
test:
	go test ./...

# Run the MCP server
run:
	go run ./cmd/main.go

# Clean build artifacts
clean:
	rm -rf bin/
	sudo rm -rf /usr/local/bin/$(BINARY_NAME)

install-local: clean build ## Install the binary to /usr/local/bin (requires sudo)
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed successfully!"
	@echo "You can now run '$(BINARY_NAME)' from anywhere."


package: build
	./scripts/sign.sh
	./scripts/package-notarize.sh


# Show help
help:
	@echo "Available targets:"
	@echo "  build  - Build the project"
	@echo "  test   - Run tests"
	@echo "  run    - Run the MCP server"
	@echo "  clean  - Clean build artifacts"
	@echo "  help   - Show this help message"