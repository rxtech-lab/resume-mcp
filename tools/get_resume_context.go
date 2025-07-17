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

func NewGetResumeContextTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("get_resume_context",
		mcp.WithDescription(`Get all available context and data structure for a resume to help AI understand how to draft templates.

This tool returns the complete resume data with all relationships and provides a schema guide for template creation.
Use this tool before creating templates to understand what fields are available and how they're structured.

Template Context Fields Available:
- .Name (string): Resume owner's name
- .Photo (string): Photo URL or path
- .Description (string): Resume description/summary
- .Contacts ([]Contact): Array of contact information
  - .Key (string): Contact type (e.g., "email", "phone", "linkedin")
  - .Value (string): Contact value
- .WorkExperiences ([]WorkExperience): Array of work experiences
  - .Company (string): Company name
  - .JobTitle (string): Job title
  - .StartDate (time.Time): Start date (use .Format "Jan 2006")
  - .EndDate (*time.Time): End date (nullable, use if .EndDate check)
  - .FeatureMaps ([]FeatureMap): Additional work details
- .Educations ([]Education): Array of education entries
  - .SchoolName (string): School name
  - .StartDate (time.Time): Start date
  - .EndDate (*time.Time): End date (nullable)
  - .FeatureMaps ([]FeatureMap): Additional education details (GPA, degree, etc.)
- .OtherExperiences ([]OtherExperience): Array of other experiences
  - .Category (string): Experience category
  - .FeatureMaps ([]FeatureMap): Experience details
- FeatureMap structure:
  - .Key (string): Feature name
  - .Value (string): Feature value (can be JSON for complex data)`),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("ID of the resume to get context for"),
		),
		mcp.WithString("include_examples",
			mcp.Description("Include template examples in the response (true/false, default: true)"),
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

		includeExamples := request.GetString("include_examples", "true")

		// Get resume with all related data
		resume, err := db.GetResumeByID(uint(resumeID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Resume not found: %v", err)), nil
		}

		// Get templates for this resume
		templates, err := db.ListTemplatesByResumeID(uint(resumeID))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get templates: %v", err)), nil
		}

		// Create context response
		contextData := map[string]interface{}{
			"resume_data": resume,
			"templates":   templates,
			"schema_guide": map[string]interface{}{
				"basic_fields": map[string]string{
					"Name":        "string - Resume owner's name",
					"Photo":       "string - Photo URL or path",
					"Description": "string - Resume description/summary",
				},
				"contact_fields": map[string]string{
					"Key":   "string - Contact type (email, phone, linkedin, etc.)",
					"Value": "string - Contact value",
				},
				"work_experience_fields": map[string]string{
					"Company":   "string - Company name",
					"JobTitle":  "string - Job title",
					"StartDate": "time.Time - Start date (use .Format \"Jan 2006\")",
					"EndDate":   "*time.Time - End date (nullable, check with if .EndDate)",
				},
				"education_fields": map[string]string{
					"SchoolName": "string - School name",
					"StartDate":  "time.Time - Start date",
					"EndDate":    "*time.Time - End date (nullable)",
				},
				"other_experience_fields": map[string]string{
					"Category": "string - Experience category",
				},
				"feature_map_fields": map[string]string{
					"Key":   "string - Feature name",
					"Value": "string - Feature value (can be JSON for complex data)",
				},
			},
			"data_counts": map[string]int{
				"contacts":          len(resume.Contacts),
				"work_experiences":  len(resume.WorkExperiences),
				"educations":        len(resume.Educations),
				"other_experiences": len(resume.OtherExperiences),
				"templates":         len(templates),
			},
		}

		if includeExamples == "true" {
			contextData["template_examples"] = map[string]interface{}{
				"basic_template": `<div class="resume">
  <h1>{{.Name}}</h1>
  <p>{{.Description}}</p>
</div>`,
				"contact_template": `{{if .Contacts}}
<div class="contact-section">
  <h2>Contact Information</h2>
  {{range .Contacts}}
  <p><strong>{{.Key}}:</strong> {{.Value}}</p>
  {{end}}
</div>
{{end}}`,
				"work_experience_template": `{{if .WorkExperiences}}
<div class="work-section">
  <h2>Work Experience</h2>
  {{range .WorkExperiences}}
  <div class="work-item">
    <h3>{{.JobTitle}} at {{.Company}}</h3>
    <p class="dates">{{.StartDate.Format "Jan 2006"}} - {{if .EndDate}}{{.EndDate.Format "Jan 2006"}}{{else}}Present{{end}}</p>
    {{if .FeatureMaps}}
    <ul>
      {{range .FeatureMaps}}
      <li><strong>{{.Key}}:</strong> {{.Value}}</li>
      {{end}}
    </ul>
    {{end}}
  </div>
  {{end}}
</div>
{{end}}`,
				"education_template": `{{if .Educations}}
<div class="education-section">
  <h2>Education</h2>
  {{range .Educations}}
  <div class="education-item">
    <h3>{{.SchoolName}}</h3>
    <p class="dates">{{.StartDate.Format "Jan 2006"}} - {{if .EndDate}}{{.EndDate.Format "Jan 2006"}}{{else}}Present{{end}}</p>
    {{if .FeatureMaps}}
    <ul>
      {{range .FeatureMaps}}
      <li><strong>{{.Key}}:</strong> {{.Value}}</li>
      {{end}}
    </ul>
    {{end}}
  </div>
  {{end}}
</div>
{{end}}`,
				"comprehensive_template": `<div class="max-w-4xl mx-auto p-8 bg-white">
  <div class="header mb-8">
    {{if .Photo}}<img src="{{.Photo}}" alt="{{.Name}}" class="w-24 h-24 rounded-full mb-4">{{end}}
    <h1 class="text-3xl font-bold text-gray-800">{{.Name}}</h1>
    <p class="text-gray-600 mt-2">{{.Description}}</p>
  </div>

  {{if .Contacts}}
  <div class="contact-section mb-6">
    <h2 class="text-xl font-semibold text-gray-700 mb-3">Contact</h2>
    <div class="grid grid-cols-2 gap-2">
      {{range .Contacts}}
      <p><span class="font-medium">{{.Key}}:</span> {{.Value}}</p>
      {{end}}
    </div>
  </div>
  {{end}}

  {{if .WorkExperiences}}
  <div class="work-section mb-6">
    <h2 class="text-xl font-semibold text-gray-700 mb-3">Work Experience</h2>
    {{range .WorkExperiences}}
    <div class="work-item mb-4 p-4 border-l-4 border-blue-500">
      <h3 class="font-semibold text-lg">{{.JobTitle}} at {{.Company}}</h3>
      <p class="text-sm text-gray-600 mb-2">{{.StartDate.Format "Jan 2006"}} - {{if .EndDate}}{{.EndDate.Format "Jan 2006"}}{{else}}Present{{end}}</p>
      {{if .FeatureMaps}}
      <ul class="list-disc list-inside text-sm">
        {{range .FeatureMaps}}
        <li><strong>{{.Key}}:</strong> {{.Value}}</li>
        {{end}}
      </ul>
      {{end}}
    </div>
    {{end}}
  </div>
  {{end}}

  {{if .Educations}}
  <div class="education-section mb-6">
    <h2 class="text-xl font-semibold text-gray-700 mb-3">Education</h2>
    {{range .Educations}}
    <div class="education-item mb-4 p-4 border-l-4 border-green-500">
      <h3 class="font-semibold text-lg">{{.SchoolName}}</h3>
      <p class="text-sm text-gray-600 mb-2">{{.StartDate.Format "Jan 2006"}} - {{if .EndDate}}{{.EndDate.Format "Jan 2006"}}{{else}}Present{{end}}</p>
      {{if .FeatureMaps}}
      <ul class="list-disc list-inside text-sm">
        {{range .FeatureMaps}}
        <li><strong>{{.Key}}:</strong> {{.Value}}</li>
        {{end}}
      </ul>
      {{end}}
    </div>
    {{end}}
  </div>
  {{end}}

  {{if .OtherExperiences}}
  <div class="other-section mb-6">
    <h2 class="text-xl font-semibold text-gray-700 mb-3">Other Experience</h2>
    {{range .OtherExperiences}}
    <div class="other-item mb-4 p-4 border-l-4 border-purple-500">
      <h3 class="font-semibold text-lg">{{.Category}}</h3>
      {{if .FeatureMaps}}
      <ul class="list-disc list-inside text-sm">
        {{range .FeatureMaps}}
        <li><strong>{{.Key}}:</strong> {{.Value}}</li>
        {{end}}
      </ul>
      {{end}}
    </div>
    {{end}}
  </div>
  {{end}}
</div>`,
			}
		}

		result := map[string]interface{}{
			"success": true,
			"message": "Resume context retrieved successfully",
			"context": contextData,
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Resume context: %s", string(resultJSON))), nil
	}

	return tool, handler
}