package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/service"
	"github.com/rxtech-lab/resume-mcp/internal/types"
	"github.com/rxtech-lab/resume-mcp/internal/utils"
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

		templateIDStr, err := request.RequireString("template_id")
		if err != nil {
			return nil, fmt.Errorf("template_id parameter is required: %w", err)
		}

		templateID, err := strconv.ParseUint(templateIDStr, 10, 32)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid template_id: %v", err)), nil
		}

		css := request.GetString("css", "")

		resume, err := db.GetResumeByID(uint(resumeID), userID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error getting resume: %v", err)), nil
		}

		template, err := db.GetTemplateByID(uint(templateID), userID)
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

		sessionID, err := db.GeneratePreview(uint(resumeID), template.TemplateData, css, userID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error generating preview: %v", err)), nil
		}

		previewURL, err := utils.GetTransactionSessionUrl(port, sessionID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error generating preview: %v", err)), nil
		}

		downloadURL, err := utils.GetDownloadSessionUrl(port, sessionID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error generating download URL: %v", err)), nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent("Preview generated successfully, and please return the following URLs in the response:\n"),
				mcp.NewTextContent(fmt.Sprintf("Preview: %s\n", previewURL)),
				mcp.NewTextContent(fmt.Sprintf("Download PDF: %s", downloadURL)),
			},
		}, nil
	}

	return tool, handler
}
