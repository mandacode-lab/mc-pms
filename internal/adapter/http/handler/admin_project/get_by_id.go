package admin_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetByIDResponse struct {
	Project *ProjectResponse `json:"project"`
}

func (h *Handler) GetByID(c *gin.Context) {
	projectID := c.Param("id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	result, err := h.adminProject.GetByID(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetByIDResponse{
		Project: toProjectResponse(result),
	})
}
