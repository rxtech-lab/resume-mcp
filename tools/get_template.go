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

func NewGetTemplateTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("get_template",
		mcp.WithDescription("Get a specific template by ID"),
		mcp.WithString("template_id",
			mcp.Required(),
			mcp.Description("ID of the template to retrieve"),
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

		template, err := db.GetTemplateByID(uint(templateID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Template not found: %v", err)), nil
		}

		result := map[string]interface{}{
			"success":  true,
			"template": template,
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Template retrieved successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}