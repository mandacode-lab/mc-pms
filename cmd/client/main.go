package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mandacode-com/mandacode-ssam/cmd/shared"
	"github.com/mandacode-com/mandacode-ssam/configs"
	_ "github.com/mandacode-com/mandacode-ssam/docs/client"
	"github.com/mandacode-com/merver"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           MandaCode Service Hub Client API
// @version         1.0
// @description     Client authentication and access verification API

// @contact.name   API Support

// @securityDefinitions.basic BasicAuth
// @description Basic Authentication using ClientAppID as username and ClientSecret as password

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load .env file if not production
	if os.Getenv("ENV") != "prod" {
		if err := godotenv.Load(".env.dev.client"); err != nil {
			log.Warn().Err(err).Msg("Could not load .env.dev.client file, using system environment variables")
		}
	}

	// Load configuration
	cfg, err := configs.LoadClientConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Setup logger
	logger := shared.SetupLogger(cfg.Env)
	logger.Info().Msg("Starting client service")

	// Setup dependencies
	adapter, err := NewAdapter(ctx, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create adapter")
	}
	defer adapter.Close()

	// Setup server
	srv := shared.SetupServer(&cfg.Server, &logger)

	// Register routes
	engine := srv.GetEngine()

	// Client access routes
	clientAccessGroup := engine.Group("/client-access")
	clientAccessHandler := adapter.ProvideClientAccessHandler()
	clientAccessHandler.RegisterRoutes(clientAccessGroup)

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
