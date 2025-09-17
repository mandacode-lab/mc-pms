package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
)

const (
	PermissionKey string = "permission"
	UserIDKey     string = "user_id"
)

// RequirePermission creates a middleware that checks if the user has the required permission
func RequirePermission(iamService out.IAMService, resource out.Resource, action out.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Extract user ID from X-User-ID header
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			err := merr.New(merr.ErrUnauthorized, "missing X-User-ID header", nil)
			c.Error(err)
			c.Abort()
			return
		}

		// Create permission object
		permission := out.Permission{
			Resource: resource,
			Action:   action,
			Context:  buildPermissionContext(c),
		}

		// Validate permission structure
		if !permission.Validate() {
			err := merr.New(merr.ErrBadRequest, "invalid permission structure", nil)
			c.Error(err)
			c.Abort()
			return
		}

		// Check permission with IAM service
		permissionInfo, err := iamService.HasPermission(ctx, userID, permission)
		if err != nil {
			err := merr.New(merr.ErrInternalServerError, "failed to check permission", err)
			c.Error(err)
			c.Abort()
			return
		}

		// Check if permission is allowed
		if !permissionInfo.Allowed {
			err := merr.New(merr.ErrForbidden, "insufficient permissions", nil)
			c.Error(err)
			c.Abort()
			return
		}

		// Store user ID and permission info in context for use in handlers
		c.Set(UserIDKey, userID)
		c.Set(PermissionKey, permissionInfo)

		c.Next()
	}
}

// buildPermissionContext builds context information for permission checking
func buildPermissionContext(c *gin.Context) map[string]any {
	context := make(map[string]any)

	// Add client IP
	context["client_ip"] = c.ClientIP()

	return context
}
