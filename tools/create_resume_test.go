package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestCreateResumeTool_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	tool, handler := NewCreateResumeTool(db)

	// Test tool creation
	if tool.Name != "create_resume" {
		t.Errorf("Expected tool name 'create_resume', got %s", tool.Name)
	}

	// Create request
	request := createTestRequest(map[string]interface{}{
		"name":        "John Doe",
		"photo":       "photo.jpg",
		"description": "Software Engineer",
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

	if !strings.Contains(textContent.Text, "Resume created successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	// Verify resume was created in database
	resumes, err := db.ListResumes()
	if err != nil {
		t.Fatalf("Failed to list resumes: %v", err)
	}

	if len(resumes) != 1 {
		t.Errorf("Expected 1 resume, got %d", len(resumes))
	}

	resume := resumes[0]
	if resume.Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got %s", resume.Name)
	}
}

func TestCreateResumeTool_CopyFromExisting(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create source resume with full data
	sourceResume := createFullTestResume(t, db)
	createTestTemplate(t, db, sourceResume.ID)

	_, handler := NewCreateResumeTool(db)

	// Create new resume by copying from existing one
	request := createTestRequest(map[string]interface{}{
		"name":                "Jane Doe",
		"photo":               "jane.jpg",
		"description":         "Senior Software Engineer",
		"copy_from_resume_id": "1",
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

	if !strings.Contains(textContent.Text, "Resume created successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	if !strings.Contains(textContent.Text, "copied_from_resume_id") {
		t.Errorf("Expected copied_from_resume_id in response, got: %s", textContent.Text)
	}

	// Verify new resume was created
	resumes, err := db.ListResumes()
	if err != nil {
		t.Fatalf("Failed to list resumes: %v", err)
	}

	if len(resumes) != 2 {
		t.Errorf("Expected 2 resumes, got %d", len(resumes))
	}

	// Find the new resume
	var newResumeID uint
	for _, resume := range resumes {
		if resume.Name == "Jane Doe" {
			newResumeID = resume.ID
			break
		}
	}

	if newResumeID == 0 {
		t.Fatal("Could not find new resume")
	}

	// Get full resume data to verify copying
	fullNewResume, err := db.GetResumeByID(newResumeID)
	if err != nil {
		t.Fatalf("Failed to get full new resume: %v", err)
	}

	// Verify basic info is new
	if fullNewResume.Name != "Jane Doe" {
		t.Errorf("Expected name 'Jane Doe', got %s", fullNewResume.Name)
	}
	if fullNewResume.Photo != "jane.jpg" {
		t.Errorf("Expected photo 'jane.jpg', got %s", fullNewResume.Photo)
	}
	if fullNewResume.Description != "Senior Software Engineer" {
		t.Errorf("Expected description 'Senior Software Engineer', got %s", fullNewResume.Description)
	}

	// Verify contacts were copied
	if len(fullNewResume.Contacts) != 2 {
		t.Errorf("Expected 2 contacts, got %d", len(fullNewResume.Contacts))
	}

	// Verify work experiences were copied
	if len(fullNewResume.WorkExperiences) != 1 {
		t.Errorf("Expected 1 work experience, got %d", len(fullNewResume.WorkExperiences))
	}

	// Verify education was copied
	if len(fullNewResume.Educations) != 1 {
		t.Errorf("Expected 1 education, got %d", len(fullNewResume.Educations))
	}

	// Verify other experiences were copied
	if len(fullNewResume.OtherExperiences) != 1 {
		t.Errorf("Expected 1 other experience, got %d", len(fullNewResume.OtherExperiences))
	}

	// Verify templates were copied
	templates, err := db.ListTemplatesByResumeID(newResumeID)
	if err != nil {
		t.Fatalf("Failed to list templates: %v", err)
	}

	if len(templates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(templates))
	}

	// Verify feature maps were copied for work experience
	if len(fullNewResume.WorkExperiences) > 0 && len(fullNewResume.WorkExperiences[0].FeatureMaps) < 1 {
		t.Errorf("Expected at least 1 work experience feature map, got %d", len(fullNewResume.WorkExperiences[0].FeatureMaps))
	}
	
	// Verify feature maps were copied for education
	if len(fullNewResume.Educations) > 0 && len(fullNewResume.Educations[0].FeatureMaps) < 1 {
		t.Errorf("Expected at least 1 education feature map, got %d", len(fullNewResume.Educations[0].FeatureMaps))
	}
	
	// Verify feature maps were copied for other experiences
	if len(fullNewResume.OtherExperiences) > 0 && len(fullNewResume.OtherExperiences[0].FeatureMaps) < 1 {
		t.Errorf("Expected at least 1 other experience feature map, got %d", len(fullNewResume.OtherExperiences[0].FeatureMaps))
	}
}

func TestCreateResumeTool_CopyFromNonExistent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewCreateResumeTool(db)

	request := createTestRequest(map[string]interface{}{
		"name":                "Jane Doe",
		"description":         "Senior Software Engineer",
		"copy_from_resume_id": "999", // Non-existent resume
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

func TestCreateResumeTool_CopyFromInvalidID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewCreateResumeTool(db)

	request := createTestRequest(map[string]interface{}{
		"name":                "Jane Doe",
		"description":         "Senior Software Engineer",
		"copy_from_resume_id": "invalid", // Invalid ID format
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

func TestCreateResumeTool_CopyEmptyResume(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	// Create empty source resume
	_ = createTestResume(t, db)

	_, handler := NewCreateResumeTool(db)

	request := createTestRequest(map[string]interface{}{
		"name":                "Jane Doe",
		"description":         "Senior Software Engineer",
		"copy_from_resume_id": "1",
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

	if !strings.Contains(textContent.Text, "Resume created successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	// Verify 2 resumes exist
	resumes, err := db.ListResumes()
	if err != nil {
		t.Fatalf("Failed to list resumes: %v", err)
	}

	if len(resumes) != 2 {
		t.Errorf("Expected 2 resumes, got %d", len(resumes))
	}
}

func TestCreateResumeTool_MissingRequiredFields(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewCreateResumeTool(db)

	tests := []struct {
		name string
		args map[string]interface{}
	}{
		{
			name: "missing name",
			args: map[string]interface{}{
				"description": "Software Engineer",
			},
		},
		{
			name: "missing description",
			args: map[string]interface{}{
				"name": "John Doe",
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

func TestCreateResumeTool_WithoutPhoto(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewCreateResumeTool(db)

	request := createTestRequest(map[string]interface{}{
		"name":        "John Doe",
		"description": "Software Engineer",
		// No photo provided
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

	if !strings.Contains(textContent.Text, "Resume created successfully") {
		t.Errorf("Expected success message, got: %s", textContent.Text)
	}

	// Verify resume was created with empty photo
	resumes, err := db.ListResumes()
	if err != nil {
		t.Fatalf("Failed to list resumes: %v", err)
	}

	if len(resumes) != 1 {
		t.Errorf("Expected 1 resume, got %d", len(resumes))
	}

	resume := resumes[0]
	if resume.Photo != "" {
		t.Errorf("Expected empty photo, got %s", resume.Photo)
	}
}