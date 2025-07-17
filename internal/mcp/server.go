package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/service"
	"github.com/rxtech-lab/resume-mcp/tools"
)

type MCPServer struct {
	server *server.MCPServer
}

func NewMCPServer(db *database.Database, port string, templateService *service.TemplateService) *MCPServer {
	mcpServer := &MCPServer{}
	mcpServer.InitializeTools(db, port, templateService)
	return mcpServer
}

func (s *MCPServer) InitializeTools(db *database.Database, port string, templateService *service.TemplateService) {
	srv := server.NewMCPServer(
		"Resume MCP Server",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	// Initialize all tools
	createResumeTool, createResumeHandler := tools.NewCreateResumeTool(db)
	srv.AddTool(createResumeTool, createResumeHandler)

	updateBasicInfoTool, updateBasicInfoHandler := tools.NewUpdateBasicInfoTool(db)
	srv.AddTool(updateBasicInfoTool, updateBasicInfoHandler)

	addContactInfoTool, addContactInfoHandler := tools.NewAddContactInfoTool(db)
	srv.AddTool(addContactInfoTool, addContactInfoHandler)

	addWorkExperienceTool, addWorkExperienceHandler := tools.NewAddWorkExperienceTool(db)
	srv.AddTool(addWorkExperienceTool, addWorkExperienceHandler)

	addEducationTool, addEducationHandler := tools.NewAddEducationTool(db)
	srv.AddTool(addEducationTool, addEducationHandler)

	addOtherExperienceTool, addOtherExperienceHandler := tools.NewAddOtherExperienceTool(db)
	srv.AddTool(addOtherExperienceTool, addOtherExperienceHandler)

	addFeatureMapTool, addFeatureMapHandler := tools.NewAddFeatureMapTool(db)
	srv.AddTool(addFeatureMapTool, addFeatureMapHandler)

	updateFeatureMapTool, updateFeatureMapHandler := tools.NewUpdateFeatureMapTool(db)
	srv.AddTool(updateFeatureMapTool, updateFeatureMapHandler)

	deleteFeatureMapTool, deleteFeatureMapHandler := tools.NewDeleteFeatureMapTool(db)
	srv.AddTool(deleteFeatureMapTool, deleteFeatureMapHandler)

	getResumeByNameTool, getResumeByNameHandler := tools.NewGetResumeByNameTool(db)
	srv.AddTool(getResumeByNameTool, getResumeByNameHandler)

	listResumesTool, listResumesHandler := tools.NewListResumesTool(db)
	srv.AddTool(listResumesTool, listResumesHandler)

	deleteResumeTool, deleteResumeHandler := tools.NewDeleteResumeTool(db)
	srv.AddTool(deleteResumeTool, deleteResumeHandler)

	generatePreviewTool, generatePreviewHandler := tools.NewGeneratePreviewTool(db, port, templateService)
	srv.AddTool(generatePreviewTool, generatePreviewHandler)

	updatePreviewStyleTool, updatePreviewStyleHandler := tools.NewUpdatePreviewStyleTool(db, port)
	srv.AddTool(updatePreviewStyleTool, updatePreviewStyleHandler)

	// Template tools
	createTemplateTool, createTemplateHandler := tools.NewCreateTemplateTool(db, templateService)
	srv.AddTool(createTemplateTool, createTemplateHandler)

	getTemplateTool, getTemplateHandler := tools.NewGetTemplateTool(db)
	srv.AddTool(getTemplateTool, getTemplateHandler)

	listTemplatesTool, listTemplatesHandler := tools.NewListTemplatesTool(db)
	srv.AddTool(listTemplatesTool, listTemplatesHandler)

	updateTemplateTool, updateTemplateHandler := tools.NewUpdateTemplateTool(db, templateService)
	srv.AddTool(updateTemplateTool, updateTemplateHandler)

	deleteTemplateTool, deleteTemplateHandler := tools.NewDeleteTemplateTool(db)
	srv.AddTool(deleteTemplateTool, deleteTemplateHandler)

	getResumeContextTool, getResumeContextHandler := tools.NewGetResumeContextTool(db)
	srv.AddTool(getResumeContextTool, getResumeContextHandler)

	s.server = srv
}

func (s *MCPServer) Start() error {
	return server.ServeStdio(s.server)
}
