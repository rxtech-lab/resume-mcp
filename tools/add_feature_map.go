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

func NewAddFeatureMapTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("add_feature_map",
		mcp.WithDescription("Add flexible key-value features to any experience (work, education, other). Use this for details like GPA, salary, responsibilities, achievements, skills, etc."),
		mcp.WithString("experience_id",
			mcp.Required(),
			mcp.Description("The ID of the experience to add features to"),
		),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("The feature key"),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("The feature value"),
		),
		mcp.WithString("category",
			mcp.Required(),
			mcp.Description("The category of the feature"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		user := types.GetAuthenticatedUser(ctx)
		userID := &user.Sub

		experienceIDStr, err := request.RequireString("experience_id")
		if err != nil {
			return nil, fmt.Errorf("experience_id parameter is required: %w", err)
		}

		experienceID, err := strconv.ParseUint(experienceIDStr, 10, 32)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid experience_id: %v", err)), nil
		}

		key, err := request.RequireString("key")
		if err != nil {
			return nil, fmt.Errorf("key parameter is required: %w", err)
		}

		value, err := request.RequireString("value")
		if err != nil {
			return nil, fmt.Errorf("value parameter is required: %w", err)
		}

		category, err := request.RequireString("category")
		if err != nil {
			return nil, fmt.Errorf("category parameter is required: %w", err)
		}

		featureMap := &models.FeatureMap{
			ExperienceID: uint(experienceID),
			Key:          key,
			Value:        value,
			Category:     category,
		}

		if err := db.AddFeatureMap(featureMap, userID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error adding feature map: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Feature map added successfully")), nil
	}

	return tool, handler
}
