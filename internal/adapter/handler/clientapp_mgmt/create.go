package clientapp_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

type CreateClientAppRequest struct {
	ServiceID   string `json:"service_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateClientAppResponse struct {
	ServiceID   string    `json:"service_id"`
	ClientAppID string    `json:"client_app_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	Secret      string    `json:"secret"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toCreateClientAppResponse(view *in.CreateClientAppView) *CreateClientAppResponse {
	return &CreateClientAppResponse{
		ServiceID:   view.ServiceID.String(),
		ClientAppID: view.ClientAppID.String(),
		Name:        view.Name,
		Description: view.Desc,
		IsActive:    view.IsActive,
		Secret:      string(view.Secret),
		CreatedAt:   view.CreatedAt,
		UpdatedAt:   view.UpdatedAt,
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
// @Failure 400 {object} common.ErrorResponse "Invalid request"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /client-apps [post]
func (h *Handler) CreateClientApp(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateClientAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid request body", err)
		c.Error(err)
		return
	}

	servicePublicID, err := serviceval.ParsePublicID(req.ServiceID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service_id format", err)
		c.Error(err)
		return
	}

	cmd := &in.CreateClientAppCommand{
		ServiceID: servicePublicID,
		Name:      req.Name,
		Desc:      req.Description,
	}

	view, err := h.clientAppMgmt.CreateClientApp(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	response := toCreateClientAppResponse(view)
	c.JSON(http.StatusCreated, response)
}
