package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
)

func NewGeneratePreviewTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("generate_preview",
		mcp.WithDescription("Generate HTML preview with template and CSS"),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("The ID of the resume to generate preview for"),
		),
		mcp.WithString("template",
			mcp.Required(),
			mcp.Description("Go template string for rendering"),
		),
		mcp.WithString("css",
			mcp.Description("CSS styles for the preview"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resumeIDStr, err := request.RequireString("resume_id")
		if err != nil {
			return nil, fmt.Errorf("resume_id parameter is required: %w", err)
		}

		resumeID, err := strconv.ParseUint(resumeIDStr, 10, 32)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resume_id: %v", err)), nil
		}

		template, err := request.RequireString("template")
		if err != nil {
			return nil, fmt.Errorf("template parameter is required: %w", err)
		}

		css := request.GetString("css", "")

		sessionID, err := db.GeneratePreview(uint(resumeID), template, css)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error generating preview: %v", err)), nil
		}

		previewURL := fmt.Sprintf("http://localhost:8080/resume/preview/%s", sessionID)
		return mcp.NewToolResultText(fmt.Sprintf("Preview generated successfully. URL: %s", previewURL)), nil
	}

	return tool, handler
}