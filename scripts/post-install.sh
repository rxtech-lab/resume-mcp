#!/bin/bash

set -e

echo "Running post-installation configuration..."

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Run the bundled config updater
echo "Updating Claude Desktop configuration..."
/usr/local/bin/config-updater

echo "Post-installation configuration completed successfully!"
echo "The resume-mcp server has been added to your Claude Desktop configuration."
echo "Please restart Claude Desktop to use the new MCP server."