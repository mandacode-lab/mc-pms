package clientaccess

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type VerifyClientResponse struct {
	ServiceID string `json:"service_id"`
}

func toVerifyClientResponse(result *in.VerifyClientResponse) *VerifyClientResponse {
	return &VerifyClientResponse{
		ServiceID: result.ServiceID.String(),
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
	username, password, hasAuth := c.Request.BasicAuth()
	if !hasAuth {
		err := merr.New(merr.ErrUnauthorized, "Basic Auth credentials required", nil)
		_ = c.Error(err)
		return
	}

	if username == "" {
		err := merr.New(merr.ErrBadRequest, "client_id cannot be empty", nil)
		_ = c.Error(err)
		return
	}

	if password == "" {
		err := merr.New(merr.ErrBadRequest, "client_secret cannot be empty", nil)
		_ = c.Error(err)
		return
	}

	// Parse client ID
	clientID, err := clientappval.ParsePublicID(username)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid client_id format", err)
		c.Error(err)
		return
	}

	// Create usecase request
	usecaseReq := &in.VerifyClientRequest{
		ClientID:     clientID,         // client_id from Basic Auth username
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

