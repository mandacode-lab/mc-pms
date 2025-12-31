package admin_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-project/internal/port/app"
)

type UpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type UpdateResponse struct {
	Project *ProjectResponse `json:"project"`
}

func (h *Handler) Update(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.adminProject.Update(c.Request.Context(), app.UpdateProjectInput{
		ProjectID:   projectID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UpdateResponse{
		Project: toProjectResponse(result),
	})
}
