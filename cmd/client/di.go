package main

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/mandacode-lab/mc-pms/ent"
	"github.com/mandacode-lab/mc-pms/internal/adapter/http/handler/client_ns"
	"github.com/mandacode-lab/mc-pms/internal/adapter/http/handler/client_project"
	"github.com/mandacode-lab/mc-pms/internal/adapter/http/middleware"
	"github.com/mandacode-lab/mc-pms/internal/adapter/iam"
	"github.com/mandacode-lab/mc-pms/internal/adapter/persistence"
	clientNSApp "github.com/mandacode-lab/mc-pms/internal/app/client_ns"
	clientProjectApp "github.com/mandacode-lab/mc-pms/internal/app/client_project"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

// Container holds all dependencies
type Container struct {
	// Infrastructure
	EntClient *ent.Client
	Logger    *zerolog.Logger

	// Adapters
	NSRepo      *persistence.NamespaceRepository
	ProjectRepo *persistence.ProjectRepository
	IAMService  *iam.NoopIAMService

	// Applications
	ClientNSApp      *clientNSApp.Application
	ClientProjectApp *clientProjectApp.Application

	// HTTP
	AuthMiddleware *middleware.AuthMiddleware
	NSHandler      *client_ns.Handler
	ProjectHandler *client_project.Handler
}

// NewContainer creates and wires all dependencies
func NewContainer(cfg *Config, logger *zerolog.Logger) (*Container, error) {
	c := &Container{
		Logger: logger,
	}

	// Setup database
	if err := c.setupDatabase(cfg); err != nil {
		return nil, err
	}

	// Setup adapters
	c.setupAdapters()

	// Setup applications
	c.setupApplications()

	// Setup HTTP layer
	c.setupHTTP()

	return c, nil
}

func (c *Container) setupDatabase(cfg *Config) error {
	// Open database connection
	drv, err := sql.Open("postgres", cfg.DB.Address())
	if err != nil {
		return err
	}

	// Create ent client
	c.EntClient = ent.NewClient(ent.Driver(drv))

	// Run migrations
	ctx := context.Background()
	if err := c.EntClient.Schema.Create(ctx); err != nil {
		return err
	}

	c.Logger.Info().Msg("database connection established and migrations applied")
	return nil
}

func (c *Container) setupAdapters() {
	// Repositories
	c.NSRepo = persistence.NewNamespaceRepository(c.EntClient)
	c.ProjectRepo = persistence.NewProjectRepository(c.EntClient)

	// IAM Service (noop for now)
	c.IAMService = iam.NewNoopIAMService()

	c.Logger.Info().Msg("adapters initialized")
}

func (c *Container) setupApplications() {
	// Client Namespace Application
	c.ClientNSApp = clientNSApp.NewApplication(c.NSRepo)

	// Client Project Application
	c.ClientProjectApp = clientProjectApp.NewApplication(c.ProjectRepo)

	c.Logger.Info().Msg("applications initialized")
}

func (c *Container) setupHTTP() {
	// Middleware
	c.AuthMiddleware = middleware.NewAuthMiddleware(c.IAMService)

	// Handlers
	c.NSHandler = client_ns.NewHandler(c.ClientNSApp)
	c.ProjectHandler = client_project.NewHandler(c.ClientProjectApp)

	c.Logger.Info().Msg("HTTP layer initialized")
}

// Close cleans up resources
func (c *Container) Close() error {
	if c.EntClient != nil {
		c.Logger.Info().Msg("closing database connection")
		return c.EntClient.Close()
	}
	return nil
}
