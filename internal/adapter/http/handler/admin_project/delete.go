package admin_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Delete(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	if err := h.adminProject.Delete(c.Request.Context(), projectID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
