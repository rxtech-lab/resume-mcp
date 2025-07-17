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
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

		featureMap := &models.FeatureMap{
			ExperienceID: uint(experienceID),
			Key:          key,
			Value:        value,
		}

		if err := db.AddFeatureMap(featureMap); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error adding feature map: %v", err)), nil
		}

		result := map[string]interface{}{
			"id":            featureMap.ID,
			"experience_id": featureMap.ExperienceID,
			"key":           featureMap.Key,
			"value":         featureMap.Value,
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Feature map added successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}