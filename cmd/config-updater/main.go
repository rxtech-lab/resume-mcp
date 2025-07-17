package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type MCPServer struct {
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env"`
}

type Config struct {
	MCPServers map[string]MCPServer `json:"mcpServers"`
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	claudeDir := filepath.Join(homeDir, "Library", "Application Support", "Claude")
	configPath := filepath.Join(claudeDir, "claude_desktop_config.json")

	// Check if Claude Desktop is installed by looking for Claude directory or config file
	if _, err := os.Stat(claudeDir); os.IsNotExist(err) {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Println("Claude Desktop not found - skipping configuration")
			fmt.Println("Please install Claude Desktop first if you want to use MCP integration with claude desktop")
			return
		}
	}

	// Create Claude directory if it doesn't exist
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		log.Fatalf("Failed to create Claude directory: %v", err)
	}

	var config Config

	// Read existing config or create new one
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			log.Fatalf("Failed to parse existing config: %v", err)
		}
	} else {
		// Create new config
		config = Config{
			MCPServers: make(map[string]MCPServer),
		}
	}

	// Add resume-mcp server configuration
	config.MCPServers["resume-mcp"] = MCPServer{
		Command: "resume-mcp",
		Args:    []string{"--port", "8123"},
		Env: map[string]string{
			"HOME": homeDir,
		},
	}

	// Write updated config
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}

	fmt.Printf("Successfully updated Claude Desktop config at %s\n", configPath)
}
