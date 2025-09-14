package clientapp_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	clientappval "github.com/mandacode-com/serengeti/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

type UpdateClientAppRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateClientAppResponse struct {
	ServiceID   string    `json:"service_id"`
	ClientAppID string    `json:"client_app_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toUpdateClientAppResponse(view *in.UpdateClientAppView) *UpdateClientAppResponse {
	return &UpdateClientAppResponse{
		ServiceID:   view.ServiceID.String(),
		ClientAppID: view.ClientAppID.String(),
		Name:        view.Name,
		Description: view.Desc,
		IsActive:    view.IsActive,
		CreatedAt:   view.CreatedAt,
		UpdatedAt:   view.UpdatedAt,
	}
}

// UpdateClientApp updates an existing client application
// @Summary Update a client app
// @Description Update an existing OAuth client application
// @Tags client-apps
// @Accept json
// @Produce json
// @Param id path string true "Client App ID"
// @Param request body UpdateClientAppRequest true "Client app update request"
// @Success 200 {object} UpdateClientAppResponse "Updated client app"
// @Failure 400 {object} common.ErrorResponse "Invalid request"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /client-apps/{id} [put]
func (h *Handler) UpdateClientApp(c *gin.Context) {
	ctx := c.Request.Context()

	clientAppID := c.Param("id")
	clientAppPublicID, err := clientappval.ParsePublicID(clientAppID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client_app_id format", err)
		c.Error(err)
		return
	}

	var req UpdateClientAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid request body", err)
		c.Error(err)
		return
	}

	cmd := &in.UpdateClientAppCommand{
		ClientAppID: clientAppPublicID,
		NewName:     &req.Name,
		NewDesc:     &req.Description,
	}

	view, err := h.clientAppMgmt.UpdateClientApp(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	response := toUpdateClientAppResponse(view)
	c.JSON(http.StatusOK, response)
}
