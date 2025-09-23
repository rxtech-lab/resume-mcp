package api

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/resume-mcp/internal/database"
	"github.com/rxtech-lab/resume-mcp/internal/service"
	types "github.com/rxtech-lab/resume-mcp/internal/types"

	auth "github.com/rxtech-lab/mcprouter-authenticator/authenticator"
	auth2 "github.com/rxtech-lab/mcprouter-authenticator/middleware"
	authTypes "github.com/rxtech-lab/mcprouter-authenticator/types"
)

type APIServer struct {
	app              *fiber.App
	db               *database.Database
	templateService  *service.TemplateService
	streamableServer *server.StreamableHTTPServer
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

	return server
}

func (s *APIServer) SetupStreamableServer(server *server.StreamableHTTPServer) {
	s.streamableServer = server

	mcprouterAuthenticator := auth.NewApikeyAuthenticator(os.Getenv("MCPROUTER_SERVER_URL"), http.DefaultClient)
	// setup middleware
	s.app.Use(auth2.FiberApikeyMiddleware(mcprouterAuthenticator, os.Getenv("MCPROUTER_SERVER_API_KEY"), func(c *fiber.Ctx, user *authTypes.User) error {
		// Store user in context for later use - adapt types.User to utils.AuthenticatedUser
		authenticatedUser := &types.AuthenticatedUser{
			Sub:   user.ID,
			Roles: []string{user.Role}, // Map single role to roles array
		}
		// Store the adapted AuthenticatedUser instead of the raw types.User
		c.Locals(types.AuthenticatedUserContextKey, authenticatedUser)
		return nil
	}))
}

func (s *APIServer) SetupRoutes() {
	// add health check
	s.app.Get("/health", s.handleHealth)
	s.app.Get("/resume/preview/:sessionId", s.handlePreview)
	if s.streamableServer != nil {
		s.app.All("/mcp", s.createAuthenticatedMCPHandler(s.streamableServer))
	}
}

func (s *APIServer) createAuthenticatedMCPHandler(streamableServer *server.StreamableHTTPServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals(types.AuthenticatedUserContextKey)
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		authenticatedUser := user.(*types.AuthenticatedUser)
		// Create a custom HTTP handler that injects authentication context
		httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Add authenticated user to context
			if authenticatedUser != nil {
				ctx = types.WithAuthenticatedUser(ctx, authenticatedUser)
			}

			// Update request with authenticated context
			r = r.WithContext(ctx)

			// Forward to the actual MCP streamable server
			streamableServer.ServeHTTP(w, r)
		})

		// Use Fiber's adaptor to convert
		return adaptor.HTTPHandler(httpHandler)(c)
	}
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

func (s *APIServer) Start(port string) (string, error) {
	// Create a listener on the specified port (0 means any available port)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return "", err
	}

	// Get the actual port assigned
	actualPort := listener.Addr().(*net.TCPAddr).Port
	log.Printf("Starting API server on port %d", actualPort)

	// Start serving in a goroutine
	go func() {
		if err := s.app.Listener(listener); err != nil {
			log.SetOutput(os.Stderr)
			log.SetFlags(0)
			log.Printf("API server error: %v", err)
			log.SetOutput(io.Discard)
		}
	}()

	return fmt.Sprintf("%d", actualPort), nil
}

func (s *APIServer) Shutdown() error {
	return s.app.Shutdown()
}
