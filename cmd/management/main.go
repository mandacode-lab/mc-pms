package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"github.com/mandacode-com/mandacode-ssam/cmd/shared"
	"github.com/mandacode-com/mandacode-ssam/configs"
	_ "github.com/mandacode-com/mandacode-ssam/docs/management"
	"github.com/mandacode-com/merver"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           MandaCode Service Hub Management API
// @version         1.0
// @description     Multi-tenant service and client application management API

// @contact.name   API Support

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
	defer func() {
		if closeErr := adapter.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("Failed to close adapter")
		}
	}()

	// Setup CORS configuration for management API
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // 12 hours
	}

	// Setup server
	server := shared.SetupServer(&cfg.Server, &logger, cors.New(corsConfig))

	// Register routes
	engine := server.GetEngine()

	// Service management routes
	serviceGroup := engine.Group("/services")
	serviceMgmtHandler := adapter.ProvideServiceMgmtHandler()
	serviceMgmtHandler.RegisterRoutes(serviceGroup)

	// Client app management routes
	clientAppGroup := engine.Group("/client-apps")
	clientAppMgmtHandler := adapter.ProvideClientAppMgmtHandler()
	clientAppMgmtHandler.RegisterRoutes(clientAppGroup)

	// Swagger documentation
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChan
		log.WithLevel(shared.SystemLevel).Msgf("received shutdown signal: %s", sig.String())
		cancel()
	}()

	srvGroup := merver.NewServerGroup(server)
	if err := srvGroup.Run(ctx); err != nil {
		logger.Fatal().Str("error", err.Error()).Msg("failed to run server group")
	}
}
