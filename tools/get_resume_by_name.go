package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
)

func NewGetResumeByNameTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("get_resume_by_name",
		mcp.WithDescription("Retrieve complete structured resume data by name. Returns all associated contacts, experiences, education, and feature maps for template generation."),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the resume to retrieve"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return nil, fmt.Errorf("name parameter is required: %w", err)
		}

		resume, err := db.GetResumeByName(name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Resume not found: %v", err)), nil
		}

		resultJSON, _ := json.Marshal(resume)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Resume found: for %s", name)),
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}

	return tool, handler
}
