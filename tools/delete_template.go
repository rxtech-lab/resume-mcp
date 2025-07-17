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

func NewDeleteTemplateTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("delete_template",
		mcp.WithDescription("Delete a template by ID"),
		mcp.WithString("template_id",
			mcp.Required(),
			mcp.Description("ID of the template to delete"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templateIDStr, err := request.RequireString("template_id")
		if err != nil {
			return nil, fmt.Errorf("template_id parameter is required: %w", err)
		}

		templateID, err := strconv.Atoi(templateIDStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid template_id: %v", err)), nil
		}

		// Check if template exists first
		template, err := db.GetTemplateByID(uint(templateID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Template not found: %v", err)), nil
		}

		if err := db.DeleteTemplate(uint(templateID)); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete template: %v", err)), nil
		}

		result := map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Template '%s' deleted successfully", template.Name),
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Template deleted successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}