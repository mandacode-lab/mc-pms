package service_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type UpdateServiceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateServiceResponse struct {
	ServiceID   string    `json:"service_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toUpdateServiceResponse(result *in.UpdateServiceResponse) *UpdateServiceResponse {
	return &UpdateServiceResponse{
		ServiceID:   result.ServiceID.String(),
		Name:        result.Name,
		Description: result.Desc,
		IsActive:    result.IsActive,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}
}

// UpdateService updates an existing service
// @Summary Update a service
// @Description Update an existing multi-tenant service
// @Tags services
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param request body UpdateServiceRequest true "Service update request"
// @Success 200 {object} UpdateServiceResponse "Updated service"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 404 {object} merrmid.ErrorResponse "Resource not found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /services/{id} [put]
func (h *Handler) UpdateService(c *gin.Context) {
	ctx := c.Request.Context()

	serviceIDStr := c.Param("id")

	// Parse service ID
	serviceID, err := serviceval.ParsePublicID(serviceIDStr)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service ID", err)
		c.Error(err)
		return
	}

	var req UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid request body", err)
		c.Error(err)
		return
	}

	usecaseReq := &in.UpdateServiceRequest{
		ServiceID: serviceID,
		NewName:   &req.Name,
		NewDesc:   &req.Description,
	}

	result, err := h.serviceMgmt.UpdateService(ctx, usecaseReq)
	if err != nil {
		c.Error(err)
		return
	}

	response := toUpdateServiceResponse(result)
	c.JSON(http.StatusOK, response)
}
