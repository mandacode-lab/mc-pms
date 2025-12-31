package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-pms/internal/port/iam"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
	ContextKeyToken     = "token"
)

type AuthMiddleware struct {
	iamService iam.IAMService
}

func NewAuthMiddleware(iamService iam.IAMService) *AuthMiddleware {
	return &AuthMiddleware{
		iamService: iamService,
	}
}

// RequireAuth extracts the token from the Authorization header and stores it in the context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, BearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, BearerPrefix)
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		// Store token in context for later use
		c.Set(ContextKeyToken, token)
		c.Next()
	}
}

// RequirePermission checks if the token holder has permission to perform the action on the resource
func (m *AuthMiddleware) RequirePermission(action iam.ActionType, resourceType iam.ResourceType, getResourceID func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, exists := c.Get(ContextKeyToken)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found in context"})
			c.Abort()
			return
		}

		tokenStr, ok := token.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid token type in context"})
			c.Abort()
			return
		}

		resourceID := getResourceID(c)
		resource := iam.NewResource(resourceType, resourceID)

		req := iam.NewAuthorizationRequest(tokenStr, action, resource)
		result, err := m.iamService.CheckAuthorization(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "authorization check failed"})
			c.Abort()
			return
		}

		if !result.Allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission denied", "reason": result.Reason})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetTokenFromContext retrieves the token from the gin context
func GetTokenFromContext(c *gin.Context) (string, bool) {
	token, exists := c.Get(ContextKeyToken)
	if !exists {
		return "", false
	}

	tokenStr, ok := token.(string)
	return tokenStr, ok
}
