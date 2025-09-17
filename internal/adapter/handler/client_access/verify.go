package client_access

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type VerifyClientResponse struct {
	IsValid     bool                 `json:"is_valid"`
	ServiceID   string               `json:"service_id,omitempty"`
	ServiceInfo *ServiceInfoResponse `json:"service_info,omitempty"`
	ClientApp   *ClientAppResponse   `json:"client_app,omitempty"`
}

type ServiceInfoResponse struct {
	ServiceID   string `json:"service_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type ClientAppResponse struct {
	ClientAppID string `json:"client_app_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

func toVerifyClientResponse(result *in.VerifyClientResponse) *VerifyClientResponse {
	if !result.IsValid {
		return &VerifyClientResponse{
			IsValid: false,
		}
	}

	return &VerifyClientResponse{
		IsValid:   true,
		ServiceID: result.ServiceID,
		ServiceInfo: &ServiceInfoResponse{
			ServiceID:   result.ServiceInfo.ServiceID.String(),
			Name:        result.ServiceInfo.Name,
			Description: result.ServiceInfo.Desc,
			IsActive:    result.ServiceInfo.IsActive,
		},
		ClientApp: &ClientAppResponse{
			ClientAppID: result.ClientAppID.String(),
			Name:        result.ClientAppInfo.Name,
			Description: result.ClientAppInfo.Desc,
			IsActive:    result.ClientAppInfo.IsActive,
		},
	}
}

// VerifyClient verifies client credentials using Basic Auth
// @Summary Verify client credentials
// @Description Verify client application credentials using Basic Authentication (client_id:client_secret)
// @Tags client-access
// @Accept json
// @Produce json
// @Security BasicAuth
// @Success 200 {object} VerifyClientResponse "Client verification result"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 401 {object} merrmid.ErrorResponse "Unauthorized - Invalid credentials"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /client-access/verify [post]
func (h *Handler) VerifyClient(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse Basic Auth credentials
	// Authorization: Basic base64(client_id:client_secret)
	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		err := merr.New(merr.ErrUnauthorized, "Basic Auth credentials required", nil)
		c.Error(err)
		return
	}

	if username == "" {
		err := merr.New(merr.ErrBadRequest, "client_id cannot be empty", nil)
		c.Error(err)
		return
	}

	if password == "" {
		err := merr.New(merr.ErrBadRequest, "client_secret cannot be empty", nil)
		c.Error(err)
		return
	}

	// Create usecase request
	usecaseReq := &in.VerifyClientRequest{
		ClientID:     username,         // client_id from Basic Auth username
		ClientSecret: []byte(password), // client_secret from Basic Auth password
	}

	// Call usecase
	result, err := h.clientAccess.VerifyClient(ctx, usecaseReq)
	if err != nil {
		c.Error(err)
		return
	}

	// Convert and return response
	response := toVerifyClientResponse(result)
	c.JSON(http.StatusOK, response)
}

