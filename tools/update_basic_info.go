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

func NewUpdateBasicInfoTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("update_basic_info",
		mcp.WithDescription("Update basic information (name, photo, description) of an existing resume. Only provide fields you want to update."),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("The ID of the resume to update"),
		),
		mcp.WithString("name",
			mcp.Description("The name of the resume owner"),
		),
		mcp.WithString("photo",
			mcp.Description("URL or path to the photo"),
		),
		mcp.WithString("description",
			mcp.Description("Brief description or summary"),
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

		name := request.GetString("name", "")
		photo := request.GetString("photo", "")
		description := request.GetString("description", "")

		resume, err := db.GetResumeByID(uint(resumeID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Resume not found: %v", err)), nil
		}

		if name != "" {
			resume.Name = name
		}
		if photo != "" {
			resume.Photo = photo
		}
		if description != "" {
			resume.Description = description
		}

		if err := db.UpdateResume(resume); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error updating resume: %v", err)), nil
		}

		result := map[string]interface{}{
			"id":          resume.ID,
			"name":        resume.Name,
			"photo":       resume.Photo,
			"description": resume.Description,
			"updated_at":  resume.UpdatedAt,
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Resume updated successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}