# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Resume MCP (Model Context Protocol) server built in Go that allows AI agents to manage resume data and generate PDF previews. The project combines an MCP server with a REST API for preview functionality.

## Core Architecture

- **MCP Server**: Built using `github.com/mark3labs/mcp-go` for creating MCP tools
- **REST API**: Fiber framework for HTTP endpoints serving PDF previews
- **Database**: GORM with SQLite for local data storage
- **Template Engine**: Go templates for resume rendering
- **PDF Generation**: On-demand PDF generation served through preview URLs

## Key Components

### Data Model
- Resume basic info (name, photo, description)
- Contact information (key-value pairs)
- Work experience (dates, company, job title)
- Education experience (dates, school name)
- Feature maps for experiences (flexible JSON data like GPA, salary, features array)
- Other experiences with categories

### MCP Tools
The server provides these MCP tools:
- `create_resume` - Create new resume with name, photo, and description
- `update_basic_info` - Update resume basic information
- `add_contact_info` - Add contact information
- `add_work_experience` - Add work experience entries
- `add_education` - Add education entries
- `add_other_experience` - Add other experience categories
- `add_feature_map` - Add flexible JSON features to experiences
- `update_feature_map` - Update feature maps
- `delete_feature_map` - Delete feature maps
- `get_resume_by_name` - Get structured resume data
- `list_resumes` - List all saved resumes
- `delete_resume` - Delete resume by ID
- `generate_preview` - Generate HTML preview with template and CSS
- `update_preview_style` - Update CSS styles for existing preview

### Workflow
1. AI agent calls `generate_resume` with Go template
2. Server stores template in DB and returns preview URL (`https://localhost:8080/resume/preview/:sid`)
3. User visits URL to view rendered PDF

## Development Commands

Based on the existing CLAUDE.md and project structure, all commands should be defined in a Makefile:

```bash
# Build the project
make build

# Run tests
make test

# Run the MCP server
make run

# Clean build artifacts
make clean
```

## Testing Strategy

Write comprehensive tests for:
- MCP tool implementations
- Database operations (GORM models)
- Template rendering
- PDF generation
- REST API endpoints

## Key Dependencies

- `github.com/mark3labs/mcp-go` - MCP server framework
- Fiber - HTTP framework for REST API
- GORM - ORM for database operations
- SQLite - Local database storage
- Go template - Template rendering engine

## File Structure

- `docs/DesignDocument.md` - Detailed technical specifications
- `go.mod` - Go module definition
- Source code will be organized in standard Go project structure

## Implementation Status

- ✅ **Complete Implementation**: All core components are implemented and working
- ✅ **MCP Server**: Fully functional with resume management tools
- ✅ **Database Layer**: GORM with SQLite, automatic migrations
- ✅ **REST API**: HTML preview with template rendering and CSS styling
- ✅ **Template System**: Go templates with CSS support for HTML generation

## Important Notes

- Use the stop hook: `make build && make test` to validate changes
- HTML generation happens in-memory and is served through preview URLs
- The server runs both MCP (stdio) and HTTP (port 8080) simultaneously
- Database file: `resume.db` (SQLite) created automatically
- All MCP tools are implemented and ready for AI agent interaction