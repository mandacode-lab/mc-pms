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
	// Apply auth middleware to all clientv1 routes
	clientv1 := r.Group("/v1")
	clientv1.Use(config.AuthMiddleware.RequireAuth())

	// Namespace routes
	namespaces := clientv1.Group("/namespaces")
	{
		// Get namespace by ID
		namespaces.GET("/:id", config.NSHandler.GetByID)
	}

	// Project routes
	projects := clientv1.Group("/projects")
	{
		// Get project by ID
		projects.GET("/:id", config.ProjectHandler.GetByID)

		// Validate project-namespace relationship
		projects.GET("/validate", config.ProjectHandler.IsValidProjectNamespace)
	}
}
