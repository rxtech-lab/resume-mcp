package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/invopop/jsonschema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/models"
)

func NewGetResumeContextTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("get_resume_context",
		mcp.WithDescription(`Get JSON schema for resume data structure to help AI understand how to draft templates.

This tool returns the complete JSON schema for the resume data model without actual resume data.
Use this tool before creating templates to understand what fields are available and their types.

The schema includes all resume fields, relationships, and nested structures:
- Resume basic info (name, photo, description)
- Contact information array
- Work experiences array with feature maps
- Education entries array with feature maps  
- Other experiences array with feature maps
- Feature maps for flexible JSON data storage

No actual resume data is returned - only the schema structure.`),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("ID of the resume (used for validation, but actual data is not returned)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resumeIDStr, err := request.RequireString("resume_id")
		if err != nil {
			return nil, fmt.Errorf("resume_id parameter is required: %w", err)
		}

		resumeID, err := strconv.Atoi(resumeIDStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resume_id: %v", err)), nil
		}

		// Validate that resume exists (but we don't return the actual data)
		_, err = db.GetResumeByID(uint(resumeID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Resume not found: %v", err)), nil
		}

		// Generate JSON schema for Resume model
		reflector := jsonschema.Reflector{
			AllowAdditionalProperties: false,
			DoNotReference:            true,
		}
		schema := reflector.Reflect(&models.Resume{})

		// Convert schema to JSON
		schemaJSON, err := json.MarshalIndent(schema, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to generate schema: %v", err)), nil
		}

		// Create context response with JSON schema
		contextData := map[string]interface{}{
			"json_schema": json.RawMessage(schemaJSON),
		}

		result := map[string]interface{}{
			"success": true,
			"message": "Resume JSON schema retrieved successfully",
			"context": contextData,
		}

		resultJSON, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Resume JSON schema retrieved successfully, and please return the following schema in the response: "),
				mcp.NewTextContent(string(resultJSON)),
			},
		}, nil
	}

	return tool, handler
}
