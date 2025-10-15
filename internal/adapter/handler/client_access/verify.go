package clientaccess

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type VerifyClientResponse struct {
	ServiceID string `json:"service_id"`
}

func toVerifyClientResponse(result *in.VerifyClientResult) *VerifyClientResponse {
	return &VerifyClientResponse{
		ServiceID: result.ServiceID.Value(),
	}
}

// VerifyClient verifies client credentials using Basic Auth
// @Summary Verify client credentials
// @Description Verify client application credentials using Basic Authentication (client_id:client_secret)
// @Tags client-access
// @Accept json
// @Produce json
// @Security BasicAuth
// @Success 200 {object} VerifyClientResponse "Client verification successful - returns service_id"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request - Invalid client_id format"
// @Failure 401 {object} merrmid.ErrorResponse "Unauthorized - Invalid client_secret"
// @Failure 403 {object} merrmid.ErrorResponse "Forbidden - Client or service inactive"
// @Failure 404 {object} merrmid.ErrorResponse "Not found - Client or service not found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /client-access/verify [post]
func (h *Handler) VerifyClient(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse Basic Auth credentials
	// Authorization: Basic base64(client_id:client_secret)
	clientIDStr, clientSecretStr, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		err := merr.New(merr.ErrUnauthorized, "Basic Auth credentials required", nil)
		_ = c.Error(err)
		return
	}

	if clientIDStr == "" {
		err := merr.New(merr.ErrBadRequest, "client_id cannot be empty", nil)
		_ = c.Error(err)
		return
	}

	if clientSecretStr == "" {
		err := merr.New(merr.ErrBadRequest, "client_secret cannot be empty", nil)
		_ = c.Error(err)
		return
	}

	// Parse client ID
	clientID, err := client.NewPublicID(clientIDStr)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client_id format", err)
		_ = c.Error(err)
		return
	}

	// Create usecase request
	input := &in.VerifyClientInput{
		ClientID:     clientID,                // client_id from Basic Auth username
		ClientSecret: []byte(clientSecretStr), // client_secret from Basic Auth password (base64 encoded)
	}

	// Call usecase
	result, err := h.clientAccess.VerifyClient(ctx, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// Convert and return response
	response := toVerifyClientResponse(result)
	c.JSON(http.StatusOK, response)
}
