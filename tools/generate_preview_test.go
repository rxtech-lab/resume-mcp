package tools

import (
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rxtech-lab/resume-mcp/internal/service"
)

func TestGeneratePreviewTool_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	templateService := service.NewTemplateService()
	resume := createTestResume(t, db)
	template := createTestTemplate(t, db, resume.ID)
	port := "8080"

	tool, handler := NewGeneratePreviewTool(db, port, templateService)

	// Test tool creation
	if tool.Name != "generate_preview" {
		t.Errorf("Expected tool name 'generate_preview', got %s", tool.Name)
	}

	// Create request
	request := createTestRequest(map[string]interface{}{
		"resume_id":   "1",
		"template_id": "1",
	})

	// Execute handler
	result, err := handler(createTestContext(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Handler returned nil result")
	}

	// Check result content - should have 3 items: message, preview URL, download URL
	if len(result.Content) < 3 {
		t.Fatalf("Handler returned insufficient content, expected 3, got %d", len(result.Content))
	}

	firstContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected first TextContent, got %T", result.Content[0])
	}

	if !strings.Contains(firstContent.Text, "Preview generated successfully") {
		t.Errorf("Expected success message, got: %s", firstContent.Text)
	}

	// Second content should be the preview URL
	secondContent, ok := result.Content[1].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected second TextContent, got %T", result.Content[1])
	}

	expectedPreviewPrefix := "Preview: http://localhost:8080/resume/preview/"
	if !strings.Contains(secondContent.Text, expectedPreviewPrefix) {
		t.Errorf("Expected preview URL to contain %s, got: %s", expectedPreviewPrefix, secondContent.Text)
	}

	// Third content should be the download URL
	thirdContent, ok := result.Content[2].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected third TextContent, got %T", result.Content[2])
	}

	expectedDownloadPrefix := "Download PDF: http://localhost:8080/resume/download/"
	if !strings.Contains(thirdContent.Text, expectedDownloadPrefix) {
		t.Errorf("Expected download URL to contain %s, got: %s", expectedDownloadPrefix, thirdContent.Text)
	}

	// Verify that a preview session was created - extract session ID from preview URL
	sessionID := strings.TrimPrefix(strings.TrimSpace(secondContent.Text), expectedPreviewPrefix)
	if sessionID == "" || sessionID == secondContent.Text {
		t.Errorf("Could not extract session ID from URL: %s", secondContent.Text)
	}

	session, err := db.GetPreviewSession(sessionID, nil)
	if err != nil {
		t.Fatalf("Failed to get preview session: %v", err)
	}

	if session.ResumeID != resume.ID {
		t.Errorf("Expected session resume ID %d, got %d", resume.ID, session.ResumeID)
	}

	if session.Template != template.TemplateData {
		t.Errorf("Expected session template to match template data")
	}
}

func TestGeneratePreviewTool_WithCSS(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	templateService := service.NewTemplateService()
	resume := createTestResume(t, db)
	_ = createTestTemplate(t, db, resume.ID)
	port := "8080"

	_, handler := NewGeneratePreviewTool(db, port, templateService)

	// Create request with CSS
	request := createTestRequest(map[string]interface{}{
		"resume_id":   "1",
		"template_id": "1",
		"css":         "body { background-color: #f0f0f0; }",
	})

	result, err := handler(createTestContext(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if result == nil {
		t.Fatal("Handler returned nil result")
	}

	// Extract session ID from preview URL (second content item)
	secondContent, ok := result.Content[1].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected second TextContent, got %T", result.Content[1])
	}

	expectedURL := "Preview: http://localhost:8080/resume/preview/"
	sessionID := strings.TrimPrefix(strings.TrimSpace(secondContent.Text), expectedURL)

	session, err := db.GetPreviewSession(sessionID, nil)
	if err != nil {
		t.Fatalf("Failed to get preview session: %v", err)
	}

	if session.CSS != "body { background-color: #f0f0f0; }" {
		t.Errorf("Expected CSS to be stored in session, got: %s", session.CSS)
	}
}

func TestGeneratePreviewTool_InvalidResumeID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	templateService := service.NewTemplateService()
	port := "8080"

	_, handler := NewGeneratePreviewTool(db, port, templateService)

	request := createTestRequest(map[string]interface{}{
		"resume_id":   "999", // Non-existent resume
		"template_id": "1",
	})

	result, err := handler(createTestContext(), request)
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

	if !strings.Contains(textContent.Text, "Error getting resume") {
		t.Errorf("Expected 'Error getting resume' error, got: %s", textContent.Text)
	}
}

func TestGeneratePreviewTool_InvalidTemplateID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	templateService := service.NewTemplateService()
	_ = createTestResume(t, db)
	port := "8080"

	_, handler := NewGeneratePreviewTool(db, port, templateService)

	request := createTestRequest(map[string]interface{}{
		"resume_id":   "1",
		"template_id": "999", // Non-existent template
	})

	result, err := handler(createTestContext(), request)
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

	if !strings.Contains(textContent.Text, "Error getting template") {
		t.Errorf("Expected 'Error getting template' error, got: %s", textContent.Text)
	}
}

func TestGeneratePreviewTool_TemplateMismatch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	templateService := service.NewTemplateService()
	_ = createTestResume(t, db)
	resume2 := createTestResume(t, db)
	resume2.Name = "Second Resume"
	db.UpdateResume(resume2, nil)

	// Create template for resume2
	_ = createTestTemplate(t, db, resume2.ID)
	port := "8080"

	_, handler := NewGeneratePreviewTool(db, port, templateService)

	// Try to use resume1 with template that belongs to resume2
	request := createTestRequest(map[string]interface{}{
		"resume_id":   "1", // resume1
		"template_id": "1", // template belongs to resume2
	})

	result, err := handler(createTestContext(), request)
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

	if !strings.Contains(textContent.Text, "Template does not belong to the specified resume") {
		t.Errorf("Expected template mismatch error, got: %s", textContent.Text)
	}
}

func TestGeneratePreviewTool_MissingRequiredFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	templateService := service.NewTemplateService()
	port := "8080"

	_, handler := NewGeneratePreviewTool(db, port, templateService)

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "missing resume_id",
			args: map[string]interface{}{
				"template_id": "1",
			},
		},
		{
			name: "missing template_id",
			args: map[string]interface{}{
				"resume_id": "1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := createTestRequest(tt.args)

			_, err := handler(createTestContext(), request)
			if err == nil {
				t.Errorf("Expected error for missing required field")
			}
		})
	}
}

func TestGeneratePreviewTool_InvalidIDs(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	templateService := service.NewTemplateService()
	port := "8080"

	_, handler := NewGeneratePreviewTool(db, port, templateService)

	tests := []struct {
		name          string
		resumeID      string
		templateID    string
		expectedError string
	}{
		{
			name:          "invalid resume_id",
			resumeID:      "invalid",
			templateID:    "1",
			expectedError: "Invalid resume_id",
		},
		{
			name:          "invalid template_id",
			resumeID:      "1",
			templateID:    "invalid",
			expectedError: "Invalid template_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := createTestRequest(map[string]interface{}{
				"resume_id":   tt.resumeID,
				"template_id": tt.templateID,
			})

			result, err := handler(createTestContext(), request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			textContent, ok := result.Content[0].(mcp.TextContent)
			if !ok {
				t.Fatalf("Expected TextContent, got %T", result.Content[0])
			}

			if !strings.Contains(textContent.Text, tt.expectedError) {
				t.Errorf("Expected '%s' error, got: %s", tt.expectedError, textContent.Text)
			}
		})
	}
}
