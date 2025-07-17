package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/models"
)

func NewAddOtherExperienceTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("add_other_experience",
		mcp.WithDescription("Add other categorized experiences to a resume (skills, awards, certifications, projects, etc.). Use feature maps to add detailed information."),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("The ID of the resume to add other experience to"),
		),
		mcp.WithString("category",
			mcp.Required(),
			mcp.Description("The category of the experience (skills, awards, certifications, etc.)"),
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

		category, err := request.RequireString("category")
		if err != nil {
			return nil, fmt.Errorf("category parameter is required: %w", err)
		}

		otherExp := &models.OtherExperience{
			ResumeID: uint(resumeID),
			Category: category,
		}

		if err := db.AddOtherExperience(otherExp); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error adding other experience: %v", err)), nil
		}

		result := map[string]interface{}{
			"id":        otherExp.ID,
			"resume_id": otherExp.ResumeID,
			"category":  otherExp.Category,
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Other experience added successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}