package client_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IsValidProjectNamespaceResponse struct {
	Valid bool `json:"valid"`
}

// IsValidProjectNamespace godoc
// @Summary Validate project-namespace relationship
// @Description Check if a project belongs to a specific namespace
// @Tags client
// @Accept json
// @Produce json
// @Param project_id query string true "Project ID"
// @Param namespace_id query string true "Namespace ID"
// @Success 200 {object} IsValidProjectNamespaceResponse "Validation result"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /projects/validate [get]
func (h *Handler) IsValidProjectNamespace(c *gin.Context) {
	projectID := c.Query("project_id")
	namespaceID := c.Query("namespace_id")

	if projectID == "" || namespaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id and namespace_id are required"})
		return
	}

	valid, err := h.clientProject.IsValidProjectNamespace(c.Request.Context(), projectID, namespaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, IsValidProjectNamespaceResponse{
		Valid: valid,
	})
}
