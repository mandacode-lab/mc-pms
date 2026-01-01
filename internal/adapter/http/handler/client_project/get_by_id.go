package client_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetByIDResponse struct {
	Project *ProjectResponse `json:"project"`
}

// GetByID godoc
// @Summary Get project by ID
// @Description Get a project by its ID (read-only)
// @Tags client
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} GetByIDResponse "Project found"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 404 {object} merrmid.ErrorResponse "Project not found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /projects/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	result, err := h.clientProject.GetByID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetByIDResponse{
		Project: toProjectResponse(result),
	})
}
