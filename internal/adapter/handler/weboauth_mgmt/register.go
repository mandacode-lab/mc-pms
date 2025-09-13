package weboauth_mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

type RegisterWebOAuthRequest struct {
	ClientAppID     string   `json:"client_app_id" binding:"required"`
	Provider        string   `json:"provider" binding:"required"`
	OAuthClientID   string   `json:"oauth_client_id" binding:"required"`
	OAuthSecret     string   `json:"oauth_secret" binding:"required"`
	RedirectURI     string   `json:"redirect_uri" binding:"required"`
	Scopes          []string `json:"scopes" binding:"required"`
}

func (h *Handler) RegisterWebOAuth(c *gin.Context) {
	ctx := c.Request.Context()

	var req RegisterWebOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientAppPublicID, err := clientappval.ParsePublicID(req.ClientAppID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client_app_id format"})
		return
	}

	provider := shared.Provider(req.Provider)
	if !provider.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider"})
		return
	}

	cmd := &in.RegisterWebOAuthCommand{
		ClientAppID:   clientAppPublicID,
		Provider:      provider,
		OAuthClientID: req.OAuthClientID,
		OAuthSecret:   []byte(req.OAuthSecret),
		RedirectURI:   req.RedirectURI,
		Scopes:        req.Scopes,
	}

	err = h.webOAuthMgmt.RegisterWebOAuth(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "WebOAuth registered successfully"})
}