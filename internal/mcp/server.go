package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
)

type MCPServer struct {
	server *server.MCPServer
	tools  *ResumeTools
}

func NewMCPServer(db *database.Database) *MCPServer {
	s := server.NewMCPServer(
		"resume-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
	)
	resumeTools := NewResumeTools(db)

	mcpServer := &MCPServer{
		server: s,
		tools:  resumeTools,
	}

	return mcpServer
}

func (s *MCPServer) Start() error {
	return server.ServeStdio(s.server)
}