package identify_user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
	clientappval "github.com/mandacode-com/mandacode-service-hub/internal/domain/clientapp/value"
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/shared"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
)

type IdentifyByCodeRequest struct {
	Provider        string `json:"provider" binding:"required"`
	OAuthCode       string `json:"oauth_code" binding:"required"`
	State           string `json:"state" binding:"required"`
	ClientAppID     string `json:"client_app_id" binding:"required"`
	ClientAppSecret string `json:"client_app_secret" binding:"required"`
}

type IdentifyByCodeResponse struct {
	ServiceID string         `json:"service_id"`
	UserID    string         `json:"user_id"`
	Nickname  string         `json:"nickname"`
	Email     string         `json:"email"`
	Provider  string         `json:"provider"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	RawData   map[string]any `json:"raw_data"`
}

func toIdentifyByCodeResponse(result *in.UserIdentityResult) *IdentifyByCodeResponse {
	return &IdentifyByCodeResponse{
		ServiceID: result.ServiceID.String(),
		UserID:    result.UserID.String(),
		Nickname:  result.Nickname,
		Email:     result.Email,
		Provider:  string(result.Provider),
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		RawData:   result.RawData,
	}
}

// IdentifyByCode identifies user by OAuth authorization code
// @Summary Identify user by code
// @Description Identify user using OAuth authorization code flow
// @Tags auth
// @Produce json
// @Param provider query string true "OAuth Provider" Enums(google, kakao, naver)
// @Param oauth_code query string true "OAuth authorization code"
// @Param state query string true "OAuth state parameter"
// @Param client_app_id query string true "Client Application ID"
// @Param client_app_secret query string true "Client Application Secret"
// @Success 200 {object} IdentifyByCodeResponse "User identity information"
// @Failure 400 {object} merrmid.ErrorResponse "Invalid request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /auth/identify-by-code [get]
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
		err := merr.New(merr.ErrBadRequest, "missing required parameters: provider, oauth_code, state, client_app_id, client_app_secret", nil)
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

	cmd := &in.IdentifyByCode{
		Provider:        provider,
		OAuthCode:       oauthCode,
		State:           state,
		ClientAppID:     clientAppPublicID,
		ClientAppSecret: []byte(clientAppSecret),
	}

	result, err := h.identifyUser.IdentifyByCode(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	response := toIdentifyByCodeResponse(result)
	c.JSON(http.StatusOK, response)
}
