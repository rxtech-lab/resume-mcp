package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
)

func NewDeleteFeatureMapTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("delete_feature_map",
		mcp.WithDescription("Delete a specific feature map by ID. Use this to remove specific details from experiences."),
		mcp.WithString("feature_map_id",
			mcp.Required(),
			mcp.Description("The ID of the feature map to delete"),
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

		if err := db.DeleteFeatureMap(uint(featureMapID)); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error deleting feature map: %v", err)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Feature map with ID %d deleted successfully", featureMapID)), nil
	}

	return tool, handler
}