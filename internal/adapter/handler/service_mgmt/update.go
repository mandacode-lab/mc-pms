package service_mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

type UpdateServiceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *Handler) UpdateService(c *gin.Context) {
	ctx := c.Request.Context()

	serviceID := c.Param("id")
	servicePublicID, err := serviceval.ParsePublicID(serviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service_id format"})
		return
	}

	var req UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serviceName, err := serviceval.NewName(req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := &in.UpdateServiceCommand{
		ServiceID:   servicePublicID,
		NewName:     &serviceName,
		NewDesc:     &req.Description,
	}

	view, err := h.serviceMgmt.UpdateService(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, view)
}