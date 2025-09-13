package weboauth_mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

func (h *Handler) ReadWebOAuth(c *gin.Context) {
	ctx := c.Request.Context()

	clientAppID := c.Query("client_app_id")
	if clientAppID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_app_id is required"})
		return
	}

	clientAppPublicID, err := clientappval.ParsePublicID(clientAppID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client_app_id format"})
		return
	}

	var provider *shared.Provider
	if providerStr := c.Query("provider"); providerStr != "" {
		p := shared.Provider(providerStr)
		if p.IsValid() {
			provider = &p
		}
	}

	query := &in.ReadWebOAuthQuery{
		ClientAppID: clientAppPublicID,
		Provider:    provider,
	}

	view, err := h.webOAuthMgmt.ReadWebOAuth(ctx, query)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, view)
}