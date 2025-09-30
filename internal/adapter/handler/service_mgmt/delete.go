package servicemgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
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

	serviceID := c.Param("id")
	if serviceID == "" {
		err := merr.New(merr.ErrBadRequest, "service_id is required", nil)
		_ = c.Error(err)
		return
	}

	publicID, err := vo.ParseServicePublicID(serviceID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service_id", err)
		_ = c.Error(err)
		return
	}

	usecaseReq := &in.DeleteServiceRequest{
		ServiceID: publicID,
	}

	if err := h.serviceMgmt.DeleteService(ctx, usecaseReq); err != nil {
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
