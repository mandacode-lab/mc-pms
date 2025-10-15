package servicemgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type UpdateServiceRequest struct {
	Name        *string `json:"name" binding:"required"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateServiceResponse struct {
	UpdatedAt time.Time `json:"updated_at"`
}

func toUpdateServiceResponse(result *in.UpdateServiceResult) *UpdateServiceResponse {
	return &UpdateServiceResponse{
		UpdatedAt: result.UpdatedAt,
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
	serviceID, err := service.NewPublicID(serviceIDStr)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service ID", err)
		_ = c.Error(err)
		return
	}

	var req UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid request body", err)
		_ = c.Error(err)
		return
	}

	var name *service.Name
	if req.Name != nil {
		n, err := service.NewName(*req.Name)
		if err != nil {
			err := merr.New(merr.ErrBadRequest, "invalid service name", err)
			_ = c.Error(err)
			return
		}
		name = &n
	}

	var desc *service.Description
	if req.Description != nil {
		d, err := service.NewDescription(*req.Description)
		if err != nil {
			err := merr.New(merr.ErrBadRequest, "invalid service description", err)
			_ = c.Error(err)
			return
		}
		desc = &d
	}

	input := &in.UpdateServiceInput{
		ServiceID:   serviceID,
		NewName:     name,
		NewDesc:     desc,
		NewIsActive: req.IsActive,
	}

	result, err := h.serviceMgmt.UpdateService(ctx, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := toUpdateServiceResponse(result)
	c.JSON(http.StatusOK, response)
}
