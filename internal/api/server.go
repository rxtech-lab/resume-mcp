package api

import (
	"io"
	"log"
	"os"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rxtech-lab/resume-mcp/internal/database"
)

type APIServer struct {
	app *fiber.App
	db  *database.Database
}

func NewAPIServer(db *database.Database) *APIServer {
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
		app: app,
		db:  db,
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

	tmpl, err := template.New("resume").Parse(session.Template)
	if err != nil {
		log.SetOutput(os.Stderr)
		log.SetFlags(0)
		log.Printf("Template parse error: %v", err)
		log.SetOutput(io.Discard)
		return c.Status(400).JSON(fiber.Map{
			"error": "Template parse error",
		})
	}

	var html string
	builder := &stringBuilder{}
	if err := tmpl.Execute(builder, session.Resume); err != nil {
		log.SetOutput(os.Stderr)
		log.SetFlags(0)
		log.Printf("Template execution error: %v", err)
		log.SetOutput(io.Discard)
		return c.Status(500).JSON(fiber.Map{
			"error": "Template execution error",
		})
	}
	html = builder.String()

	var cssStyle string
	if session.CSS != "" {
		cssStyle = "<style>" + session.CSS + "</style>"
	}

	fullHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Resume Preview</title>
    <script src="https://cdn.tailwindcss.com"></script>
    ` + cssStyle + `
</head>
<body>
    ` + html + `
</body>
</html>`

	c.Set("Content-Type", "text/html")
	return c.SendString(fullHTML)
}

func (s *APIServer) handleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"service": "resume-mcp",
	})
}

func (s *APIServer) Start(port string) error {
	return s.app.Listen(":" + port)
}

func (s *APIServer) Shutdown() error {
	return s.app.Shutdown()
}

type stringBuilder struct {
	content string
}

func (sb *stringBuilder) Write(p []byte) (n int, err error) {
	sb.content += string(p)
	return len(p), nil
}

func (sb *stringBuilder) String() string {
	return sb.content
}