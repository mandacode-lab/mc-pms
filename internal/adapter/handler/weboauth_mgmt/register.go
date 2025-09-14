package weboauth_mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/serengeti/internal/adapter/handler/common"
	clientappval "github.com/mandacode-com/serengeti/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti/internal/domain/shared"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

type RegisterWebOAuthRequest struct {
	ClientAppID   string   `json:"client_app_id" binding:"required"`
	Provider      string   `json:"provider" binding:"required"`
	OAuthClientID string   `json:"oauth_client_id" binding:"required"`
	OAuthSecret   string   `json:"oauth_secret" binding:"required"`
	RedirectURI   string   `json:"redirect_uri" binding:"required"`
	Scopes        []string `json:"scopes" binding:"required"`
}

type RegisterWebOAuthResponse struct {
	Message string `json:"message"`
}

func toRegisterWebOAuthResponse() *RegisterWebOAuthResponse {
	return &RegisterWebOAuthResponse{
		Message: "WebOAuth registered successfully",
	}
}

// RegisterWebOAuth registers OAuth configuration for a client app
// @Summary Register OAuth config
// @Description Register OAuth configuration for a client application
// @Tags weboauth
// @Accept json
// @Produce json
// @Param request body RegisterWebOAuthRequest true "OAuth registration request"
// @Success 201 {object} RegisterWebOAuthResponse "OAuth registered successfully"
// @Failure 400 {object} merrmid.ErrorResponse "Invalid request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /weboauth/register [post]
func (h *Handler) RegisterWebOAuth(c *gin.Context) {
	ctx := c.Request.Context()

	var req RegisterWebOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid request body", err)
		c.Error(err)
		return
	}

	clientAppPublicID, err := clientappval.ParsePublicID(req.ClientAppID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client_app_id format", err)
		c.Error(err)
		return
	}

	provider := shared.Provider(req.Provider)
	if !provider.IsValid() {
		err := merr.New(merr.ErrBadRequest, "invalid provider", nil)
		c.Error(err)
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

	response := toRegisterWebOAuthResponse()
	c.JSON(http.StatusCreated, response)
}
