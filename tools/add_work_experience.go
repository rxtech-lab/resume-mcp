package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/models"
)

func NewAddWorkExperienceTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("add_work_experience",
		mcp.WithDescription("Add work experience to a resume with company, job title, and date range. Use feature maps to add additional details like responsibilities or achievements."),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("The ID of the resume to add work experience to"),
		),
		mcp.WithString("company",
			mcp.Required(),
			mcp.Description("The company name"),
		),
		mcp.WithString("job_title",
			mcp.Required(),
			mcp.Description("The job title"),
		),
		mcp.WithString("type",
			mcp.Description("Type of work experience: fulltime, parttime, or internship (default: fulltime)"),
			mcp.WithStringEnumItems(
				[]string{"fulltime", "parttime", "internship"},
			),
		),
		mcp.WithString("start_date",
			mcp.Required(),
			mcp.Description("Start date in YYYY-MM-DD format"),
		),
		mcp.WithString("end_date",
			mcp.Description("End date in YYYY-MM-DD format (optional for current job)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		resumeIDStr, err := request.RequireString("resume_id")
		if err != nil {
			return nil, fmt.Errorf("resume_id parameter is required: %w", err)
		}

		resumeID, err := strconv.ParseUint(resumeIDStr, 10, 32)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid resume_id: %v", err)), nil
		}

		company, err := request.RequireString("company")
		if err != nil {
			return nil, fmt.Errorf("company parameter is required: %w", err)
		}

		jobTitle, err := request.RequireString("job_title")
		if err != nil {
			return nil, fmt.Errorf("job_title parameter is required: %w", err)
		}

		workType := request.GetString("type", "fulltime")
		// Validate type
		if workType != "fulltime" && workType != "parttime" && workType != "internship" {
			return mcp.NewToolResultError("Invalid type. Must be: fulltime, parttime, or internship"), nil
		}

		startDateStr, err := request.RequireString("start_date")
		if err != nil {
			return nil, fmt.Errorf("start_date parameter is required: %w", err)
		}

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid start_date format: %v", err)), nil
		}

		endDateStr := request.GetString("end_date", "")
		var endDate *time.Time
		if endDateStr != "" {
			parsedEndDate, err := time.Parse("2006-01-02", endDateStr)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid end_date format: %v", err)), nil
			}
			endDate = &parsedEndDate
		}

		workExp := &models.WorkExperience{
			ResumeID:  uint(resumeID),
			Company:   company,
			JobTitle:  jobTitle,
			Type:      workType,
			StartDate: startDate,
			EndDate:   endDate,
		}

		if err := db.AddWorkExperience(workExp); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error adding work experience: %v", err)), nil
		}

		result := map[string]interface{}{
			"id":         workExp.ID,
			"resume_id":  workExp.ResumeID,
			"company":    workExp.Company,
			"job_title":  workExp.JobTitle,
			"type":       workExp.Type,
			"start_date": workExp.StartDate.Format("2006-01-02"),
		}

		if workExp.EndDate != nil {
			result["end_date"] = workExp.EndDate.Format("2006-01-02")
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Work experience added successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}
