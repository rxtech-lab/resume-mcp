package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/models"
)

func NewCreateResumeTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("create_resume",
		mcp.WithDescription("Create a new resume with basic information"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the resume owner"),
		),
		mcp.WithString("photo",
			mcp.Description("URL or path to the photo"),
		),
		mcp.WithString("description",
			mcp.Required(),
			mcp.Description("Brief description or summary"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return nil, fmt.Errorf("name parameter is required: %w", err)
		}

		description, err := request.RequireString("description")
		if err != nil {
			return nil, fmt.Errorf("description parameter is required: %w", err)
		}

		photo := request.GetString("photo", "")

		resume := &models.Resume{
			Name:        name,
			Photo:       photo,
			Description: description,
		}

		if err := db.CreateResume(resume); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error creating resume: %v", err)), nil
		}

		result := map[string]interface{}{
			"id":          resume.ID,
			"name":        resume.Name,
			"photo":       resume.Photo,
			"description": resume.Description,
			"created_at":  resume.CreatedAt,
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Resume created successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}
