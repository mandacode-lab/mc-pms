package admin_ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetByIDResponse struct {
	Namespace *NamespaceResponse `json:"namespace"`
}

func (h *Handler) GetByID(c *gin.Context) {
	namespaceID := c.Param("id")
	if namespaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace_id is required"})
		return
	}

	result, err := h.adminNS.GetByID(c.Request.Context(), namespaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetByIDResponse{
		Namespace: toNamespaceResponse(result),
	})
}
