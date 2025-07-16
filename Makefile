.PHONY: build test run clean help

# Default target
all: build

# Build the project
build:
	go build -o bin/resume-mcp ./cmd/main.go

# Run tests
test:
	go test ./...

# Run the MCP server
run:
	go run ./cmd/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Show help
help:
	@echo "Available targets:"
	@echo "  build  - Build the project"
	@echo "  test   - Run tests"
	@echo "  run    - Run the MCP server"
	@echo "  clean  - Clean build artifacts"
	@echo "  help   - Show this help message"