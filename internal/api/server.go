package api

import (
	"io"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/service"
)

type APIServer struct {
	app             *fiber.App
	db              *database.Database
	templateService *service.TemplateService
}

func NewAPIServer(db *database.Database, templateService *service.TemplateService) *APIServer {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			log.SetOutput(os.Stderr)
			log.SetFlags(0)
			log.Printf("API Error: %v", err)
			log.SetOutput(io.Discard)
			return c.Status(500).JSON(fiber.Map{
				"error": "Internal server error",
			})
		},
	})

	app.Use(cors.New())

	server := &APIServer{
		app:             app,
		db:              db,
		templateService: templateService,
	}

	server.setupRoutes()
	return server
}

func (s *APIServer) setupRoutes() {
	s.app.Get("/resume/preview/:sessionId", s.handlePreview)
	s.app.Get("/health", s.handleHealth)
}

func (s *APIServer) handlePreview(c *fiber.Ctx) error {
	sessionID := c.Params("sessionId")

	session, err := s.db.GetPreviewSession(sessionID)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.SetFlags(0)
		log.Printf("Preview session not found: %v", err)
		log.SetOutput(io.Discard)
		return c.Status(404).JSON(fiber.Map{
			"error": "Preview session not found",
		})
	}

	fullHTML, err := s.templateService.GeneratePreview(session.Template, session.CSS, session.Resume)

	// set content type to html
	c.Set("Content-Type", "text/html")
	return c.SendString(fullHTML)
}

func (s *APIServer) handleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "resume-mcp",
	})
}

func (s *APIServer) Start(port string) error {
	log.Printf("Starting API server on port %s", port)
	return s.app.Listen(":" + port)
}

func (s *APIServer) Shutdown() error {
	return s.app.Shutdown()
}
