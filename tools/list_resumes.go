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
		mcp.WithDescription("List all saved resumes with their IDs and names. Use this to find available resumes before generating previews."),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resumes, err := db.ListResumes()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error listing resumes: %v", err)), nil
		}

		simpleResumes := []string{}
		for _, resume := range resumes {
			simpleResumes = append(simpleResumes, fmt.Sprintf("%d: %s", resume.ID, resume.Name))
		}

		resultJSON, _ := json.Marshal(simpleResumes)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(fmt.Sprintf("Resumes found: %d", len(resumes))),
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}

	return tool, handler
}
