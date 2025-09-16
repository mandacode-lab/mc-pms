package service_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type CreateServiceRequest struct {
	Name        string `json:"name" binding:"required" example:"My Service"`
	Description string `json:"description" example:"Service description"`
}

type CreateServiceResponse struct {
	ServiceID   string    `json:"service_id" example:"srv_1234567890abcdef"`
	Name        string    `json:"name" example:"My Service"`
	Description string    `json:"description" example:"Service description"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func toCreateServiceResponse(result *in.CreateServiceResponse) *CreateServiceResponse {
	return &CreateServiceResponse{
		ServiceID:   result.ServiceID.String(),
		Name:        result.Name,
		Description: result.Desc,
		IsActive:    result.IsActive,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
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
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 404 {object} merrmid.ErrorResponse "Resource not found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /services [post]
func (h *Handler) CreateService(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid request body", err)
		c.Error(err)
		return
	}

	usecaseReq := &in.CreateServiceRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	result, err := h.serviceMgmt.CreateService(ctx, usecaseReq)
	if err != nil {
		c.Error(err)
		return
	}

	response := toCreateServiceResponse(result)
	c.JSON(http.StatusCreated, response)
}
