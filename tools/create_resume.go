package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/models"
)

func NewCreateResumeTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("create_resume",
		mcp.WithDescription("Create a new resume with basic information including name, photo, and description. Optionally copy all data from an existing resume. Returns the created resume ID for use with other tools."),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the resume owner"),
		),
		mcp.WithString("photo",
			mcp.Description("URL or path to the photo"),
		),
		mcp.WithString("description",
			mcp.Required(),
			mcp.Description("Brief description or summary"),
		),
		mcp.WithString("copy_from_resume_id",
			mcp.Description("Optional: ID of an existing resume to copy all data from (contacts, work experiences, education, etc.). The new resume will have the provided name, photo, and description, but all other data will be copied from the source resume. If user ask to create a new resume base on the existing one, please use this parameter and don't need to add all the data manually."),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return nil, fmt.Errorf("name parameter is required: %w", err)
		}

		description, err := request.RequireString("description")
		if err != nil {
			return nil, fmt.Errorf("description parameter is required: %w", err)
		}

		photo := request.GetString("photo", "")
		copyFromResumeIDStr := request.GetString("copy_from_resume_id", "")

		// Create the new resume with basic information
		resume := &models.Resume{
			Name:        name,
			Photo:       photo,
			Description: description,
		}

		if err := db.CreateResume(resume); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error creating resume: %v", err)), nil
		}

		// If copy_from_resume_id is provided, copy all data from the source resume
		if copyFromResumeIDStr != "" {
			copyFromResumeID, err := strconv.Atoi(copyFromResumeIDStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid copy_from_resume_id: %v", err)), nil
			}

			// Get the source resume with all related data
			sourceResume, err := db.GetResumeByID(uint(copyFromResumeID))
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Source resume not found: %v", err)), nil
			}

			// Copy contacts
			for _, contact := range sourceResume.Contacts {
				newContact := &models.Contact{
					ResumeID: resume.ID,
					Key:      contact.Key,
					Value:    contact.Value,
				}
				if err := db.AddContact(newContact); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Error copying contact: %v", err)), nil
				}
			}

			// Copy work experiences
			for _, workExp := range sourceResume.WorkExperiences {
				newWorkExp := &models.WorkExperience{
					ResumeID:  resume.ID,
					Company:   workExp.Company,
					JobTitle:  workExp.JobTitle,
					Type:      workExp.Type,
					StartDate: workExp.StartDate,
					EndDate:   workExp.EndDate,
				}
				if err := db.AddWorkExperience(newWorkExp); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Error copying work experience: %v", err)), nil
				}

				// Copy feature maps for this work experience
				for _, featureMap := range workExp.FeatureMaps {
					newFeatureMap := &models.FeatureMap{
						ExperienceID: newWorkExp.ID,
						Key:          featureMap.Key,
						Value:        featureMap.Value,
					}
					if err := db.AddFeatureMap(newFeatureMap); err != nil {
						return mcp.NewToolResultError(fmt.Sprintf("Error copying work experience feature map: %v", err)), nil
					}
				}
			}

			// Copy education
			for _, education := range sourceResume.Educations {
				newEducation := &models.Education{
					ResumeID:   resume.ID,
					SchoolName: education.SchoolName,
					Type:       education.Type,
					StartDate:  education.StartDate,
					EndDate:    education.EndDate,
				}
				if err := db.AddEducation(newEducation); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Error copying education: %v", err)), nil
				}

				// Copy feature maps for this education
				for _, featureMap := range education.FeatureMaps {
					newFeatureMap := &models.FeatureMap{
						ExperienceID: newEducation.ID,
						Key:          featureMap.Key,
						Value:        featureMap.Value,
					}
					if err := db.AddFeatureMap(newFeatureMap); err != nil {
						return mcp.NewToolResultError(fmt.Sprintf("Error copying education feature map: %v", err)), nil
					}
				}
			}

			// Copy other experiences
			for _, otherExp := range sourceResume.OtherExperiences {
				newOtherExp := &models.OtherExperience{
					ResumeID: resume.ID,
					Category: otherExp.Category,
				}
				if err := db.AddOtherExperience(newOtherExp); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Error copying other experience: %v", err)), nil
				}

				// Copy feature maps for this other experience
				for _, featureMap := range otherExp.FeatureMaps {
					newFeatureMap := &models.FeatureMap{
						ExperienceID: newOtherExp.ID,
						Key:          featureMap.Key,
						Value:        featureMap.Value,
					}
					if err := db.AddFeatureMap(newFeatureMap); err != nil {
						return mcp.NewToolResultError(fmt.Sprintf("Error copying other experience feature map: %v", err)), nil
					}
				}
			}

			// Copy templates
			sourceTemplates, err := db.ListTemplatesByResumeID(uint(copyFromResumeID))
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Error getting source templates: %v", err)), nil
			}

			for _, template := range sourceTemplates {
				newTemplate := &models.Template{
					ResumeID:     resume.ID,
					Name:         template.Name,
					Description:  template.Description,
					TemplateData: template.TemplateData,
				}
				if err := db.CreateTemplate(newTemplate); err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Error copying template: %v", err)), nil
				}
			}
		}

		result := map[string]interface{}{
			"id":          resume.ID,
			"name":        resume.Name,
			"photo":       resume.Photo,
			"description": resume.Description,
			"created_at":  resume.CreatedAt,
		}

		if copyFromResumeIDStr != "" {
			result["copied_from_resume_id"] = copyFromResumeIDStr
			result["message"] = fmt.Sprintf("Resume created successfully and copied data from resume ID %s", copyFromResumeIDStr)
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Resume created successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}
