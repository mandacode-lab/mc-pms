package clientmgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type CreateClientAppRequest struct {
	ServiceID   string `json:"service_id" binding:"required" example:"srv_1234567890abcdef"`
	Name        string `json:"name" binding:"required" example:"My App"`
	Description string `json:"description" example:"Description of my application"`
}

type CreateClientAppResponse struct {
	ClientResponse
	Secret string `json:"secret" example:"secret_xyz789abc123def456"` // Raw secret, only returned on creation
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
	serviceID, err := service.NewPublicID(req.ServiceID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service ID", err)
		_ = c.Error(err)
		return
	}
	name, err := client.NewName(req.Name)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client name", err)
		_ = c.Error(err)
		return
	}
	var desc client.Description
	if req.Description != "" {
		desc, err = client.NewDescription(req.Description)
		if err != nil {
			err := merr.New(merr.ErrBadRequest, "invalid client description", err)
			_ = c.Error(err)
			return
		}
	}

	input := &in.CreateClientInput{
		ServiceID: serviceID,
		Name:      name,
		Desc:      desc,
	}

	result, err := h.clientMgmt.CreateClient(ctx, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := &CreateClientAppResponse{
		ClientResponse: toClientResponse(result.ClientView),
		Secret:         string(result.Secret),
	}
	c.JSON(http.StatusCreated, response)
}
