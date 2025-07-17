package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/models"
)

func createTestTemplate(t *testing.T, db *database.Database, resumeID uint) *models.Template {
	template := &models.Template{
		ResumeID:     resumeID,
		Name:         "Test Template",
		Description:  "A test template",
		TemplateData: "<h1>{{.Name}}</h1>",
	}
	err := db.CreateTemplate(template)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}
	return template
}

func TestGetTemplateTool_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	resume := createTestResume(t, db)
	_ = createTestTemplate(t, db, resume.ID)

	tool, handler := NewGetTemplateTool(db)

	// Test tool creation
	if tool.Name != "get_template" {
		t.Errorf("Expected tool name 'get_template', got %s", tool.Name)
	}

	// Create request
	request := createTestRequest(map[string]interface{}{
		"template_id": "1",
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

	if !strings.Contains(textContent.Text, "Template retrieved successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	// Check that the response contains template information
	if !strings.Contains(textContent.Text, "Test Template") {
		t.Errorf("Expected template name in response, got: %s", textContent.Text)
	}
	
	if !strings.Contains(textContent.Text, "{{.Name}}") {
		t.Errorf("Expected template data in response, got: %s", textContent.Text)
	}
}

func TestGetTemplateTool_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewGetTemplateTool(db)

	request := createTestRequest(map[string]interface{}{
		"template_id": "999", // Non-existent template
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

	if !strings.Contains(textContent.Text, "Template not found") {
		t.Errorf("Expected 'Template not found' error, got: %s", textContent.Text)
	}
}

func TestGetTemplateTool_InvalidTemplateID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewGetTemplateTool(db)

	request := createTestRequest(map[string]interface{}{
		"template_id": "invalid", // Invalid template ID
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

	if !strings.Contains(textContent.Text, "Invalid template_id") {
		t.Errorf("Expected 'Invalid template_id' error, got: %s", textContent.Text)
	}
}

func TestGetTemplateTool_MissingTemplateID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewGetTemplateTool(db)

	request := createTestRequest(map[string]interface{}{
		// Missing template_id
	})

	_, err := handler(context.Background(), request)
	if err == nil {
		t.Errorf("Expected error for missing template_id")
	}
}

func TestGetTemplateTool_WithAllFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	resume := createTestResume(t, db)
	template := &models.Template{
		ResumeID:     resume.ID,
		Name:         "Comprehensive Template",
		Description:  "A template with all fields filled",
		TemplateData: `<div class="resume">
			<h1>{{.Name}}</h1>
			<p>{{.Description}}</p>
			{{range .Contacts}}
			<span>{{.Key}}: {{.Value}}</span>
			{{end}}
		</div>`,
	}
	err := db.CreateTemplate(template)
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	_, handler := NewGetTemplateTool(db)

	request := createTestRequest(map[string]interface{}{
		"template_id": "1",
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

	// Verify all template fields are in the response
	if !strings.Contains(textContent.Text, "Comprehensive Template") {
		t.Errorf("Expected template name in response")
	}
	if !strings.Contains(textContent.Text, "A template with all fields filled") {
		t.Errorf("Expected template description in response")
	}
	if !strings.Contains(textContent.Text, "class=\\\"resume\\\"") {
		t.Errorf("Expected template data in response")
	}
}