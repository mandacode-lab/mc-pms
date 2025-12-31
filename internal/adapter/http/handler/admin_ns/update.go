package admin_ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
)

type UpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type UpdateResponse struct {
	Namespace *NamespaceResponse `json:"namespace"`
}

func (h *Handler) Update(c *gin.Context) {
	namespaceID := c.Param("id")
	if namespaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace_id is required"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.adminNS.Update(c.Request.Context(), app.UpdateNamespaceInput{
		NamespaceID: namespaceID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UpdateResponse{
		Namespace: toNamespaceResponse(result),
	})
}
