package clientapp_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
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

func toUpdateClientAppResponse(result *in.UpdateClientAppResponse) *UpdateClientAppResponse {
	return &UpdateClientAppResponse{
		ServiceID:   result.ServiceID.String(),
		ClientAppID: result.ClientAppID.String(),
		Name:        result.Name,
		Description: result.Desc,
		IsActive:    result.IsActive,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
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
// @Failure 400 {object} merrmid.ErrorResponse "Invalid request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /client-apps/{id} [put]
func (h *Handler) UpdateClientApp(c *gin.Context) {
	ctx := c.Request.Context()

	clientAppIDStr := c.Param("id")

	// Parse client app ID
	clientAppID, err := clientappval.ParsePublicID(clientAppIDStr)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client app ID", err)
		c.Error(err)
		return
	}

	var req UpdateClientAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid request body", err)
		c.Error(err)
		return
	}

	usecaseReq := &in.UpdateClientAppRequest{
		ClientAppID: clientAppID,
		NewName:     &req.Name,
		NewDesc:     &req.Description,
	}

	result, err := h.clientAppMgmt.UpdateClientApp(ctx, usecaseReq)
	if err != nil {
		c.Error(err)
		return
	}

	response := toUpdateClientAppResponse(result)
	c.JSON(http.StatusOK, response)
}
