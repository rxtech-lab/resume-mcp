package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/service"
)

func NewGeneratePreviewTool(db *database.Database, port string, templateService *service.TemplateService) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("generate_preview",
		mcp.WithDescription("Generate HTML preview of a resume using a saved template. Returns a preview URL. Templates include Tailwind CSS for styling."),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("The ID of the resume to generate preview for"),
		),
		mcp.WithString("template_id",
			mcp.Required(),
			mcp.Description("The ID of the template to use for rendering"),
		),
		mcp.WithString("css",
			mcp.Description("Additional CSS styles for the preview (optional, Tailwind CSS classes are available in templates)"),
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

		templateIDStr, err := request.RequireString("template_id")
		if err != nil {
			return nil, fmt.Errorf("template_id parameter is required: %w", err)
		}

		templateID, err := strconv.ParseUint(templateIDStr, 10, 32)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid template_id: %v", err)), nil
		}

		css := request.GetString("css", "")

		resume, err := db.GetResumeByID(uint(resumeID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error getting resume: %v", err)), nil
		}

		template, err := db.GetTemplateByID(uint(templateID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error getting template: %v", err)), nil
		}

		// Verify template belongs to the same resume
		if template.ResumeID != uint(resumeID) {
			return mcp.NewToolResultError("Template does not belong to the specified resume"), nil
		}

		_, err = templateService.GeneratePreview(template.TemplateData, css, *resume)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error generating preview: %v", err)), nil
		}

		sessionID, err := db.GeneratePreview(uint(resumeID), template.TemplateData, css)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error generating preview: %v", err)), nil
		}

		previewURL := fmt.Sprintf("http://localhost:%s/resume/preview/%s", port, sessionID)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Preview generated successfully, and please return the following URL in the response: "),
				mcp.NewTextContent(previewURL),
			},
		}, nil
	}

	return tool, handler
}
