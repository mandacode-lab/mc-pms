package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mandacode-com/merver"
	"github.com/mandacode-com/mandacode-ssam/cmd/shared"
	"github.com/mandacode-com/mandacode-ssam/configs"
	_ "github.com/mandacode-com/mandacode-ssam/docs/management"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           MandaCode Service Hub Management API
// @version         1.0
// @description     Multi-tenant service and client application management API

// @contact.name   API Support

// @BasePath  /v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.


func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load .env file if not production
	if os.Getenv("ENV") != "prod" {
		if err := godotenv.Load(".env.dev.management"); err != nil {
			log.Warn().Err(err).Msg("Could not load .env.dev.management file, using system environment variables")
		}
	}

	// Load configuration
	cfg, err := configs.LoadManagementConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Setup logger
	logger := shared.SetupLogger(cfg.Env)
	logger.Info().Msg("Starting management service")

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

	// Swagger documentation
	srv.GetEngine().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
