package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/models"
	"github.com/rxtech-lab/resume-mcp/internal/service"
)

func setupTestDB(t *testing.T) *database.Database {
	db, err := database.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

func createTestResume(t *testing.T, db *database.Database) *models.Resume {
	resume := &models.Resume{
		Name:        "Test User",
		Photo:       "test.jpg",
		Description: "Test Description",
	}
	err := db.CreateResume(resume)
	if err != nil {
		t.Fatalf("Failed to create test resume: %v", err)
	}
	return resume
}

// createTestRequest creates a proper CallToolRequest for testing
func createTestRequest(arguments map[string]interface{}) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: arguments,
		},
	}
}

func TestCreateTemplateTool_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()
	resume := createTestResume(t, db)

	tool, handler := NewCreateTemplateTool(db, templateService)

	// Test tool creation
	if tool.Name != "create_template" {
		t.Errorf("Expected tool name 'create_template', got %s", tool.Name)
	}

	// Create request
	request := createTestRequest(map[string]interface{}{
		"resume_id":     "1",
		"name":          "Test Template",
		"description":   "A test template",
		"template_data": "<h1>{{.Name}}</h1><p>{{.Description}}</p>",
	})

	// Execute handler
	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Handler returned nil result")
	}

	// Check result content
	if len(result.Content) == 0 {
		t.Fatal("Handler returned empty content")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}
	if !strings.Contains(textContent.Text, "Created template successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	// Verify template was created in database
	templates, err := db.ListTemplatesByResumeID(resume.ID)
	if err != nil {
		t.Fatalf("Failed to list templates: %v", err)
	}

	if len(templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(templates))
	}

	template := templates[0]
	if template.Name != "Test Template" {
		t.Errorf("Expected template name 'Test Template', got %s", template.Name)
	}
	if template.Description != "A test template" {
		t.Errorf("Expected description 'A test template', got %s", template.Description)
	}
	if template.ResumeID != resume.ID {
		t.Errorf("Expected resume ID %d, got %d", resume.ID, template.ResumeID)
	}
}

func TestCreateTemplateTool_InvalidResumeID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()

	_, handler := NewCreateTemplateTool(db, templateService)

	request := createTestRequest(map[string]interface{}{
		"resume_id":     "999", // Non-existent resume
		"name":          "Test Template",
		"template_data": "<h1>{{.Name}}</h1>",
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	// Should return error result
	if len(result.Content) == 0 {
		t.Fatal("Expected error content")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}
	if !strings.Contains(textContent.Text, "Resume not found") {
		t.Errorf("Expected 'Resume not found' error, got: %s", textContent.Text)
	}
}

func TestCreateTemplateTool_InvalidTemplate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()
	createTestResume(t, db)

	_, handler := NewCreateTemplateTool(db, templateService)

	request := createTestRequest(map[string]interface{}{
		"resume_id":     "1",
		"name":          "Invalid Template",
		"template_data": "{{.NonExistentField}}", // This should cause a validation error
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	// Should return validation error
	if len(result.Content) == 0 {
		t.Fatal("Expected error content")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}
	if !strings.Contains(textContent.Text, "Template validation failed") {
		t.Errorf("Expected template validation error, got: %s", textContent.Text)
	}
}

func TestCreateTemplateTool_MissingRequiredFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()

	_, handler := NewCreateTemplateTool(db, templateService)

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "missing resume_id",
			args: map[string]interface{}{
				"name":          "Test Template",
				"template_data": "<h1>Test</h1>",
			},
		},
		{
			name: "missing name",
			args: map[string]interface{}{
				"resume_id":     "1",
				"template_data": "<h1>Test</h1>",
			},
		},
		{
			name: "missing template_data",
			args: map[string]interface{}{
				"resume_id": "1",
				"name":      "Test Template",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := createTestRequest(tt.args)

			_, err := handler(context.Background(), request)
			if err == nil {
				t.Errorf("Expected error for missing required field")
			}
		})
	}
}

func TestCreateTemplateTool_WithOptionalDescription(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()
	createTestResume(t, db)

	_, handler := NewCreateTemplateTool(db, templateService)

	request := createTestRequest(map[string]interface{}{
		"resume_id":     "1",
		"name":          "Template Without Description",
		"template_data": "<h1>{{.Name}}</h1>",
		// No description provided
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Handler returned nil result")
	}

	// Verify template was created with empty description
	templates, err := db.ListTemplatesByResumeID(1)
	if err != nil {
		t.Fatalf("Failed to list templates: %v", err)
	}

	if len(templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(templates))
	}

	template := templates[0]
	if template.Description != "" {
		t.Errorf("Expected empty description, got %s", template.Description)
	}
}

func TestCreateTemplateTool_CopyFromExistingResume(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()
	
	// Create source resume with full data
	_ = createFullTestResume(t, db)
	
	// Create target resume (empty)
	targetResume := createTestResume(t, db)

	_, handler := NewCreateTemplateTool(db, templateService)

	request := createTestRequest(map[string]interface{}{
		"resume_id":             "2", // Target resume ID
		"copy_from_resume_id":   "1", // Source resume ID
		"name":                  "Copied Template",
		"description":           "Template with copied data",
		"template_data":         "<h1>{{.Name}}</h1>{{range .WorkExperiences}}<div>{{.JobTitle}} at {{.Company}}</div>{{end}}",
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Handler returned nil result")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	if !strings.Contains(textContent.Text, "Created template successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	if !strings.Contains(textContent.Text, "copied_from_resume_id") {
		t.Errorf("Expected copied_from_resume_id in response, got: %s", textContent.Text)
	}

	// Verify that target resume now has the copied data
	fullTargetResume, err := db.GetResumeByID(targetResume.ID)
	if err != nil {
		t.Fatalf("Failed to get target resume: %v", err)
	}

	// Check that work experiences were copied
	if len(fullTargetResume.WorkExperiences) != 1 {
		t.Errorf("Expected 1 work experience, got %d", len(fullTargetResume.WorkExperiences))
	}

	// Check that education was copied
	if len(fullTargetResume.Educations) != 1 {
		t.Errorf("Expected 1 education, got %d", len(fullTargetResume.Educations))
	}

	// Check that other experiences were copied
	if len(fullTargetResume.OtherExperiences) != 1 {
		t.Errorf("Expected 1 other experience, got %d", len(fullTargetResume.OtherExperiences))
	}

	// Verify feature maps were copied for work experience
	if len(fullTargetResume.WorkExperiences) > 0 && len(fullTargetResume.WorkExperiences[0].FeatureMaps) < 1 {
		t.Errorf("Expected at least 1 work experience feature map, got %d", len(fullTargetResume.WorkExperiences[0].FeatureMaps))
	}

	// Verify feature maps were copied for education
	if len(fullTargetResume.Educations) > 0 && len(fullTargetResume.Educations[0].FeatureMaps) < 1 {
		t.Errorf("Expected at least 1 education feature map, got %d", len(fullTargetResume.Educations[0].FeatureMaps))
	}

	// Verify feature maps were copied for other experiences
	if len(fullTargetResume.OtherExperiences) > 0 && len(fullTargetResume.OtherExperiences[0].FeatureMaps) < 1 {
		t.Errorf("Expected at least 1 other experience feature map, got %d", len(fullTargetResume.OtherExperiences[0].FeatureMaps))
	}

	// Verify template was created for target resume
	templates, err := db.ListTemplatesByResumeID(targetResume.ID)
	if err != nil {
		t.Fatalf("Failed to list templates: %v", err)
	}

	if len(templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(templates))
	}

	template := templates[0]
	if template.Name != "Copied Template" {
		t.Errorf("Expected template name 'Copied Template', got %s", template.Name)
	}
}

func TestCreateTemplateTool_CopyFromNonExistentResume(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()
	createTestResume(t, db)

	_, handler := NewCreateTemplateTool(db, templateService)

	request := createTestRequest(map[string]interface{}{
		"resume_id":           "1",
		"copy_from_resume_id": "999", // Non-existent resume
		"name":                "Test Template",
		"template_data":       "<h1>{{.Name}}</h1>",
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	if !strings.Contains(textContent.Text, "Source resume not found") {
		t.Errorf("Expected 'Source resume not found' error, got: %s", textContent.Text)
	}
}

func TestCreateTemplateTool_CopyFromInvalidID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()
	createTestResume(t, db)

	_, handler := NewCreateTemplateTool(db, templateService)

	request := createTestRequest(map[string]interface{}{
		"resume_id":           "1",
		"copy_from_resume_id": "invalid", // Invalid ID format
		"name":                "Test Template",
		"template_data":       "<h1>{{.Name}}</h1>",
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	if !strings.Contains(textContent.Text, "Invalid copy_from_resume_id") {
		t.Errorf("Expected 'Invalid copy_from_resume_id' error, got: %s", textContent.Text)
	}
}

func TestCreateTemplateTool_CopyFromEmptyResume(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()
	
	// Create source resume (empty)
	createTestResume(t, db) // ID 1
	
	// Create target resume 
	createTestResume(t, db) // ID 2

	_, handler := NewCreateTemplateTool(db, templateService)

	request := createTestRequest(map[string]interface{}{
		"resume_id":           "2",
		"copy_from_resume_id": "1", // Empty source resume
		"name":                "Template from Empty",
		"template_data":       "<h1>{{.Name}}</h1>",
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Handler returned nil result")
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	if !strings.Contains(textContent.Text, "Created template successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	// Verify template was created
	templates, err := db.ListTemplatesByResumeID(2)
	if err != nil {
		t.Fatalf("Failed to list templates: %v", err)
	}

	if len(templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(templates))
	}
}