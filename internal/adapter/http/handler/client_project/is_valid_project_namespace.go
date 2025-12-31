package client_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IsValidProjectNamespaceResponse struct {
	Valid bool `json:"valid"`
}

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
