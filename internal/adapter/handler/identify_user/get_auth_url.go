package identify_user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
	clientappval "github.com/mandacode-com/mandacode-service-hub/internal/domain/clientapp/value"
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/shared"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
)

type GetAuthURLRequest struct {
	Provider        string `json:"provider" binding:"required" example:"google"`
	ClientAppID     string `json:"client_app_id" binding:"required" example:"app_abcdef1234567890"`
	ClientAppSecret string `json:"client_app_secret" binding:"required" example:"secret_xyz789abc123def456"`
}

type GetAuthURLResponse struct {
	AuthURL string `json:"auth_url" example:"https://accounts.google.com/o/oauth2/v2/auth?client_id=abc&redirect_uri=def&response_type=code&scope=openid+profile+email&state=xyz"`
	State   string `json:"state" example:"random_state_string_123"`
}

func toGetAuthURLResponse(view *in.GetAuthURLView) *GetAuthURLResponse {
	return &GetAuthURLResponse{
		AuthURL: view.AuthURL,
		State:   view.State,
	}
}

// GetAuthURL generates OAuth authorization URL
// @Summary Get OAuth authorization URL
// @Description Generate OAuth authorization URL for specified provider and client app
// @Tags auth
// @Produce json
// @Param provider query string true "OAuth Provider" Enums(google, kakao, naver)
// @Param client_app_id query string true "Client Application ID"
// @Param client_app_secret query string true "Client Application Secret"
// @Success 200 {object} GetAuthURLResponse "Authorization URL and state"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 404 {object} merrmid.ErrorResponse "Resource not found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /auth/auth-url [get]
func (h *Handler) GetAuthURL(c *gin.Context) {
	ctx := c.Request.Context()

	// Get parameters from query string
	providerStr := c.Query("provider")
	clientAppID := c.Query("client_app_id")
	clientAppSecret := c.Query("client_app_secret")

	// Validate required parameters
	if providerStr == "" || clientAppID == "" || clientAppSecret == "" {
		err := merr.New(merr.ErrBadRequest, "missing required parameters: provider, client_app_id, client_app_secret", nil)
		c.Error(err)
		return
	}

	provider := shared.Provider(providerStr)
	if !provider.IsValid() {
		err := merr.New(merr.ErrBadRequest, "invalid provider", nil)
		c.Error(err)
		return
	}

	clientAppPublicID, err := clientappval.ParsePublicID(clientAppID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client_app_id format", err)
		c.Error(err)
		return
	}

	cmd := &in.GetAuthURL{
		Provider:        provider,
		ClientAppID:     clientAppPublicID,
		ClientAppSecret: []byte(clientAppSecret),
	}

	view, err := h.identifyUser.GetAuthURL(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	response := toGetAuthURLResponse(view)
	c.JSON(http.StatusOK, response)
}
