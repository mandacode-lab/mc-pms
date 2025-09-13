package service_mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

type CreateServiceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *Handler) CreateService(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serviceName, err := serviceval.NewName(req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := &in.CreateServiceCommand{
		Name:        serviceName,
		Description: req.Description,
	}

	view, err := h.serviceMgmt.CreateService(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, view)
}