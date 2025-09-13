package clientapp_mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

type UpdateClientAppRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *Handler) UpdateClientApp(c *gin.Context) {
	ctx := c.Request.Context()

	clientAppID := c.Param("id")
	clientAppPublicID, err := clientappval.ParsePublicID(clientAppID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client_app_id format"})
		return
	}

	var req UpdateClientAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, view)
}