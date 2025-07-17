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

func NewAddEducationTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("add_education",
		mcp.WithDescription("Add education experience to a resume with school name and date range. Use feature maps to add details like degree, GPA, or coursework."),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("The ID of the resume to add education to"),
		),
		mcp.WithString("school_name",
			mcp.Required(),
			mcp.Description("The name of the school"),
		),
		mcp.WithString("type",
			mcp.Description("Type of education: fulltime, parttime, or internship (default: fulltime)"),
			mcp.WithStringEnumItems(
				[]string{"fulltime", "parttime", "internship"},
			),
		),
		mcp.WithString("start_date",
			mcp.Required(),
			mcp.Description("Start date in YYYY-MM-DD format"),
		),
		mcp.WithString("end_date",
			mcp.Description("End date in YYYY-MM-DD format (optional for current education)"),
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

		schoolName, err := request.RequireString("school_name")
		if err != nil {
			return nil, fmt.Errorf("school_name parameter is required: %w", err)
		}

		eduType := request.GetString("type", "fulltime")
		// Validate type
		if eduType != "fulltime" && eduType != "parttime" && eduType != "internship" {
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

		education := &models.Education{
			ResumeID:   uint(resumeID),
			SchoolName: schoolName,
			Type:       eduType,
			StartDate:  startDate,
			EndDate:    endDate,
		}

		if err := db.AddEducation(education); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error adding education: %v", err)), nil
		}

		result := map[string]interface{}{
			"id":          education.ID,
			"resume_id":   education.ResumeID,
			"school_name": education.SchoolName,
			"type":        education.Type,
			"start_date":  education.StartDate.Format("2006-01-02"),
		}

		if education.EndDate != nil {
			result["end_date"] = education.EndDate.Format("2006-01-02")
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Education added successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}
