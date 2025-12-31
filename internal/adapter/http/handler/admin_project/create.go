package admin_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
)

type CreateRequest struct {
	NamespaceID string `json:"namespace_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type CreateResponse struct {
	Project *ProjectResponse `json:"project"`
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.adminProject.Create(c.Request.Context(), app.CreateProjectInput{
		NamespaceID: req.NamespaceID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateResponse{
		Project: toProjectResponse(result),
	})
}
