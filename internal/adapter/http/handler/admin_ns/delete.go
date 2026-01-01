package admin_ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary Delete namespace
// @Description Delete a namespace by ID
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Namespace ID"
// @Success 204 "Namespace deleted successfully"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /namespaces/{id} [delete]
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
