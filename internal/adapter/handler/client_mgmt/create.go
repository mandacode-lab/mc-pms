package clientmgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
)

type CreateClientAppRequest struct {
	ServiceID   string `json:"service_id" binding:"required" example:"srv_1234567890abcdef"`
	Name        string `json:"name" binding:"required" example:"My App"`
	Description string `json:"description" example:"Description of my application"`
}

type CreateClientAppResponse struct {
	ServiceID   string    `json:"service_id" example:"srv_1234567890abcdef"`
	ClientAppID string    `json:"client_app_id" example:"app_abcdef1234567890"`
	Name        string    `json:"name" example:"My App"`
	Description string    `json:"description" example:"Description of my application"`
	IsActive    bool      `json:"is_active" example:"true"`
	Secret      string    `json:"secret" example:"secret_xyz789abc123def456"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func toCreateClientAppResponse(result *in.CreateClientAppResponse) *CreateClientAppResponse {
	return &CreateClientAppResponse{
		ServiceID:   result.ServiceID.String(),
		ClientAppID: result.ClientAppID.String(),
		Name:        result.Name,
		Description: result.Desc,
		IsActive:    result.IsActive,
		Secret:      string(result.Secret),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}
}

// CreateClientApp creates a new client application
// @Summary Create a new client app
// @Description Create a new OAuth client application
// @Tags client-apps
// @Accept json
// @Produce json
// @Param request body CreateClientAppRequest true "Client app creation request"
// @Success 201 {object} CreateClientAppResponse "Created client app with secret"
// @Failure 400 {object} merrmid.ErrorResponse "Invalid request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /client-apps [post]
func (h *Handler) CreateClientApp(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateClientAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid request body", err)
		_ = c.Error(err)
		return
	}

	// Parse service ID
	serviceID, err := serviceval.ParsePublicID(req.ServiceID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service ID", err)
		_ = c.Error(err)
		return
	}

	usecaseReq := &in.CreateClientAppRequest{
		ServiceID: serviceID,
		Name:      req.Name,
		Desc:      req.Description,
	}

	result, err := h.clientAppMgmt.CreateClientApp(ctx, usecaseReq)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := toCreateClientAppResponse(result)
	c.JSON(http.StatusCreated, response)
}
