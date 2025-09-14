package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/mandacode-com/merver"
	"github.com/mandacode-com/serengeti-integrated/cmd/shared"
	"github.com/mandacode-com/serengeti-integrated/configs"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := configs.LoadAuthConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Setup logger
	logger := shared.SetupLogger(cfg.Env)
	logger.Info().Msg("Starting auth service")

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
	authGroup := v1.Group("/auth")

	identifyUserHandler := adapter.ProvideIdentifyUserHandler()
	identifyUserHandler.RegisterRoutes(authGroup)

	// Start server
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

