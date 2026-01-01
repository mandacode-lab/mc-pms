package main

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/mandacode-lab/mc-pms/ent"
	"github.com/mandacode-lab/mc-pms/internal/adapter/http/handler/admin_ns"
	"github.com/mandacode-lab/mc-pms/internal/adapter/http/handler/admin_project"
	"github.com/mandacode-lab/mc-pms/internal/adapter/http/middleware"
	"github.com/mandacode-lab/mc-pms/internal/adapter/iam"
	"github.com/mandacode-lab/mc-pms/internal/adapter/idgen"
	"github.com/mandacode-lab/mc-pms/internal/adapter/persistence"
	adminNSApp "github.com/mandacode-lab/mc-pms/internal/app/admin_ns"
	adminProjectApp "github.com/mandacode-lab/mc-pms/internal/app/admin_project"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

// Container holds all dependencies
type Container struct {
	// Infrastructure
	EntClient *ent.Client
	Logger    *zerolog.Logger

	// Adapters
	NSRepo          *persistence.NamespaceRepository
	ProjectRepo     *persistence.ProjectRepository
	NSIDGenerator   *idgen.ULIDGenerator
	ProjIDGenerator *idgen.ULIDGenerator
	IAMService      *iam.NoopIAMService

	// Applications
	AdminNSApp      *adminNSApp.Application
	AdminProjectApp *adminProjectApp.Application

	// HTTP
	AuthMiddleware *middleware.AuthMiddleware
	NSHandler      *admin_ns.Handler
	ProjectHandler *admin_project.Handler
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

	// ID Generators
	c.NSIDGenerator = idgen.NewULIDGenerator()
	c.ProjIDGenerator = idgen.NewULIDGenerator()

	// IAM Service (noop for now)
	c.IAMService = iam.NewNoopIAMService()

	c.Logger.Info().Msg("adapters initialized")
}

func (c *Container) setupApplications() {
	// Admin Namespace Application
	c.AdminNSApp = adminNSApp.NewApplication(
		c.NSRepo,
		c.NSIDGenerator,
	)

	// Admin Project Application
	c.AdminProjectApp = adminProjectApp.NewApplication(
		c.ProjectRepo,
		c.ProjIDGenerator,
	)

	c.Logger.Info().Msg("applications initialized")
}

func (c *Container) setupHTTP() {
	// Middleware
	c.AuthMiddleware = middleware.NewAuthMiddleware(c.IAMService)

	// Handlers
	c.NSHandler = admin_ns.NewHandler(c.AdminNSApp)
	c.ProjectHandler = admin_project.NewHandler(c.AdminProjectApp)

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
