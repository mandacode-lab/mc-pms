package admin_ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Delete(c *gin.Context) {
	namespaceID := c.Param("id")
	if namespaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace_id is required"})
		return
	}

	if err := h.adminNS.Delete(c.Request.Context(), namespaceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
