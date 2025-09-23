package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/types"
)

func NewDeleteTemplateTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("delete_template",
		mcp.WithDescription("Delete a template by ID"),
		mcp.WithString("template_id",
			mcp.Required(),
			mcp.Description("ID of the template to delete"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		user := types.GetAuthenticatedUser(ctx)
		userID := &user.Sub

		templateIDStr, err := request.RequireString("template_id")
		if err != nil {
			return nil, fmt.Errorf("template_id parameter is required: %w", err)
		}

		templateID, err := strconv.Atoi(templateIDStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid template_id: %v", err)), nil
		}

		// Check if template exists first
		template, err := db.GetTemplateByID(uint(templateID), userID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Template not found: %v", err)), nil
		}

		if err := db.DeleteTemplate(uint(templateID), userID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete template: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Template deleted successfully: %s", template.Name)), nil
	}

	return tool, handler
}
