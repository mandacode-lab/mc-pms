package servicemgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type CreateServiceRequest struct {
	Name        string `json:"name" binding:"required" example:"My Service"`
	Description string `json:"description" example:"Service description"`
}

type CreateServiceResponse struct {
	ServiceResponse
}

func toCreateServiceResponse(result *in.CreateServiceResult) *CreateServiceResponse {
	return &CreateServiceResponse{
		ServiceResponse: toServiceResponse(&result.ServiceView),
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
		_ = c.Error(err)
		return
	}

	name, err := service.NewName(req.Name)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service name", err)
		_ = c.Error(err)
		return
	}

	desc, err := service.NewDescription(req.Description)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service description", err)
		_ = c.Error(err)
		return
	}

	input := &in.CreateServiceInput{
		Name:        name,
		Description: desc,
	}

	result, err := h.serviceMgmt.CreateService(ctx, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := toCreateServiceResponse(result)
	c.JSON(http.StatusCreated, response)
}
