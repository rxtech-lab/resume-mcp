package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rxtech-lab/resume-mcp/internal/api"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/mcp"
)

func main() {
	log.Println("Starting Resume MCP Server...")

	db, err := database.NewDatabase("resume.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	mcpServer := mcp.NewMCPServer(db)
	
	apiServer := api.NewAPIServer(db)

	go func() {
		log.Println("Starting API server on port 8080...")
		if err := apiServer.Start("8080"); err != nil {
			log.Fatal("Failed to start API server:", err)
		}
	}()

	go func() {
		log.Println("Starting MCP server...")
		if err := mcpServer.Start(); err != nil {
			log.Fatal("Failed to start MCP server:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down servers...")
	if err := apiServer.Shutdown(); err != nil {
		log.Printf("Error shutting down API server: %v", err)
	}

	log.Println("Resume MCP Server stopped")
}