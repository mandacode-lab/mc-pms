package clientapp_mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

func (h *Handler) RefreshSecret(c *gin.Context) {
	ctx := c.Request.Context()

	clientAppID := c.Param("id")
	clientAppPublicID, err := clientappval.ParsePublicID(clientAppID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client_app_id format"})
		return
	}

	cmd := &in.RefreshSecretCommand{
		ClientAppID: clientAppPublicID,
	}

	view, err := h.clientAppMgmt.RefreshSecret(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, view)
}