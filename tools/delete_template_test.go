package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestDeleteTemplateTool_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	resume := createTestResume(t, db)
	template := createTestTemplate(t, db, resume.ID)

	tool, handler := NewDeleteTemplateTool(db)

	// Test tool creation
	if tool.Name != "delete_template" {
		t.Errorf("Expected tool name 'delete_template', got %s", tool.Name)
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

	if !strings.Contains(textContent.Text, "Template deleted successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	if !strings.Contains(textContent.Text, template.Name) {
		t.Errorf("Expected template name in response, got: %s", textContent.Text)
	}

	// Verify template was deleted from database
	_, err = db.GetTemplateByID(template.ID)
	if err == nil {
		t.Errorf("Expected template to be deleted, but it still exists")
	}
}

func TestDeleteTemplateTool_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewDeleteTemplateTool(db)

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

func TestDeleteTemplateTool_InvalidTemplateID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewDeleteTemplateTool(db)

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

func TestDeleteTemplateTool_MissingTemplateID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewDeleteTemplateTool(db)

	request := createTestRequest(map[string]interface{}{
		// Missing template_id
	})

	_, err := handler(context.Background(), request)
	if err == nil {
		t.Errorf("Expected error for missing template_id")
	}
}

func TestDeleteTemplateTool_MultipleDeletions(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	resume := createTestResume(t, db)
	template1 := createTestTemplate(t, db, resume.ID)
	template2 := createTestTemplate(t, db, resume.ID)
	template2.Name = "Second Template"
	db.UpdateTemplate(template2)

	_, handler := NewDeleteTemplateTool(db)

	// Delete first template
	request1 := createTestRequest(map[string]interface{}{
		"template_id": "1",
	})

	result1, err := handler(context.Background(), request1)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result1 == nil {
		t.Fatal("Handler returned nil result")
	}

	// Verify first template was deleted
	_, err = db.GetTemplateByID(template1.ID)
	if err == nil {
		t.Errorf("Expected first template to be deleted")
	}

	// Verify second template still exists
	_, err = db.GetTemplateByID(template2.ID)
	if err != nil {
		t.Errorf("Expected second template to still exist: %v", err)
	}

	// Delete second template
	request2 := createTestRequest(map[string]interface{}{
		"template_id": "2",
	})

	result2, err := handler(context.Background(), request2)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result2 == nil {
		t.Fatal("Handler returned nil result")
	}

	// Verify second template was also deleted
	_, err = db.GetTemplateByID(template2.ID)
	if err == nil {
		t.Errorf("Expected second template to be deleted")
	}

	// Verify no templates remain for the resume
	templates, err := db.ListTemplatesByResumeID(resume.ID)
	if err != nil {
		t.Fatalf("Failed to list templates: %v", err)
	}

	if len(templates) != 0 {
		t.Errorf("Expected no templates to remain, got %d", len(templates))
	}
}