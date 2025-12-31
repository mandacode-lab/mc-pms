package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-pms/internal/adapter/http/handler/client_ns"
	"github.com/mandacode-com/mandacode-pms/internal/adapter/http/handler/client_project"
	"github.com/mandacode-com/mandacode-pms/internal/adapter/http/middleware"
)

type ClientRouterConfig struct {
	NSHandler      *client_ns.Handler
	ProjectHandler *client_project.Handler
	AuthMiddleware *middleware.AuthMiddleware
}

func SetupClientRouter(r *gin.Engine, config *ClientRouterConfig) {
	// Apply auth middleware to all client routes
	client := r.Group("/client")
	client.Use(config.AuthMiddleware.RequireAuth())

	// Namespace routes
	namespaces := client.Group("/namespaces")
	{
		// Get namespace by ID
		namespaces.GET("/:id", config.NSHandler.GetByID)
	}

	// Project routes
	projects := client.Group("/projects")
	{
		// Get project by ID
		projects.GET("/:id", config.ProjectHandler.GetByID)

		// Validate project-namespace relationship
		projects.GET("/validate", config.ProjectHandler.IsValidProjectNamespace)
	}
}
