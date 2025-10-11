package clientaccess

import (
	"net/http"

	"github.com/gin-gonic/gin"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

// CheckServiceExists checks if a service exists
// @Summary Check if service exists
// @Description Check if a service exists by service ID. Returns 204 if exists, 404 if not found.
// @Tags client-access
// @Param service_id path string true "Service ID" example:"srv_1234567890abcdef"
// @Success 204 "Service exists"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request - Invalid service_id format"
// @Failure 404 {object} merrmid.ErrorResponse "Service not found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /client-access/service/{service_id}/exists [get]
func (h *Handler) CheckServiceExists(c *gin.Context) {
	ctx := c.Request.Context()

	// Get service_id from path parameter
	serviceIDStr := c.Param("service_id")
	if serviceIDStr == "" {
		err := merr.New(merr.ErrBadRequest, "service_id is required", nil)
		_ = c.Error(err)
		return
	}

	// Parse service ID
	serviceID, err := vo.NewServicePublicID(serviceIDStr)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service_id format", err)
		_ = c.Error(err)
		return
	}

	// Create usecase request
	usecaseReq := &in.CheckServiceExistsRequest{
		ServiceID: serviceID,
	}

	// Call usecase
	result, err := h.clientAccess.CheckServiceExists(ctx, usecaseReq)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !result {
		err := merr.New(merr.ErrNotFound, "service not found", nil)
		_ = c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
