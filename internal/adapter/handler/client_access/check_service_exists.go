package clientaccess

import (
	"net/http"

	"github.com/gin-gonic/gin"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

// CheckServiceExists checks if a service exists
// @Summary Check if service exists
// @Description Check if a service exists by service ID
// @Tags client-access
// @Accept json
// @Produce json
// @Param service_id path string true "Service ID" example:"srv_1234567890abcdef"
// @Success 200 {object} CheckServiceExistsResponse "Service existence check result"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request - Invalid service_id format"
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
	serviceID := vo.NewServicePublicID(serviceIDStr)

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

	c.JSON(http.StatusOK, gin.H{"exists": true})
}
