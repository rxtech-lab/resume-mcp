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
		mcp.WithDescription("Generate HTML preview of a resume using Go templates. Returns a preview URL. Templates now include Tailwind CSS for styling."),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("The ID of the resume to generate preview for"),
		),
		mcp.WithString("template",
			mcp.Required(),
			mcp.Description("Go template string for rendering. Use Tailwind CSS classes for styling (CDN included). Call get_resume_by_name to see available context variables."),
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

		template, err := request.RequireString("template")
		if err != nil {
			return nil, fmt.Errorf("template parameter is required: %w", err)
		}

		css := request.GetString("css", "")

		resume, err := db.GetResumeByID(uint(resumeID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error getting resume: %v", err)), nil
		}

		_, err = templateService.GeneratePreview(template, css, *resume)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error generating preview: %v", err)), nil
		}

		sessionID, err := db.GeneratePreview(uint(resumeID), template, css)
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
