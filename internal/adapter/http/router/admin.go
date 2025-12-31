package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-pms/internal/adapter/http/handler/admin_ns"
	"github.com/mandacode-com/mandacode-pms/internal/adapter/http/handler/admin_project"
	"github.com/mandacode-com/mandacode-pms/internal/adapter/http/middleware"
	"github.com/mandacode-com/mandacode-pms/internal/port/iam"
)

type AdminRouterConfig struct {
	NSHandler      *admin_ns.Handler
	ProjectHandler *admin_project.Handler
	AuthMiddleware *middleware.AuthMiddleware
}

func SetupAdminRouter(r *gin.Engine, config *AdminRouterConfig) {
	// Apply auth middleware to all admin routes
	admin := r.Group("/admin")
	admin.Use(config.AuthMiddleware.RequireAuth())

	// Namespace routes
	namespaces := admin.Group("/namespaces")
	{
		// Create namespace - requires permission for all namespaces
		namespaces.POST("",
			config.AuthMiddleware.RequirePermission(
				iam.ActionNamespaceCreate,
				iam.ResourceTypeNamespace,
				func(c *gin.Context) string { return "*" },
			),
			config.NSHandler.Create,
		)

		// Get namespace by ID
		namespaces.GET("/:id",
			config.AuthMiddleware.RequirePermission(
				iam.ActionNamespaceRead,
				iam.ResourceTypeNamespace,
				func(c *gin.Context) string { return c.Param("id") },
			),
			config.NSHandler.GetByID,
		)

		// Update namespace
		namespaces.POST("/:id",
			config.AuthMiddleware.RequirePermission(
				iam.ActionNamespaceUpdate,
				iam.ResourceTypeNamespace,
				func(c *gin.Context) string { return c.Param("id") },
			),
			config.NSHandler.Update,
		)

		// Delete namespace
		namespaces.DELETE("/:id",
			config.AuthMiddleware.RequirePermission(
				iam.ActionNamespaceDelete,
				iam.ResourceTypeNamespace,
				func(c *gin.Context) string { return c.Param("id") },
			),
			config.NSHandler.Delete,
		)

		// Find/List namespaces - query params determine if it's find or list all
		namespaces.GET("",
			config.AuthMiddleware.RequirePermission(
				iam.ActionNamespaceList,
				iam.ResourceTypeNamespace,
				func(c *gin.Context) string { return "*" },
			),
			func(c *gin.Context) {
				// If name query param exists, use Find, otherwise ListAll
				if c.Query("name") != "" {
					config.NSHandler.Find(c)
				} else {
					config.NSHandler.ListAll(c)
				}
			},
		)
	}

	// Project routes
	projects := admin.Group("/projects")
	{
		// Create project
		projects.POST("",
			config.AuthMiddleware.RequirePermission(
				iam.ActionProjectCreate,
				iam.ResourceTypeProject,
				func(c *gin.Context) string { return "*" },
			),
			config.ProjectHandler.Create,
		)

		// Get project by ID
		projects.GET("/:id",
			config.AuthMiddleware.RequirePermission(
				iam.ActionProjectRead,
				iam.ResourceTypeProject,
				func(c *gin.Context) string { return c.Param("id") },
			),
			config.ProjectHandler.GetByID,
		)

		// Update project
		projects.POST("/:id",
			config.AuthMiddleware.RequirePermission(
				iam.ActionProjectUpdate,
				iam.ResourceTypeProject,
				func(c *gin.Context) string { return c.Param("id") },
			),
			config.ProjectHandler.Update,
		)

		// Delete project
		projects.DELETE("/:id",
			config.AuthMiddleware.RequirePermission(
				iam.ActionProjectDelete,
				iam.ResourceTypeProject,
				func(c *gin.Context) string { return c.Param("id") },
			),
			config.ProjectHandler.Delete,
		)

		// Find/List projects - query params determine if it's find or list all
		projects.GET("",
			config.AuthMiddleware.RequirePermission(
				iam.ActionProjectList,
				iam.ResourceTypeProject,
				func(c *gin.Context) string { return "*" },
			),
			func(c *gin.Context) {
				// If name or namespace_id query params exist, use Find, otherwise ListAll
				if c.Query("name") != "" || c.Query("namespace_id") != "" {
					config.ProjectHandler.Find(c)
				} else {
					config.ProjectHandler.ListAll(c)
				}
			},
		)
	}
}
