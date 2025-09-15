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
	_ "github.com/mandacode-com/mandacode-ssam/docs/core"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           MandaCode Service Hub Core API
// @version         1.0
// @description     Multi-tenant service management and client application API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

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
