package weboauth_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
	clientappval "github.com/mandacode-com/serengeti/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti/internal/domain/shared"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

type WebOAuthResponse struct {
	ClientAppID   string    `json:"client_app_id"`
	Provider      string    `json:"provider"`
	OAuthClientID string    `json:"oauth_client_id"`
	RedirectURI   string    `json:"redirect_uri"`
	Scopes        []string  `json:"scopes"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ReadWebOAuthResponse struct {
	WebOAuths []WebOAuthResponse `json:"web_oauths"`
}

func toWebOAuthResponse(info in.WebOAuthInfo) WebOAuthResponse {
	return WebOAuthResponse{
		ClientAppID:   info.ClientAppID.String(),
		Provider:      string(info.Provider),
		OAuthClientID: info.OAuthClientID,
		RedirectURI:   info.RedirectURI,
		Scopes:        info.Scopes,
		CreatedAt:     info.CreatedAt,
		UpdatedAt:     info.UpdatedAt,
	}
}

func toReadWebOAuthResponse(view *in.ReadWebOAuthView) *ReadWebOAuthResponse {
	webOAuths := make([]WebOAuthResponse, len(view.WebOAuths))
	for i, info := range view.WebOAuths {
		webOAuths[i] = toWebOAuthResponse(info)
	}
	return &ReadWebOAuthResponse{
		WebOAuths: webOAuths,
	}
}

// ReadWebOAuth reads OAuth configurations for a client app
// @Summary Read OAuth configs
// @Description Read OAuth configurations for a client application
// @Tags weboauth
// @Produce json
// @Param client_app_id query string true "Client App ID"
// @Param provider query string false "OAuth Provider" Enums(google, kakao, naver)
// @Success 200 {object} ReadWebOAuthResponse "OAuth configurations"
// @Failure 400 {object} merrmid.ErrorResponse "Invalid request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /weboauth [get]
func (h *Handler) ReadWebOAuth(c *gin.Context) {
	ctx := c.Request.Context()

	clientAppID := c.Query("client_app_id")
	if clientAppID == "" {
		err := merr.New(merr.ErrBadRequest, "client_app_id is required", nil)
		c.Error(err)
		return
	}

	clientAppPublicID, err := clientappval.ParsePublicID(clientAppID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client_app_id format", err)
		c.Error(err)
		return
	}

	var provider *shared.Provider
	if providerStr := c.Query("provider"); providerStr != "" {
		p := shared.Provider(providerStr)
		if p.IsValid() {
			provider = &p
		}
	}

	query := &in.ReadWebOAuthQuery{
		ClientAppID: clientAppPublicID,
		Provider:    provider,
	}

	view, err := h.webOAuthMgmt.ReadWebOAuth(ctx, query)
	if err != nil {
		c.Error(err)
		return
	}

	response := toReadWebOAuthResponse(view)
	c.JSON(http.StatusOK, response)
}
