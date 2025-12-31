package admin_project

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
	Project *ProjectResponse `json:"project"`
}

// Update godoc
// @Summary Update project
// @Description Update an existing project
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body UpdateRequest true "Update project request"
// @Success 200 {object} UpdateResponse "Project updated successfully"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 404 {object} merrmid.ErrorResponse "Project not found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /projects/{id} [post]
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
