package tools

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/models"
)

func createFullTestResume(t *testing.T, db *database.Database) *models.Resume {
	resume := createTestResume(t, db)
	
	// Add contact info
	contact1 := &models.Contact{
		ResumeID: resume.ID,
		Key:      "email",
		Value:    "test@example.com",
	}
	contact2 := &models.Contact{
		ResumeID: resume.ID,
		Key:      "phone",
		Value:    "+1234567890",
	}
	db.AddContact(contact1)
	db.AddContact(contact2)
	
	// Add work experience
	workExp := &models.WorkExperience{
		ResumeID:  resume.ID,
		Company:   "Tech Corp",
		JobTitle:  "Software Engineer",
		StartDate: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   &time.Time{},
	}
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	workExp.EndDate = &endDate
	db.AddWorkExperience(workExp)
	
	// Add work experience feature map
	workFeature := &models.FeatureMap{
		ExperienceID: workExp.ID,
		Key:          "skills",
		Value:        "Go, Python, React",
	}
	db.AddFeatureMap(workFeature)
	
	// Add education
	education := &models.Education{
		ResumeID:   resume.ID,
		SchoolName: "University of Technology",
		StartDate:  time.Date(2016, 9, 1, 0, 0, 0, 0, time.UTC),
		EndDate:    &time.Time{},
	}
	eduEndDate := time.Date(2020, 6, 30, 0, 0, 0, 0, time.UTC)
	education.EndDate = &eduEndDate
	db.AddEducation(education)
	
	// Add education feature map
	eduFeature := &models.FeatureMap{
		ExperienceID: education.ID,
		Key:          "degree",
		Value:        "Bachelor of Science in Computer Science",
	}
	db.AddFeatureMap(eduFeature)
	
	// Add other experience
	otherExp := &models.OtherExperience{
		ResumeID: resume.ID,
		Category: "Projects",
	}
	db.AddOtherExperience(otherExp)
	
	// Add other experience feature map
	otherFeature := &models.FeatureMap{
		ExperienceID: otherExp.ID,
		Key:          "project_name",
		Value:        "E-commerce Platform",
	}
	db.AddFeatureMap(otherFeature)
	
	return resume
}

func TestGetResumeContextTool_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	resume := createFullTestResume(t, db)
	
	// Create a template for this resume
	_ = createTestTemplate(t, db, resume.ID)

	tool, handler := NewGetResumeContextTool(db)

	// Test tool creation
	if tool.Name != "get_resume_context" {
		t.Errorf("Expected tool name 'get_resume_context', got %s", tool.Name)
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
	if len(result.Content) < 2 {
		t.Fatal("Handler returned insufficient content")
	}

	messageContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent for message, got %T", result.Content[0])
	}

	if !strings.Contains(messageContent.Text, "Resume JSON schema retrieved successfully") {
		t.Errorf("Expected success message, got: %s", messageContent.Text)
	}

	// Check the JSON data content
	jsonContent, ok := result.Content[1].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent for JSON data, got %T", result.Content[1])
	}

	// Check that JSON schema is included
	if !strings.Contains(jsonContent.Text, "json_schema") {
		t.Errorf("Expected json_schema in response, got: %s", jsonContent.Text)
	}
	
	// Check that schema contains expected structure
	if !strings.Contains(jsonContent.Text, "properties") {
		t.Errorf("Expected JSON schema properties in response, got: %s", jsonContent.Text)
	}
}

func TestGetResumeContextTool_SchemaValidation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	_ = createFullTestResume(t, db)

	_, handler := NewGetResumeContextTool(db)

	// Create request
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

	if len(result.Content) < 2 {
		t.Fatal("Handler returned insufficient content")
	}

	jsonContent, ok := result.Content[1].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent for JSON data, got %T", result.Content[1])
	}

	// Should contain JSON schema
	if !strings.Contains(jsonContent.Text, "json_schema") {
		t.Errorf("Expected json_schema in response, got: %s", jsonContent.Text)
	}
	
	// Should contain JSON schema structure
	if !strings.Contains(jsonContent.Text, "properties") {
		t.Errorf("Expected JSON schema properties in response, got: %s", jsonContent.Text)
	}
}

func TestGetResumeContextTool_EmptyResume(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	// Create minimal resume with no additional data
	_ = createTestResume(t, db)

	_, handler := NewGetResumeContextTool(db)

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

	if len(result.Content) < 2 {
		t.Fatal("Handler returned insufficient content")
	}

	messageContent, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent for message, got %T", result.Content[0])
	}

	if !strings.Contains(messageContent.Text, "Resume JSON schema retrieved successfully") {
		t.Errorf("Expected success message, got: %s", messageContent.Text)
	}

	jsonContent, ok := result.Content[1].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent for JSON data, got %T", result.Content[1])
	}

	// Should contain JSON schema even for empty resume
	if !strings.Contains(jsonContent.Text, "json_schema") {
		t.Errorf("Expected json_schema in response, got: %s", jsonContent.Text)
	}
}

func TestGetResumeContextTool_InvalidResumeID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewGetResumeContextTool(db)

	request := createTestRequest(map[string]interface{}{
		"resume_id": "999", // Non-existent resume
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

func TestGetResumeContextTool_InvalidResumeIDFormat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewGetResumeContextTool(db)

	request := createTestRequest(map[string]interface{}{
		"resume_id": "invalid", // Invalid format
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

func TestGetResumeContextTool_MissingResumeID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	_, handler := NewGetResumeContextTool(db)

	request := createTestRequest(map[string]interface{}{
		// Missing resume_id
	})

	_, err := handler(context.Background(), request)
	if err == nil {
		t.Errorf("Expected error for missing resume_id")
	}
}

func TestGetResumeContextTool_SchemaStructure(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	_ = createFullTestResume(t, db)

	_, handler := NewGetResumeContextTool(db)

	request := createTestRequest(map[string]interface{}{
		"resume_id": "1",
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if len(result.Content) < 2 {
		t.Fatal("Handler returned insufficient content")
	}

	jsonContent, ok := result.Content[1].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent for JSON data, got %T", result.Content[1])
	}

	// Check that schema contains expected JSON schema structure
	if !strings.Contains(jsonContent.Text, "properties") {
		t.Errorf("Expected JSON schema properties in response, got: %s", jsonContent.Text)
	}
}

func TestGetResumeContextTool_JSONSchemaContent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	
	_ = createFullTestResume(t, db)

	_, handler := NewGetResumeContextTool(db)

	request := createTestRequest(map[string]interface{}{
		"resume_id": "1",
	})

	result, err := handler(context.Background(), request)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if len(result.Content) < 2 {
		t.Fatal("Handler returned insufficient content")
	}

	jsonContent, ok := result.Content[1].(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent for JSON data, got %T", result.Content[1])
	}

	// Check that JSON schema contains expected structure
	expectedSchemaElements := []string{
		"type",
		"properties",
		"required",
	}

	for _, element := range expectedSchemaElements {
		if !strings.Contains(jsonContent.Text, element) {
			t.Errorf("Expected JSON schema to contain %s, got: %s", element, jsonContent.Text)
		}
	}

	// Check that schema contains Resume model fields
	if !strings.Contains(jsonContent.Text, "\"name\"") {
		t.Errorf("Expected JSON schema to contain name field, got: %s", jsonContent.Text)
	}
}