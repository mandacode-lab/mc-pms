package clientmgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type RefreshSecretResponse struct {
	Secret string `json:"secret"`
}

func toRefreshSecretResponse(result *in.RefreshSecretResponse) *RefreshSecretResponse {
	return &RefreshSecretResponse{
		Secret: string(result.Secret),
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

	clientAppIDStr := c.Param("id")

	// Parse client app ID
	clientAppID, err := vo.ParseClientAppPublicID(clientAppIDStr)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client app ID", err)
		_ = c.Error(err)
		return
	}

	usecaseReq := &in.RefreshSecretRequest{
		ClientAppID: clientAppID,
	}

	result, err := h.clientAppMgmt.RefreshSecret(ctx, usecaseReq)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := toRefreshSecretResponse(result)
	c.JSON(http.StatusOK, response)
}
