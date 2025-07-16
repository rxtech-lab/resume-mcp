package mcp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/models"
)

func newTextResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.NewTextContent(text),
		},
	}
}

func newErrorResult(message string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []interface{}{
			mcp.NewTextContent(message),
		},
		IsError: true,
	}
}

type ResumeTools struct {
	db *database.Database
}

func NewResumeTools(db *database.Database) *ResumeTools {
	return &ResumeTools{db: db}
}

func (r *ResumeTools) CreateResumeWithUser(name, photo, description string) (*mcp.CallToolResult, error) {
	resume := &models.Resume{
		Name:        name,
		Photo:       photo,
		Description: description,
	}

	if err := r.db.CreateResume(resume); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to create resume: %v", err)), nil
	}

	return newTextResult(fmt.Sprintf("Resume created successfully with ID: %d", resume.ID)), nil
}

func (r *ResumeTools) UpdateBasicInfo(resumeID uint, name, photo, description string) (*mcp.CallToolResult, error) {
	resume, err := r.db.GetResumeByID(resumeID)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Resume not found: %v", err)), nil
	}

	if name != "" {
		resume.Name = name
	}
	if photo != "" {
		resume.Photo = photo
	}
	if description != "" {
		resume.Description = description
	}

	if err := r.db.UpdateResume(resume); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to update resume: %v", err)), nil
	}

	return newTextResult("Resume updated successfully"), nil
}

func (r *ResumeTools) AddContactInfo(resumeID uint, key, value string) (*mcp.CallToolResult, error) {
	contact := &models.Contact{
		ResumeID: resumeID,
		Key:      key,
		Value:    value,
	}

	if err := r.db.AddContact(contact); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to add contact info: %v", err)), nil
	}

	return newTextResult(fmt.Sprintf("Contact info added successfully with ID: %d", contact.ID)), nil
}

func (r *ResumeTools) AddWorkExperience(resumeID uint, company, jobTitle string, startDate, endDate *time.Time) (*mcp.CallToolResult, error) {
	experience := &models.WorkExperience{
		ResumeID:  resumeID,
		Company:   company,
		JobTitle:  jobTitle,
		StartDate: *startDate,
		EndDate:   endDate,
	}

	if err := r.db.AddWorkExperience(experience); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to add work experience: %v", err)), nil
	}

	return newTextResult(fmt.Sprintf("Work experience added successfully with ID: %d", experience.ID)), nil
}

func (r *ResumeTools) AddEducation(resumeID uint, schoolName string, startDate, endDate *time.Time) (*mcp.CallToolResult, error) {
	education := &models.Education{
		ResumeID:   resumeID,
		SchoolName: schoolName,
		StartDate:  *startDate,
		EndDate:    endDate,
	}

	if err := r.db.AddEducation(education); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to add education: %v", err)), nil
	}

	return newTextResult(fmt.Sprintf("Education added successfully with ID: %d", education.ID)), nil
}

func (r *ResumeTools) AddOtherExperience(resumeID uint, category string) (*mcp.CallToolResult, error) {
	experience := &models.OtherExperience{
		ResumeID: resumeID,
		Category: category,
	}

	if err := r.db.AddOtherExperience(experience); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to add other experience: %v", err)), nil
	}

	return newTextResult(fmt.Sprintf("Other experience added successfully with ID: %d", experience.ID)), nil
}

func (r *ResumeTools) AddFeatureMap(experienceID uint, key string, value interface{}) (*mcp.CallToolResult, error) {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to marshal value: %v", err)), nil
	}

	featureMap := &models.FeatureMap{
		ExperienceID: experienceID,
		Key:          key,
		Value:        string(valueJSON),
	}

	if err := r.db.AddFeatureMap(featureMap); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to add feature map: %v", err)), nil
	}

	return newTextResult(fmt.Sprintf("Feature map added successfully with ID: %d", featureMap.ID)), nil
}

func (r *ResumeTools) UpdateFeatureMap(featureMapID uint, key string, value interface{}) (*mcp.CallToolResult, error) {
	valueJSON, err := json.Marshal(value)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to marshal value: %v", err)), nil
	}

	featureMap := &models.FeatureMap{
		ID:    featureMapID,
		Key:   key,
		Value: string(valueJSON),
	}

	if err := r.db.UpdateFeatureMap(featureMap); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to update feature map: %v", err)), nil
	}

	return newTextResult("Feature map updated successfully"), nil
}

func (r *ResumeTools) DeleteFeatureMap(featureMapID uint) (*mcp.CallToolResult, error) {
	if err := r.db.DeleteFeatureMap(featureMapID); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to delete feature map: %v", err)), nil
	}

	return newTextResult("Feature map deleted successfully"), nil
}

func (r *ResumeTools) GetResumeByName(name string) (*mcp.CallToolResult, error) {
	resume, err := r.db.GetResumeByName(name)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Resume not found: %v", err)), nil
	}

	result := map[string]interface{}{
		"resume": []map[string]interface{}{
			{
				"name": "basic-info",
				"data": map[string]interface{}{
					"name":        resume.Name,
					"photo":       resume.Photo,
					"description": resume.Description,
				},
			},
			{
				"name": "contact",
				"data": r.formatContacts(resume.Contacts),
			},
			{
				"name": "work-experience",
				"data": r.formatWorkExperiences(resume.WorkExperiences),
			},
			{
				"name": "education",
				"data": r.formatEducations(resume.Educations),
			},
			{
				"name": "other-experience",
				"data": r.formatOtherExperiences(resume.OtherExperiences),
			},
		},
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}

	return newTextResult(string(resultJSON)), nil
}

func (r *ResumeTools) ListResumes() (*mcp.CallToolResult, error) {
	resumes, err := r.db.ListResumes()
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to list resumes: %v", err)), nil
	}

	result := make([]map[string]interface{}, len(resumes))
	for i, resume := range resumes {
		result[i] = map[string]interface{}{
			"id":   resume.ID,
			"name": resume.Name,
		}
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return newErrorResult(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}

	return newTextResult(string(resultJSON)), nil
}

func (r *ResumeTools) DeleteResume(resumeID uint) (*mcp.CallToolResult, error) {
	if err := r.db.DeleteResume(resumeID); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to delete resume: %v", err)), nil
	}

	return newTextResult("Resume deleted successfully"), nil
}

func (r *ResumeTools) GeneratePreview(resumeID uint, template string, css string) (*mcp.CallToolResult, error) {
	sessionID := uuid.New().String()
	session := &models.PreviewSession{
		ID:       sessionID,
		ResumeID: resumeID,
		Template: template,
		CSS:      css,
	}

	if err := r.db.CreatePreviewSession(session); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to create preview session: %v", err)), nil
	}

	previewURL := fmt.Sprintf("http://localhost:8080/resume/preview/%s", sessionID)
	return newTextResult(previewURL), nil
}

func (r *ResumeTools) UpdatePreviewStyle(sessionID string, css string) (*mcp.CallToolResult, error) {
	if err := r.db.UpdatePreviewSessionCSS(sessionID, css); err != nil {
		return newErrorResult(fmt.Sprintf("Failed to update preview style: %v", err)), nil
	}

	return newTextResult("Preview style updated successfully"), nil
}

func (r *ResumeTools) formatContacts(contacts []models.Contact) map[string]interface{} {
	result := make(map[string]interface{})
	for _, contact := range contacts {
		result[contact.Key] = contact.Value
	}
	return result
}

func (r *ResumeTools) formatWorkExperiences(experiences []models.WorkExperience) []map[string]interface{} {
	result := make([]map[string]interface{}, len(experiences))
	for i, exp := range experiences {
		result[i] = map[string]interface{}{
			"id":         exp.ID,
			"company":    exp.Company,
			"job_title":  exp.JobTitle,
			"start_date": exp.StartDate,
			"end_date":   exp.EndDate,
		}
		
		if len(exp.FeatureMaps) > 0 {
			features := make(map[string]interface{})
			for _, fm := range exp.FeatureMaps {
				var value interface{}
				json.Unmarshal([]byte(fm.Value), &value)
				features[fm.Key] = value
			}
			result[i]["features"] = features
		}
	}
	return result
}

func (r *ResumeTools) formatEducations(educations []models.Education) []map[string]interface{} {
	result := make([]map[string]interface{}, len(educations))
	for i, edu := range educations {
		result[i] = map[string]interface{}{
			"id":          edu.ID,
			"school_name": edu.SchoolName,
			"start_date":  edu.StartDate,
			"end_date":    edu.EndDate,
		}
		
		if len(edu.FeatureMaps) > 0 {
			features := make(map[string]interface{})
			for _, fm := range edu.FeatureMaps {
				var value interface{}
				json.Unmarshal([]byte(fm.Value), &value)
				features[fm.Key] = value
			}
			result[i]["features"] = features
		}
	}
	return result
}

func (r *ResumeTools) formatOtherExperiences(experiences []models.OtherExperience) []map[string]interface{} {
	result := make([]map[string]interface{}, len(experiences))
	for i, exp := range experiences {
		result[i] = map[string]interface{}{
			"id":       exp.ID,
			"category": exp.Category,
		}
		
		if len(exp.FeatureMaps) > 0 {
			features := make(map[string]interface{})
			for _, fm := range exp.FeatureMaps {
				var value interface{}
				json.Unmarshal([]byte(fm.Value), &value)
				features[fm.Key] = value
			}
			result[i]["features"] = features
		}
	}
	return result
}