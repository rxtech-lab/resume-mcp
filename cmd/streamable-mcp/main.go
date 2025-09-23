package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rxtech-lab/resume-mcp/internal/api"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/mcp"
	"github.com/rxtech-lab/resume-mcp/internal/service"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	db, err := database.NewPostgresDatabase(os.Getenv("POSTGRES_URL"))
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	templateService := service.NewTemplateService()

	// Create API server first
	apiServer := api.NewAPIServer(db, templateService)

	// Create MCP server with the actual port
	mcpServer := mcp.NewMCPServer(db, port, templateService)
	streamableServer := mcpServer.StartStreamable()
	apiServer.SetupStreamableServer(streamableServer)
	apiServer.SetupRoutes()

	// Start API server and get the actual port
	_, err = apiServer.Start(port)
	if err != nil {
		log.SetFlags(0)
		log.Fatal("Failed to start API server:", err)
	}

	go func() {
		if err := mcpServer.Start(); err != nil {
			log.SetOutput(os.Stderr)
			log.SetFlags(0)
			log.Fatal("Failed to start MCP server:", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	if err := apiServer.Shutdown(); err != nil {
		log.SetOutput(os.Stderr)
		log.SetFlags(0)
		log.Printf("Error shutting down API server: %v", err)
	}
}
