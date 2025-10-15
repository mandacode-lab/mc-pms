package clientmgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type RefreshSecretResponse struct {
	Secret string `json:"secret"`
}

func toRefreshSecretResponse(result *in.RefreshSecretResult) *RefreshSecretResponse {
	return &RefreshSecretResponse{
		Secret: string(result.RawSecret),
	}
}

// RefreshSecret generates a new secret for a client application
// @Summary Refresh client app secret
// @Description Generate a new secret for an OAuth client application
// @Tags client-apps
// @Produce json
// @Param id path string true "Client App ID"
// @Success 200 {object} RefreshSecretResponse "New secret"
// @Failure 400 {object} merrmid.ErrorResponse "Invalid request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /client-apps/{id}/refresh-secret [post]
func (h *Handler) RefreshSecret(c *gin.Context) {
	ctx := c.Request.Context()

	clientIDStr := c.Param("id")
	if clientIDStr == "" {
		err := merr.New(merr.ErrBadRequest, "client_app_id is required", nil)
		_ = c.Error(err)
		return
	}

	clientID, err := client.NewPublicID(clientIDStr)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client_app_id", err)
		_ = c.Error(err)
		return
	}

	input := &in.RefreshSecretInput{
		ClientID: clientID,
	}

	result, err := h.clientMgmt.RefreshSecret(ctx, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := toRefreshSecretResponse(result)
	c.JSON(http.StatusOK, response)
}
