package main

import (
	"flag"
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
	// get port from cmd line
	port := flag.String("port", "0", "Port to listen on (0 for any available port)")
	flag.Parse()
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get home directory:", err)
	}

	db, err := database.NewDatabase(homePath + "/resume.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	templateService := service.NewTemplateService()
	
	// Create API server first
	apiServer := api.NewAPIServer(db, templateService)
	
	// Start API server and get the actual port
	actualPort, err := apiServer.Start(*port)
	if err != nil {
		log.SetFlags(0)
		log.Fatal("Failed to start API server:", err)
	}
	
	// Create MCP server with the actual port
	mcpServer := mcp.NewMCPServer(db, actualPort, templateService)

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
