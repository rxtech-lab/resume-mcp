package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
)

func NewDeleteResumeTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("delete_resume",
		mcp.WithDescription("Delete a resume and all associated data (contacts, experiences, feature maps) by ID. This action cannot be undone."),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("The ID of the resume to delete"),
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

		if err := db.DeleteResume(uint(resumeID)); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error deleting resume: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Resume with ID %d deleted successfully", resumeID)), nil
	}

	return tool, handler
}