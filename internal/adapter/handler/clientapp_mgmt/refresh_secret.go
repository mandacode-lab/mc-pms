package clientapp_mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

type RefreshSecretResponse struct {
	Secret string `json:"secret"`
}

func toRefreshSecretResponse(view *in.RefreshSecretView) *RefreshSecretResponse {
	return &RefreshSecretResponse{
		Secret: string(view.Secret),
	}
}

// RefreshSecret generates a new secret for a client application
// @Summary Refresh client app secret
// @Description Generate a new secret for an OAuth client application
// @Tags client-apps
// @Produce json
// @Param id path string true "Client App ID"
// @Success 200 {object} RefreshSecretResponse "New secret"
// @Failure 400 {object} common.ErrorResponse "Invalid request"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /client-apps/{id}/refresh-secret [post]
func (h *Handler) RefreshSecret(c *gin.Context) {
	ctx := c.Request.Context()

	clientAppID := c.Param("id")
	clientAppPublicID, err := clientappval.ParsePublicID(clientAppID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client_app_id format", err)
		c.Error(err)
		return
	}

	cmd := &in.RefreshSecretCommand{
		ClientAppID: clientAppPublicID,
	}

	view, err := h.clientAppMgmt.RefreshSecret(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	response := toRefreshSecretResponse(view)
	c.JSON(http.StatusOK, response)
}
