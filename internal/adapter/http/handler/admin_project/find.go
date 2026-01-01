package admin_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
)

type FindResponse struct {
	Projects []*ProjectResponse `json:"projects"`
}

// Find godoc
// @Summary Find projects
// @Description Find projects with optional filters
// @Tags admin
// @Accept json
// @Produce json
// @Param name query string false "Filter by name contains"
// @Param namespace_id query string false "Filter by namespace ID"
// @Success 200 {object} FindResponse "Projects found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /projects [get]
func (h *Handler) Find(c *gin.Context) {
	nameContains := c.Query("name")
	namespaceID := c.Query("namespace_id")

	results, err := h.adminProject.Find(c.Request.Context(), app.FindProjectInput{
		NameContains: nameContains,
		NamespaceID:  namespaceID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	projects := make([]*ProjectResponse, len(results))
	for i, p := range results {
		projects[i] = toProjectResponse(p)
	}

	c.JSON(http.StatusOK, FindResponse{
		Projects: projects,
	})
}
