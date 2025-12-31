package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mandacode-com/mandacode-pms/internal/adapter/http/router"
	"github.com/mandacode-com/mandacode-pms/internal/config"
	infrahttp "github.com/mandacode-com/mandacode-pms/internal/infra/http"
	"github.com/mandacode-com/merver"
	"github.com/rs/zerolog"
)

// @title MandaCode PMS Client API
// @version 1.0
// @description Client API for MandaCode Project Management System - read-only access to namespaces and projects
// @termsOfService https://mandacode.com/terms

// @contact.name API Support
// @contact.url https://mandacode.com/support
// @contact.email support@mandacode.com

// @host pms.mandacode.com
// @BasePath /api/client/v1
// @schemes http https

// @tag.name client
// @tag.description Client read-only operations for namespace and project access

func main() {
	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger.Info().Msg("starting client server")

	// Load configuration from environment variables
	cfg := &Config{}
	if err := config.Load(cfg); err != nil {
		logger.Fatal().Err(err).Msg("failed to load configuration")
	}

	// Initialize dependency container
	container, err := NewContainer(cfg, &logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize container")
	}
	defer func() {
		if err := container.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close container")
		}
	}()

	// Create HTTP server
	httpServer := infrahttp.New(&cfg.HTTP, &logger)

	// Setup routes
	router.SetupClientRouter(httpServer.GetEngine(), &router.ClientRouterConfig{
		NSHandler:      container.NSHandler,
		ProjectHandler: container.ProjectHandler,
		AuthMiddleware: container.AuthMiddleware,
	})

	logger.Info().Msg("routes configured")

	// Setup graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		logger.Info().Msgf("received shutdown signal: %s", sig.String())
		cancel()
	}()

	// Create server group and run
	srvGroup := merver.NewServerGroup(httpServer)
	if err := srvGroup.Run(ctx); err != nil {
		logger.Fatal().Err(err).Msg("failed to run server group")
	}

	logger.Info().Msg("client server stopped")
}
