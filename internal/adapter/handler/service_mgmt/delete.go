package servicemgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

// DeleteService deletes a service
// @Summary Delete a service
// @Description Delete a service by ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Success 204 "Service deleted successfully"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 404 {object} merrmid.ErrorResponse "Service not found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /services/{id} [delete]
func (h *Handler) DeleteService(c *gin.Context) {
	ctx := c.Request.Context()

	serviceIDStr := c.Param("id")
	if serviceIDStr == "" {
		err := merr.New(merr.ErrBadRequest, "service ID is required", nil)
		_ = c.Error(err)
		return
	}

	serviceID, err := service.NewPublicID(serviceIDStr)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service ID", err)
		_ = c.Error(err)
		return
	}

	input := &in.DeleteServiceInput{
		ServiceID: serviceID,
	}

	if err := h.serviceMgmt.DeleteService(ctx, input); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
