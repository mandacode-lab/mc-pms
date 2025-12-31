package admin_ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
)

type CreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateResponse struct {
	Namespace *NamespaceResponse `json:"namespace"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.adminNS.Create(c.Request.Context(), app.CreateNamespaceInput{
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateResponse{
		Namespace: toNamespaceResponse(result),
	})
}
