package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/mandacode-com/mandacode-ssam/cmd/shared"
	"github.com/mandacode-com/mandacode-ssam/configs"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	// Load .env file if not production
	if os.Getenv("ENV") != "prod" {
		if err := godotenv.Load(".env.dev.initialize"); err != nil {
			log.Warn().Err(err).Msg("Could not load .env.dev.initialize file, using system environment variables")
		}
	}

	// Load configuration
	cfg, err := configs.LoadInitializeConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Setup logger
	logger := shared.SetupLogger(cfg.Env)
	logger.Info().Msg("Starting initialize service")

	// Setup dependencies
	adapter, err := NewAdapter(ctx, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create adapter")
	}
	defer adapter.Close()

	// Run initialization
	initUsecase := adapter.ProvideInitializeUsecase()

	req := &in.InitializeSystemRequest{
		ServiceName:     cfg.SystemClient.SystemClientName,
		ClientAppName:   cfg.SystemClient.SystemClientName,
		ClientAppID:     cfg.SystemClient.SystemClientID,
		ClientAppSecret: []byte(cfg.SystemClient.SystemClientSecret),
	}

	result, err := initUsecase.InitializeSystem(ctx, req)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize system")
	}

	logger.Info().
		Str("service_id", result.ServiceID).
		Str("client_id", result.ClientAppID).
		Msg("System initialized successfully")
}

