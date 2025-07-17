package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
)

func NewListTemplatesTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("list_templates",
		mcp.WithDescription("List all templates for a specific resume"),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("ID of the resume to list templates for"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resumeIDStr, err := request.RequireString("resume_id")
		if err != nil {
			return nil, fmt.Errorf("resume_id parameter is required: %w", err)
		}

		resumeID, err := strconv.Atoi(resumeIDStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resume_id: %v", err)), nil
		}

		templates, err := db.ListTemplatesByResumeID(uint(resumeID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list templates: %v", err)), nil
		}

		result := map[string]interface{}{
			"success":   true,
			"templates": templates,
			"count":     len(templates),
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Templates listed successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}