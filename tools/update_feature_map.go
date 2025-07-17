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

func NewUpdateFeatureMapTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("update_feature_map",
		mcp.WithDescription("Update an existing feature map by ID. Use this to modify specific details attached to experiences."),
		mcp.WithString("feature_map_id",
			mcp.Required(),
			mcp.Description("The ID of the feature map to update"),
		),
		mcp.WithString("key",
			mcp.Description("The feature key"),
		),
		mcp.WithString("value",
			mcp.Description("The feature value"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		featureMapIDStr, err := request.RequireString("feature_map_id")
		if err != nil {
			return nil, fmt.Errorf("feature_map_id parameter is required: %w", err)
		}

		featureMapID, err := strconv.ParseUint(featureMapIDStr, 10, 32)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid feature_map_id: %v", err)), nil
		}

		key := request.GetString("key", "")
		value := request.GetString("value", "")

		featureMap, err := db.GetFeatureMapByID(uint(featureMapID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Feature map not found: %v", err)), nil
		}

		if key != "" {
			featureMap.Key = key
		}
		if value != "" {
			featureMap.Value = value
		}

		if err := db.UpdateFeatureMap(featureMap); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error updating feature map: %v", err)), nil
		}

		result := map[string]interface{}{
			"id":            featureMap.ID,
			"experience_id": featureMap.ExperienceID,
			"key":           featureMap.Key,
			"value":         featureMap.Value,
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Feature map updated successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}