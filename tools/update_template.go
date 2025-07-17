package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/service"
)

func NewUpdateTemplateTool(db *database.Database, templateService *service.TemplateService) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("update_template",
		mcp.WithDescription("Update an existing template. If user ask to remove sections or data from the final preview, please use this tool to update the template and don't try to delete the data first. Only delete the data if you are sure about the data is not needed."),
		mcp.WithString("template_id",
			mcp.Required(),
			mcp.Description("ID of the template to update"),
		),
		mcp.WithString("name",
			mcp.Description("New name for the template (optional)"),
		),
		mcp.WithString("description",
			mcp.Description("New description for the template (optional)"),
		),
		mcp.WithString("template_data",
			mcp.Description("New template data (optional)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		templateIDStr, err := request.RequireString("template_id")
		if err != nil {
			return nil, fmt.Errorf("template_id parameter is required: %w", err)
		}

		templateID, err := strconv.Atoi(templateIDStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid template_id: %v", err)), nil
		}

		template, err := db.GetTemplateByID(uint(templateID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Template not found: %v", err)), nil
		}

		// Update fields if provided
		name := request.GetString("name", "")
		if name != "" {
			template.Name = name
		}

		description := request.GetString("description", "")
		if description != "" {
			template.Description = description
		}

		templateData := request.GetString("template_data", "")
		if templateData != "" {
			// Validate new template by testing it
			resume, err := db.GetResumeByID(template.ResumeID)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Resume not found: %v", err)), nil
			}

			_, err = templateService.GeneratePreview(templateData, "", *resume)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Template validation failed: %v. Please check your Go template syntax and ensure all referenced fields exist on the resume model.", err)), nil
			}

			template.TemplateData = templateData
		}

		if err := db.UpdateTemplate(template); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update template: %v", err)), nil
		}

		result := map[string]interface{}{
			"success": true,
			"message": fmt.Sprintf("Template '%s' updated successfully", template.Name),
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Template updated successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}
