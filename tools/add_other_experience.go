package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/models"
	"github.com/rxtech-lab/resume-mcp/internal/types"
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
		user := types.GetAuthenticatedUser(ctx)
		userID := &user.Sub

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

		if err := db.AddOtherExperience(otherExp, userID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error adding other experience: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Other experience added successfully")), nil
	}

	return tool, handler
}
