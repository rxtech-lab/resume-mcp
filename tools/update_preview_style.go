package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
)

func NewUpdatePreviewStyleTool(db *database.Database, port string) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("update_preview_style",
		mcp.WithDescription("Update CSS styles for an existing preview session. Tailwind CSS classes are available for styling."),
		mcp.WithString("session_id",
			mcp.Required(),
			mcp.Description("The session ID of the preview to update"),
		),
		mcp.WithString("css",
			mcp.Required(),
			mcp.Description("New CSS styles for the preview"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		sessionID, err := request.RequireString("session_id")
		if err != nil {
			return nil, fmt.Errorf("session_id parameter is required: %w", err)
		}

		css, err := request.RequireString("css")
		if err != nil {
			return nil, fmt.Errorf("css parameter is required: %w", err)
		}

		if err := db.UpdatePreviewSessionCSS(sessionID, css); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error updating preview style: %v", err)), nil
		}

		previewURL := fmt.Sprintf("http://localhost:%s/resume/preview/%s", port, sessionID)
		return mcp.NewToolResultText(fmt.Sprintf("Preview style updated successfully. URL: %s", previewURL)), nil
	}

	return tool, handler
}