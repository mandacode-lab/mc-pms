package weboauth_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
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
		ClientAppID:   info.ClientAppID,
		Provider:      info.Provider,
		OAuthClientID: info.OAuthClientID,
		RedirectURI:   info.RedirectURI,
		Scopes:        info.Scopes,
		CreatedAt:     info.CreatedAt,
		UpdatedAt:     info.UpdatedAt,
	}
}

func toReadWebOAuthResponse(result *in.ReadWebOAuthResponse) *ReadWebOAuthResponse {
	webOAuths := make([]WebOAuthResponse, len(result.WebOAuths))
	for i, info := range result.WebOAuths {
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

	var provider *string
	if providerStr := c.Query("provider"); providerStr != "" {
		provider = &providerStr
	}

	usecaseReq := &in.ReadWebOAuthRequest{
		ClientAppID: clientAppID,
		Provider:    provider,
	}

	result, err := h.webOAuthMgmt.ReadWebOAuth(ctx, usecaseReq)
	if err != nil {
		c.Error(err)
		return
	}

	response := toReadWebOAuthResponse(result)
	c.JSON(http.StatusOK, response)
}
