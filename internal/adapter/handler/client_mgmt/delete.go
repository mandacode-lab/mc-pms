package clientmgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

// DeleteClientApp deletes a client application
// @Summary Delete a client application
// @Description Delete a client application by ID
// @Tags client-apps
// @Accept json
// @Produce json
// @Param id path string true "Client App ID"
// @Success 204 "Client app deleted successfully"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 404 {object} merrmid.ErrorResponse "Client app not found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /client-apps/{id} [delete]
func (h *Handler) DeleteClientApp(c *gin.Context) {
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

	input := &in.DeleteClientInput{
		ClientID: clientID,
	}

	if err := h.clientMgmt.DeleteClient(ctx, input); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
