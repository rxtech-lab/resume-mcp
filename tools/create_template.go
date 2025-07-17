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
	"github.com/rxtech-lab/resume-mcp/internal/service"
)

func NewCreateTemplateTool(db *database.Database, templateService *service.TemplateService) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("create_template",
		mcp.WithDescription(`Create a new template for a resume. The template uses Go template syntax with access to resume data.

Example template:
<div class="max-w-4xl mx-auto p-8 bg-white">
  <h1 class="text-3xl font-bold text-gray-800">{{.Name}}</h1>
  <p class="text-gray-600 mt-2">{{.Description}}</p>
  
  {{if .Contacts}}
  <div class="mt-6">
    <h2 class="text-xl font-semibold text-gray-700">Contact</h2>
    {{range .Contacts}}
    <p>{{.Key}}: {{.Value}}</p>
    {{end}}
  </div>
  {{end}}
  
  {{if .WorkExperiences}}
  <div class="mt-6">
    <h2 class="text-xl font-semibold text-gray-700">Work Experience</h2>
    {{range .WorkExperiences}}
    <div class="mb-4">
      <h3 class="font-semibold">{{.JobTitle}} at {{.Company}}</h3>
      <p class="text-sm text-gray-600">{{.StartDate.Format "Jan 2006"}} - {{if .EndDate}}{{.EndDate.Format "Jan 2006"}}{{else}}Present{{end}}</p>
      {{range .FeatureMaps}}
      <p>{{.Key}}: {{.Value}}</p>
      {{end}}
    </div>
    {{end}}
  </div>
  {{end}}
</div>`),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("ID of the resume this template belongs to"),
		),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the template"),
		),
		mcp.WithString("description",
			mcp.Description("Description of what this template does"),
		),
		mcp.WithString("template_data",
			mcp.Required(),
			mcp.Description("Go template string for rendering the resume HTML"),
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

		name, err := request.RequireString("name")
		if err != nil {
			return nil, fmt.Errorf("name parameter is required: %w", err)
		}

		templateData, err := request.RequireString("template_data")
		if err != nil {
			return nil, fmt.Errorf("template_data parameter is required: %w", err)
		}

		description := request.GetString("description", "")

		// Validate resume exists
		resume, err := db.GetResumeByID(uint(resumeID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Resume not found: %v", err)), nil
		}

		// Validate template by testing it
		_, err = templateService.GeneratePreview(templateData, "", *resume)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Template validation failed: %v. Please check your Go template syntax and ensure all referenced fields exist on the resume model.", err)), nil
		}

		template := &models.Template{
			ResumeID:     uint(resumeID),
			Name:         name,
			Description:  description,
			TemplateData: templateData,
		}

		if err := db.CreateTemplate(template); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create template: %v", err)), nil
		}

		result := map[string]interface{}{
			"success":     true,
			"template_id": template.ID,
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Created template successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}
