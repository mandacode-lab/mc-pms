package identify_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

type IdentifyByCodeRequest struct {
	Provider        string `json:"provider" binding:"required"`
	OAuthCode       string `json:"oauth_code" binding:"required"`
	State           string `json:"state" binding:"required"`
	ClientAppID     string `json:"client_app_id" binding:"required"`
	ClientAppSecret string `json:"client_app_secret" binding:"required"`
}

func (h *Handler) IdentifyByCode(c *gin.Context) {
	ctx := c.Request.Context()

	// Get parameters from query string
	providerStr := c.Query("provider")
	oauthCode := c.Query("oauth_code")
	state := c.Query("state")
	clientAppID := c.Query("client_app_id")
	clientAppSecret := c.Query("client_app_secret")

	// Validate required parameters
	if providerStr == "" || oauthCode == "" || state == "" || clientAppID == "" || clientAppSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required parameters: provider, oauth_code, state, client_app_id, client_app_secret"})
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

	cmd := &in.IdentifyByCode{
		Provider:        provider,
		OAuthCode:       oauthCode,
		State:           state,
		ClientAppID:     clientAppPublicID,
		ClientAppSecret: []byte(clientAppSecret),
	}

	view, err := h.identifyUser.IdentifyByCode(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, view)
}