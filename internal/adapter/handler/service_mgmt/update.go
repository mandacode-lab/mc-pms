package service_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	serviceval "github.com/mandacode-com/serengeti/internal/domain/service/value"
	"github.com/mandacode-com/serengeti/internal/port/in"
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

func toUpdateServiceResponse(view *in.UpdateServiceView) *UpdateServiceResponse {
	return &UpdateServiceResponse{
		ServiceID:   view.ServiceID.String(),
		Name:        view.Name,
		Description: view.Desc,
		IsActive:    view.IsActive,
		CreatedAt:   view.CreatedAt,
		UpdatedAt:   view.UpdatedAt,
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
// @Failure 400 {object} common.ErrorResponse "Bad request"
// @Failure 404 {object} common.ErrorResponse "Resource not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /services/{id} [put]
func (h *Handler) UpdateService(c *gin.Context) {
	ctx := c.Request.Context()

	serviceID := c.Param("id")
	servicePublicID, err := serviceval.ParsePublicID(serviceID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service_id format", nil)
		c.Error(err)
		return
	}

	var req UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid request body", err)
		c.Error(err)
		return
	}

	serviceName, err := serviceval.NewName(req.Name)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service name", err)
		c.Error(err)
		return
	}

	cmd := &in.UpdateServiceCommand{
		ServiceID: servicePublicID,
		NewName:   &serviceName,
		NewDesc:   &req.Description,
	}

	view, err := h.serviceMgmt.UpdateService(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	response := toUpdateServiceResponse(view)
	c.JSON(http.StatusOK, response)
}
