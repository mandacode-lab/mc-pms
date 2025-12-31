package admin_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Delete godoc
// @Summary Delete project
// @Description Delete a project by ID
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 204 "Project deleted successfully"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /projects/{id} [delete]
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
