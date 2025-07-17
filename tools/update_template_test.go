package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rxtech-lab/resume-mcp/internal/service"
)

func TestUpdateTemplateTool_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()
	resume := createTestResume(t, db)
	template := createTestTemplate(t, db, resume.ID)

	tool, handler := NewUpdateTemplateTool(db, templateService)

	// Test tool creation
	if tool.Name != "update_template" {
		t.Errorf("Expected tool name 'update_template', got %s", tool.Name)
	}

	// Create request to update all fields
	request := createTestRequest(map[string]interface{}{
		"template_id":   "1",
		"name":          "Updated Template",
		"description":   "Updated description",
		"template_data": "<h2>{{.Name}}</h2><p>Updated: {{.Description}}</p>",
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

	if !strings.Contains(textContent.Text, "Template updated successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	// Verify template was updated in database
	updatedTemplate, err := db.GetTemplateByID(template.ID)
	if err != nil {
		t.Fatalf("Failed to get updated template: %v", err)
	}

	if updatedTemplate.Name != "Updated Template" {
		t.Errorf("Expected updated name 'Updated Template', got %s", updatedTemplate.Name)
	}
	if updatedTemplate.Description != "Updated description" {
		t.Errorf("Expected updated description 'Updated description', got %s", updatedTemplate.Description)
	}
	if !strings.Contains(updatedTemplate.TemplateData, "<h2>{{.Name}}</h2>") {
		t.Errorf("Expected updated template data, got %s", updatedTemplate.TemplateData)
	}
}

func TestUpdateTemplateTool_PartialUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()
	resume := createTestResume(t, db)
	template := createTestTemplate(t, db, resume.ID)

	_, handler := NewUpdateTemplateTool(db, templateService)

	// Update only the name
	request := createTestRequest(map[string]interface{}{
		"template_id": "1",
		"name":        "Only Name Updated",
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Handler returned nil result")
	}

	// Verify only name was updated
	updatedTemplate, err := db.GetTemplateByID(template.ID)
	if err != nil {
		t.Fatalf("Failed to get updated template: %v", err)
	}

	if updatedTemplate.Name != "Only Name Updated" {
		t.Errorf("Expected updated name 'Only Name Updated', got %s", updatedTemplate.Name)
	}
	// Other fields should remain the same
	if updatedTemplate.Description != template.Description {
		t.Errorf("Description should not have changed")
	}
	if updatedTemplate.TemplateData != template.TemplateData {
		t.Errorf("Template data should not have changed")
	}
}

func TestUpdateTemplateTool_InvalidTemplate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()
	resume := createTestResume(t, db)
	createTestTemplate(t, db, resume.ID)

	_, handler := NewUpdateTemplateTool(db, templateService)

	request := createTestRequest(map[string]interface{}{
		"template_id":   "1",
		"template_data": "{{.InvalidField}}", // Invalid template
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

func TestUpdateTemplateTool_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()

	_, handler := NewUpdateTemplateTool(db, templateService)

	request := createTestRequest(map[string]interface{}{
		"template_id": "999", // Non-existent template
		"name":        "Updated Name",
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

func TestUpdateTemplateTool_InvalidTemplateID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()

	_, handler := NewUpdateTemplateTool(db, templateService)

	request := createTestRequest(map[string]interface{}{
		"template_id": "invalid", // Invalid template ID
		"name":        "Updated Name",
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	if !strings.Contains(textContent.Text, "Invalid template_id") {
		t.Errorf("Expected 'Invalid template_id' error, got: %s", textContent.Text)
	}
}

func TestUpdateTemplateTool_MissingTemplateID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()

	_, handler := NewUpdateTemplateTool(db, templateService)

	request := createTestRequest(map[string]interface{}{
		// Missing template_id
		"name": "Updated Name",
	})

	_, err := handler(context.Background(), request)
	if err == nil {
		t.Errorf("Expected error for missing template_id")
	}
}

func TestUpdateTemplateTool_EmptyUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	templateService := service.NewTemplateService()
	resume := createTestResume(t, db)
	template := createTestTemplate(t, db, resume.ID)

	_, handler := NewUpdateTemplateTool(db, templateService)

	// Request with no optional fields
	request := createTestRequest(map[string]interface{}{
		"template_id": "1",
		// No other fields to update
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Handler returned nil result")
	}

	// Verify template was not changed
	unchangedTemplate, err := db.GetTemplateByID(template.ID)
	if err != nil {
		t.Fatalf("Failed to get template: %v", err)
	}

	if unchangedTemplate.Name != template.Name {
		t.Errorf("Name should not have changed")
	}
	if unchangedTemplate.Description != template.Description {
		t.Errorf("Description should not have changed")
	}
	if unchangedTemplate.TemplateData != template.TemplateData {
		t.Errorf("Template data should not have changed")
	}
}