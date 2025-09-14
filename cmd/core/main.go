package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/mandacode-com/serengeti-integrated/cmd/shared"
	"github.com/mandacode-com/serengeti-integrated/configs"
	"github.com/mandacode-com/merver"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load .env file if not production
	if os.Getenv("ENV") != "prod" {
		if err := godotenv.Load(".env.dev.core"); err != nil {
			log.Warn().Err(err).Msg("Could not load .env.dev.core file, using system environment variables")
		}
	}

	// Load configuration
	cfg, err := configs.LoadCoreConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Setup logger
	logger := shared.SetupLogger(cfg.Env)
	logger.Info().Msg("Starting core service")

	// Setup dependencies
	adapter, err := NewAdapter(ctx, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create adapter")
	}
	defer adapter.Close()

	// Setup server
	srv := shared.SetupServer(&cfg.Server, &logger)

	// Register routes
	v1 := srv.GetEngine().Group("/v1")

	// Service management routes
	serviceGroup := v1.Group("/services")
	serviceMgmtHandler := adapter.ProvideServiceMgmtHandler()
	serviceMgmtHandler.RegisterRoutes(serviceGroup)

	// Client app management routes
	clientAppGroup := v1.Group("/client-apps")
	clientAppMgmtHandler := adapter.ProvideClientAppMgmtHandler()
	clientAppMgmtHandler.RegisterRoutes(clientAppGroup)

	// WebOAuth management routes
	webOAuthGroup := v1.Group("/weboauth")
	webOAuthMgmtHandler := adapter.ProvideWebOAuthMgmtHandler()
	webOAuthMgmtHandler.RegisterRoutes(webOAuthGroup)

	// User management routes
	userGroup := v1.Group("/users")
	userMgmtHandler := adapter.ProvideUserMgmtHandler()
	userMgmtHandler.RegisterRoutes(userGroup)

	// Setup graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		log.WithLevel(shared.SystemLevel).Msgf("received shutdown signal: %s", sig.String())
		cancel()
	}()

	srvGroup := merver.NewServerGroup(srv)
	if err := srvGroup.Run(ctx); err != nil {
		logger.Fatal().Str("error", err.Error()).Msg("failed to run server group")
	}
}
