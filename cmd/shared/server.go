package shared

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	merrmid "github.com/mandacode-com/merr/middleware"
	mervermid "github.com/mandacode-com/merver/middleware"
	"github.com/mandacode-com/mandacode-service-hub/configs"
	"github.com/rs/zerolog"
)

type HTTPServer struct {
	http   *http.Server
	engine *gin.Engine
	logger *zerolog.Logger
}

// Run implements merver.Server.
func (s *HTTPServer) Run(ctx context.Context) error {
	s.logger.WithLevel(SystemLevel).Str("addr", s.http.Addr).Msg("Starting HTTP server")
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error().Err(err).Msg("Failed to start HTTP server")
		return err
	}
	return nil
}

// Stop implements merver.Server.
func (s *HTTPServer) Stop(ctx context.Context) error {
	s.logger.WithLevel(SystemLevel).Msg("Stopping HTTP server")
	if err := s.http.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("Failed to gracefully shutdown HTTP server")
		return err
	}
	s.logger.WithLevel(SystemLevel).Msg("HTTP server stopped gracefully")
	return nil
}

func (s *HTTPServer) GetEngine() *gin.Engine {
	return s.engine
}

func SetupServer(cfg *configs.ServerConfig, logger *zerolog.Logger) *HTTPServer {
	engine := gin.New()

	// Middleware stack
	engine.Use(mervermid.GinZeroLogger(logger))
	engine.Use(merrmid.GinErrorHandler())
	engine.Use(gin.Recovery())

	// Health endpoints
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204)
	})

	httpServer := &http.Server{
		Addr:    cfg.GenerateAddress(),
		Handler: engine,
	}

	return &HTTPServer{
		http:   httpServer,
		engine: engine,
		logger: logger,
	}
}

