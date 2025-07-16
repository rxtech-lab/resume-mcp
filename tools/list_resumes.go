package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
)

func NewListResumesTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("list_resumes",
		mcp.WithDescription("List all saved resumes"),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resumes, err := db.ListResumes()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error listing resumes: %v", err)), nil
		}

		resultJSON, _ := json.Marshal(resumes)
		return mcp.NewToolResultText(fmt.Sprintf("Resumes: %s", string(resultJSON))), nil
	}

	return tool, handler
}