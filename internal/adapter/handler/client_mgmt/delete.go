package clientmgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
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

	clientAppID := c.Param("id")
	if clientAppID == "" {
		err := merr.New(merr.ErrBadRequest, "client_app_id is required", nil)
		_ = c.Error(err)
		return
	}

	publicID, err := clientappval.NewPublicIDFromString(clientAppID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client_app_id", err)
		_ = c.Error(err)
		return
	}

	usecaseReq := &in.DeleteClientAppRequest{
		ClientAppID: publicID,
	}

	if err := h.clientAppMgmt.DeleteClientApp(ctx, usecaseReq); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}