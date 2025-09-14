package identify_user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
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

type IdentifyByTokenResponse struct {
	ServiceID string         `json:"service_id"`
	UserID    string         `json:"user_id"`
	Nickname  string         `json:"nickname"`
	Email     string         `json:"email"`
	Provider  string         `json:"provider"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	RawData   map[string]any `json:"raw_data"`
}

func toIdentifyByTokenResponse(view *in.UserIdentityView) *IdentifyByTokenResponse {
	return &IdentifyByTokenResponse{
		ServiceID: view.ServiceID.String(),
		UserID:    view.UserID.String(),
		Nickname:  view.Nickname,
		Email:     view.Email,
		Provider:  string(view.Provider),
		CreatedAt: view.CreatedAt,
		UpdatedAt: view.UpdatedAt,
		RawData:   view.RawData,
	}
}

// IdentifyByToken identifies user by OAuth access token
// @Summary Identify user by token
// @Description Identify user using OAuth access token
// @Tags auth
// @Produce json
// @Param provider query string true "OAuth Provider" Enums(google, kakao, naver)
// @Param token query string true "OAuth access token"
// @Param client_app_id query string true "Client Application ID"
// @Param client_app_secret query string true "Client Application Secret"
// @Success 200 {object} IdentifyByTokenResponse "User identity information"
// @Failure 400 {object} common.ErrorResponse "Invalid request"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /auth/identify-by-token [get]
func (h *Handler) IdentifyByToken(c *gin.Context) {
	ctx := c.Request.Context()

	// Get parameters from query string
	providerStr := c.Query("provider")
	token := c.Query("token")
	clientAppID := c.Query("client_app_id")
	clientAppSecret := c.Query("client_app_secret")

	// Validate required parameters
	if providerStr == "" || token == "" || clientAppID == "" || clientAppSecret == "" {
		err := merr.New(merr.ErrBadRequest, "missing required parameters: provider, token, client_app_id, client_app_secret", nil)
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

	response := toIdentifyByTokenResponse(view)
	c.JSON(http.StatusOK, response)
}