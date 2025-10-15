package clientmgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type UpdateClientRequest struct {
	Name        *string `json:"name" binding:"required"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

type UpdateClientResponse struct {
	UpdatedAt time.Time `json:"updated_at"`
}

func toUpdateClientResponse(result *in.UpdateClientResult) *UpdateClientResponse {
	return &UpdateClientResponse{
		UpdatedAt: result.UpdatedAt,
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

	clientIDStr := c.Param("id")

	// Parse client app ID
	clientID, err := client.NewPublicID(clientIDStr)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client app ID", err)
		_ = c.Error(err)
		return
	}

	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid request body", err)
		_ = c.Error(err)
		return
	}

	var newName *client.Name
	if req.Name != nil {
		name, err := client.NewName(*req.Name)
		if err != nil {
			err := merr.New(merr.ErrBadRequest, "invalid client name", err)
			_ = c.Error(err)
			return
		}
		newName = &name
	}

	var newDesc *client.Description
	if req.Description != nil {
		if *req.Description != "" {
			desc, err := client.NewDescription(*req.Description)
			if err != nil {
				err := merr.New(merr.ErrBadRequest, "invalid client description", err)
				_ = c.Error(err)
				return
			}
			newDesc = &desc
		}
	}

	input := &in.UpdateClientInput{
		ClientID:    clientID,
		NewName:     newName,
		NewDesc:     newDesc,
		NewIsActive: req.IsActive,
	}

	result, err := h.clientMgmt.UpdateClient(ctx, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := toUpdateClientResponse(result)
	c.JSON(http.StatusOK, response)
}
