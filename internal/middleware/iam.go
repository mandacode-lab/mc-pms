package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/merr"
	"github.com/rs/zerolog/log"
)

const (
	PermissionKey string = "permission"
	UserIDKey     string = "user_id"
)

type PermissionMiddleware struct {
	IAMService out.IAMService
	Enabled    bool
}

func NewPermissionMiddleware(iamService out.IAMService, enabled bool) *PermissionMiddleware {
	return &PermissionMiddleware{
		IAMService: iamService,
		Enabled:    enabled,
	}
}

// RequirePermission creates a middleware that checks if the user has the required permission
func (m *PermissionMiddleware) RequirePermission(resource out.Resource, action out.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If IAM is disabled, skip all permission checks and allow the request
		if !m.Enabled {
			log.Warn().Msg("IAM is disabled; skipping permission checks")
			// Set a default user ID if none provided when IAM is disabled
			userID := c.GetHeader("X-User-ID")
			if userID == "" {
				userID = "system" // Default user when IAM is disabled
			}
			c.Set(UserIDKey, userID)
			c.Set(PermissionKey, true)
			c.Next()
			return
		}

		ctx := c.Request.Context()

		// Extract user ID from X-User-ID header
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			err := merr.New(merr.ErrUnauthorized, "missing X-User-ID header", nil)
			_ = c.Error(err)
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
			_ = c.Error(err)
			c.Abort()
			return
		}

		// Check permission with IAM service
		allowed, err := m.IAMService.HasPermission(ctx, userID, permission)
		if err != nil {
			err := merr.New(merr.ErrInternalServerError, "failed to check permission", err)
			_ = c.Error(err)
			c.Abort()
			return
		}

		// Check if permission is allowed
		if !allowed {
			err := merr.New(merr.ErrForbidden, "insufficient permissions", nil)
			_ = c.Error(err)
			c.Abort()
			return
		}

		// Store user ID and permission info in context for use in handlers
		c.Set(UserIDKey, userID)
		c.Set(PermissionKey, allowed)

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
