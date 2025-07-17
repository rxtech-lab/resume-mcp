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

func NewAddContactInfoTool(db *database.Database) (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("add_contact_info",
		mcp.WithDescription("Add contact information to a resume as key-value pairs (e.g., email, phone, linkedin, github, etc.)."),
		mcp.WithString("resume_id",
			mcp.Required(),
			mcp.Description("The ID of the resume to add contact info to"),
		),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("The contact type (email, phone, linkedin, etc.)"),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("The contact value"),
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

		key, err := request.RequireString("key")
		if err != nil {
			return nil, fmt.Errorf("key parameter is required: %w", err)
		}

		value, err := request.RequireString("value")
		if err != nil {
			return nil, fmt.Errorf("value parameter is required: %w", err)
		}

		contact := &models.Contact{
			ResumeID: uint(resumeID),
			Key:      key,
			Value:    value,
		}

		if err := db.AddContact(contact); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error adding contact info: %v", err)), nil
		}

		result := map[string]interface{}{
			"id":        contact.ID,
			"resume_id": contact.ResumeID,
			"key":       contact.Key,
			"value":     contact.Value,
		}

		resultJSON, _ := json.Marshal(result)
		return mcp.NewToolResultText(fmt.Sprintf("Contact info added successfully: %s", string(resultJSON))), nil
	}

	return tool, handler
}