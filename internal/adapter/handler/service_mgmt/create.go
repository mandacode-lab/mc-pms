package service_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	serviceval "github.com/mandacode-com/serengeti/internal/domain/service/value"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

type CreateServiceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateServiceResponse struct {
	ServiceID   string    `json:"service_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toCreateServiceResponse(view *in.CreateServiceView) *CreateServiceResponse {
	return &CreateServiceResponse{
		ServiceID:   view.ServiceID.String(),
		Name:        view.Name,
		Description: view.Desc,
		IsActive:    view.IsActive,
		CreatedAt:   view.CreatedAt,
		UpdatedAt:   view.UpdatedAt,
	}
}

// CreateService creates a new service
// @Summary Create a new service
// @Description Create a new multi-tenant service
// @Tags services
// @Accept json
// @Produce json
// @Param request body CreateServiceRequest true "Service creation request"
// @Success 201 {object} CreateServiceResponse "Created service"
// @Failure 400 {object} common.ErrorResponse "Bad request"
// @Failure 404 {object} common.ErrorResponse "Resource not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /services [post]
func (h *Handler) CreateService(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateServiceRequest
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

	cmd := &in.CreateServiceCommand{
		Name:        serviceName,
		Description: req.Description,
	}

	view, err := h.serviceMgmt.CreateService(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	response := toCreateServiceResponse(view)
	c.JSON(http.StatusCreated, response)
}
