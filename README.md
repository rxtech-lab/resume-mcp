# Resume MCP Server

A Model Context Protocol (MCP) server built in Go that allows AI agents to manage resume data and generate HTML previews. The project combines an MCP server with a REST API for preview functionality.

## Features

- **Resume Management**: Create, update, and delete resumes with structured data
- **Template System**: Create and manage Go templates for resume rendering
- **HTML Preview**: Generate HTML previews with custom styling
- **Flexible Data Model**: Support for contacts, work experience, education, and custom feature maps
- **Copy Functionality**: Duplicate existing resumes and templates with all related data
- **REST API**: HTTP endpoints for preview generation and styling

## Installation

### Download Pre-built Binary

Download the latest release from: https://github.com/rxtech-lab/resume-mcp/releases

For macOS users, download `resume-mcp_macOS_arm64.pkg` and run the installer.

### Build from Source

```bash
# Clone the repository
git clone https://github.com/rxtech-lab/resume-mcp.git
cd resume-mcp

# Build the project
make build

# Run tests
make test

# Run the MCP server
make run
```

## Usage

### MCP Tools

The server provides comprehensive MCP tools for resume management:

#### Resume Management
- `create_resume` - Create new resume with basic info (supports copying from existing)
- `update_basic_info` - Update resume name, photo, and description
- `get_resume_by_name` - Retrieve resume data by name
- `list_resumes` - List all saved resumes
- `delete_resume` - Delete resume by ID

#### Contact Information
- `add_contact_info` - Add contact details (email, phone, etc.)

#### Experience Management
- `add_work_experience` - Add work experience entries
- `add_education` - Add education entries
- `add_other_experience` - Add other experience categories

#### Feature Maps
- `add_feature_map` - Add flexible JSON features to experiences
- `update_feature_map` - Update existing feature maps
- `delete_feature_map` - Delete feature maps

#### Template System
- `create_template` - Create Go templates for resume rendering (supports copying data)
- `get_template` - Retrieve template by ID
- `list_templates` - List templates for a resume
- `update_template` - Update existing templates
- `delete_template` - Delete templates

#### Preview and PDF Generation
- `generate_preview` - Generate HTML preview using template and resume data (returns preview and download URLs)
- `update_preview_style` - Update CSS styles for existing previews
- `get_resume_context` - Get comprehensive resume data and schema guide for template creation

When calling `generate_preview`, you'll receive both:
- A preview URL to view the resume in browser (includes a download button)
- A download URL to directly download the PDF version

### Copy Functionality

Both `create_resume` and `create_template` tools support copying from existing data:

```bash
# Copy all data from an existing resume when creating a new one
create_resume(name="New Resume", copy_from_resume_id="1")

# Copy data from one resume to another when creating a template
create_template(resume_id="2", copy_from_resume_id="1", name="My Template")
```

### Template Examples

The server supports Go template syntax for resume rendering:

```html
<div class="resume">
  <h1>{{.Name}}</h1>
  <p>{{.Description}}</p>
  
  {{if .Contacts}}
  <div class="contact">
    {{range .Contacts}}
    <p>{{.Key}}: {{.Value}}</p>
    {{end}}
  </div>
  {{end}}
  
  {{if .WorkExperiences}}
  <div class="experience">
    {{range .WorkExperiences}}
    <div>
      <h3>{{.JobTitle}} at {{.Company}}</h3>
      <p>{{.StartDate.Format "Jan 2006"}} - {{if .EndDate}}{{.EndDate.Format "Jan 2006"}}{{else}}Present{{end}}</p>
      {{range .FeatureMaps}}
      <p>{{.Key}}: {{.Value}}</p>
      {{end}}
    </div>
    {{end}}
  </div>
  {{end}}
</div>
```

## Architecture

### Core Components

- **MCP Server**: Built using `github.com/mark3labs/mcp-go`
- **REST API**: Fiber framework for HTTP endpoints
- **Database**: GORM with SQLite for local storage
- **Template Engine**: Go templates with Tailwind CSS support
- **Preview Generation**: On-demand HTML generation

### Data Model

- **Resume**: Basic info (name, photo, description)
- **Contact**: Key-value pairs for contact information
- **WorkExperience**: Job history with dates and details
- **Education**: Educational background
- **OtherExperience**: Flexible categories for additional experiences
- **FeatureMap**: Custom JSON data for any experience type
- **Template**: Go templates for resume rendering

### Workflow

1. AI agent creates resume using MCP tools
2. Agent adds contacts, experiences, and feature maps
3. Agent creates templates for rendering
4. Agent generates HTML previews with custom styling
5. User visits preview URL to view rendered resume

## Development

### Prerequisites

- Go 1.21 or higher
- Make

### Commands

```bash
# Build the project
make build

# Run tests
make test

# Run the MCP server
make run

# Clean build artifacts
make clean

# Install locally (requires sudo)
make install-local

# Create package (macOS)
make package
```

### Testing

The project includes comprehensive unit tests for all MCP tools:

```bash
# Run all tests
make test

# Run tests with verbose output
go test ./tools -v

# Run specific test
go test ./tools -run TestCreateResumeTool
```

## API

### HTTP Endpoints

The server runs an HTTP API on port 8080 with the following endpoints:

- `GET /resume/preview/:sid` - View generated HTML preview with download button
- `GET /resume/download/:sid` - Download resume as PDF (pixel-perfect with preview)
- `GET /health` - Health check endpoint

### PDF Generation

The server supports PDF generation using headless Chrome via chromedp:

- **Local Mode**: Uses locally installed Chrome/Chromium (automatic)
- **Remote Mode**: Connect to remote Chrome instance via WebSocket

#### Environment Variables

- `CHROMEDP_REMOTE_URL` - Optional WebSocket URL for remote Chrome (e.g., `ws://chromedp:9222`)
- `BASE_URL` - Base URL for generating preview/download links (e.g., `https://resume.example.com`)

#### PDF Features

- Generated PDF is pixel-perfect with web preview
- Uses the same HTML and CSS styling
- Download button excluded from PDF output
- Stateless - PDF generated in-memory, no local files created

### Kubernetes Deployment

When deploying to Kubernetes, the deployment includes a Chrome sidecar container:

```yaml
# Chrome runs as a sidecar on port 9222
- name: chrome
  image: zenika/alpine-chrome:latest
  args:
    - --headless
    - --no-sandbox
    - --disable-dev-shm-usage
    - --remote-debugging-port=9222
```

Set `CHROMEDP_REMOTE_URL=ws://localhost:9222` to use the sidecar.

### Database

- **Type**: SQLite
- **Location**: `resume.db` (created automatically)
- **Migrations**: Automatic on startup

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `make test`
5. Build the project: `make build`
6. Submit a pull request

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Support

For issues and questions, please visit: https://github.com/rxtech-lab/resume-mcp/issues