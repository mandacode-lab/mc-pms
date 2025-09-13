package identify_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

type IdentifyByTokenRequest struct {
	Provider        string `json:"provider" binding:"required"`
	Token           string `json:"token" binding:"required"`
	ClientAppID     string `json:"client_app_id" binding:"required"`
	ClientAppSecret string `json:"client_app_secret" binding:"required"`
}

func (h *Handler) IdentifyByToken(c *gin.Context) {
	ctx := c.Request.Context()

	// Get parameters from query string
	providerStr := c.Query("provider")
	token := c.Query("token")
	clientAppID := c.Query("client_app_id")
	clientAppSecret := c.Query("client_app_secret")

	// Validate required parameters
	if providerStr == "" || token == "" || clientAppID == "" || clientAppSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters: provider, token, client_app_id, client_app_secret"})
		return
	}

	provider := shared.Provider(providerStr)
	if !provider.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
		return
	}

	clientAppPublicID, err := clientappval.ParsePublicID(clientAppID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client_app_id format"})
		return
	}

	cmd := &in.IdentifyByToken{
		Provider:        provider,
		Token:           token,
		ClientApp:       clientAppPublicID,
		ClientAppSecret: []byte(clientAppSecret),
	}

	view, err := h.identifyUser.IdentifyByToken(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, view)
}