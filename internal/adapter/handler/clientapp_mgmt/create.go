package clientapp_mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

type CreateClientAppRequest struct {
	ServiceID   string `json:"service_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *Handler) CreateClientApp(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateClientAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	servicePublicID, err := serviceval.ParsePublicID(req.ServiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service_id format"})
		return
	}

	cmd := &in.CreateClientAppCommand{
		ServiceID: servicePublicID,
		Name:      req.Name,
		Desc:      req.Description,
	}

	view, err := h.clientAppMgmt.CreateClientApp(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, view)
}