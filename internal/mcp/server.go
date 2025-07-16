package mcp

import (
	"log"

	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/tools"
)

type MCPServer struct {
	server *server.MCPServer
}

func NewMCPServer(db *database.Database) *MCPServer {
	mcpServer := &MCPServer{}

	mcpServer.InitializeTools()
	return mcpServer
}

func (s *MCPServer) InitializeTools() {
	srv := server.NewMCPServer(
		"Resume MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	db, err := database.NewDatabase("resume.db")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	createResumeTool, createResumeHandler := tools.NewCreateResumeTool(db)
	srv.AddTool(createResumeTool, createResumeHandler)

	s.server = srv
}

func (s *MCPServer) Start() error {
	return server.ServeStdio(s.server)
}
