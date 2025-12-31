package http

import (
	"context"
	stdhttp "net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-pms/internal/config"
	merrmid "github.com/mandacode-com/merr/middleware"
	mervermid "github.com/mandacode-com/merver/middleware"
	"github.com/rs/zerolog"
)

// HTTPServer wraps gin engine with merver.Server interface
type HTTPServer struct {
	http   *stdhttp.Server
	engine *gin.Engine
	logger *zerolog.Logger
}

// Run implements merver.Server
func (s *HTTPServer) Run(ctx context.Context) error {
	s.logger.Info().Str("addr", s.http.Addr).Msg("starting HTTP server")
	if err := s.http.ListenAndServe(); err != nil && err != stdhttp.ErrServerClosed {
		s.logger.Error().Err(err).Msg("failed to start HTTP server")
		return err
	}
	return nil
}

// Stop implements merver.Server
func (s *HTTPServer) Stop(ctx context.Context) error {
	s.logger.Info().Msg("stopping HTTP server")
	if err := s.http.Shutdown(ctx); err != nil {
		s.logger.Error().Err(err).Msg("failed to gracefully shutdown HTTP server")
		return err
	}
	s.logger.Info().Msg("HTTP server stopped gracefully")
	return nil
}

// GetEngine returns the underlying gin.Engine for route registration
func (s *HTTPServer) GetEngine() *gin.Engine {
	return s.engine
}

// New creates a new HTTP server with the given configuration
func New(cfg *config.HTTP, logger *zerolog.Logger) *HTTPServer {
	engine := gin.New()

	// Disable automatic redirects
	engine.RedirectTrailingSlash = false
	engine.RedirectFixedPath = false

	// Setup CORS if enabled
	if cfg.CORS.Enabled {
		corsConfig := cors.Config{
			AllowOrigins:     cfg.CORS.AllowedOrigins,
			AllowMethods:     cfg.CORS.AllowedMethods,
			AllowHeaders:     cfg.CORS.AllowedHeaders,
			ExposeHeaders:    cfg.CORS.ExposeHeaders,
			AllowCredentials: cfg.CORS.AllowCredentials,
			MaxAge:           time.Duration(cfg.CORS.MaxAge) * time.Second,
		}
		engine.Use(cors.New(corsConfig))
		logger.Info().
			Strs("allowed_origins", cfg.CORS.AllowedOrigins).
			Bool("allow_credentials", cfg.CORS.AllowCredentials).
			Msg("CORS middleware enabled")
	}

	// Middleware stack
	engine.Use(mervermid.GinZeroLogger(logger))
	engine.Use(merrmid.GinErrorHandler())
	engine.Use(gin.Recovery())

	// Health check endpoints
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	engine.GET("/ready", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ready"})
	})

	engine.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(204)
	})

	httpServer := &stdhttp.Server{
		Addr:         cfg.Address(),
		Handler:      engine,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &HTTPServer{
		http:   httpServer,
		engine: engine,
		logger: logger,
	}
}
