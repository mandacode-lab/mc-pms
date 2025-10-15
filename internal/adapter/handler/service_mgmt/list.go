package servicemgmt

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type ListServicesResponse struct {
	Services []ServiceResponse `json:"services"`
	Total    int               `json:"total"`
}

// ListServices lists all services
// @Summary List services
// @Description List all services with optional filters
// @Tags services
// @Accept json
// @Produce json
// @Param name query string false "Filter by name (contains)"
// @Param is_active query boolean false "Filter by active status"
// @Param limit query int false "Limit number of results" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} ListServicesResponse "List of services"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /services [get]
func (h *Handler) ListServices(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters
	var nameContains *string
	if name := c.Query("name"); name != "" {
		nameContains = &name
	}

	var isActive *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		active, err := strconv.ParseBool(activeStr)
		if err != nil {
			err := merr.New(merr.ErrBadRequest, "invalid is_active parameter", err)
			_ = c.Error(err)
			return
		}
		isActive = &active
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit < 1 {
			err := merr.New(merr.ErrBadRequest, "invalid limit parameter", err)
			_ = c.Error(err)
			return
		}
		limit = parsedLimit
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil || parsedOffset < 0 {
			err := merr.New(merr.ErrBadRequest, "invalid offset parameter", err)
			_ = c.Error(err)
			return
		}
		offset = parsedOffset
	}

	input := &in.ListServicesInput{
		NameContains: nameContains,
		IsActive:     isActive,
		Limit:        limit,
		Offset:       offset,
	}

	result, err := h.serviceMgmt.ListServices(ctx, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	services := make([]ServiceResponse, 0, len(result.Services))
	for _, svc := range result.Services {
		services = append(services, toServiceResponse(&svc))
	}

	response := ListServicesResponse{
		Services: services,
		Total:    len(services),
	}

	c.JSON(http.StatusOK, response)
}
