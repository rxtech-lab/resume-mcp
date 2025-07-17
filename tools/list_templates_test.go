package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestListTemplatesTool_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	resume := createTestResume(t, db)
	
	// Create multiple templates
	createTestTemplate(t, db, resume.ID)
	template2 := createTestTemplate(t, db, resume.ID)
	template2.Name = "Second Template"
	template2.Description = "Another test template"
	db.UpdateTemplate(template2)

	tool, handler := NewListTemplatesTool(db)

	// Test tool creation
	if tool.Name != "list_templates" {
		t.Errorf("Expected tool name 'list_templates', got %s", tool.Name)
	}

	// Create request
	request := createTestRequest(map[string]interface{}{
		"resume_id": "1",
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

	if !strings.Contains(textContent.Text, "Templates listed successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	// Check that both templates are in the response
	if !strings.Contains(textContent.Text, "Test Template") {
		t.Errorf("Expected first template in response, got: %s", textContent.Text)
	}
	
	if !strings.Contains(textContent.Text, "Second Template") {
		t.Errorf("Expected second template in response, got: %s", textContent.Text)
	}

	// Check count
	if !strings.Contains(textContent.Text, "\"count\":2") {
		t.Errorf("Expected count of 2 templates, got: %s", textContent.Text)
	}
}

func TestListTemplatesTool_EmptyList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	_ = createTestResume(t, db) // Create resume but don't create any templates

	_, handler := NewListTemplatesTool(db)

	request := createTestRequest(map[string]interface{}{
		"resume_id": "1",
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

	if !strings.Contains(textContent.Text, "Templates listed successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	// Check count is 0
	if !strings.Contains(textContent.Text, "\"count\":0") {
		t.Errorf("Expected count of 0 templates, got: %s", textContent.Text)
	}
}

func TestListTemplatesTool_InvalidResumeID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewListTemplatesTool(db)

	request := createTestRequest(map[string]interface{}{
		"resume_id": "invalid",
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	if !strings.Contains(textContent.Text, "Invalid resume_id") {
		t.Errorf("Expected 'Invalid resume_id' error, got: %s", textContent.Text)
	}
}

func TestListTemplatesTool_MissingResumeID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewListTemplatesTool(db)

	request := createTestRequest(map[string]interface{}{
		// Missing resume_id
	})

	_, err := handler(context.Background(), request)
	if err == nil {
		t.Errorf("Expected error for missing resume_id")
	}
}

func TestListTemplatesTool_NonExistentResume(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewListTemplatesTool(db)

	request := createTestRequest(map[string]interface{}{
		"resume_id": "999", // Non-existent resume
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	// Should still return success with empty list
	textContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	if !strings.Contains(textContent.Text, "Templates listed successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	if !strings.Contains(textContent.Text, "\"count\":0") {
		t.Errorf("Expected count of 0 templates, got: %s", textContent.Text)
	}
}