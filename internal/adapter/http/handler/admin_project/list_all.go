package admin_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListAllResponse struct {
	Projects []*ProjectResponse `json:"projects"`
}

// ListAll godoc
// @Summary List all projects
// @Description List all projects without filters
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} ListAllResponse "Projects listed successfully"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /projects/all [get]
func (h *Handler) ListAll(c *gin.Context) {
	results, err := h.adminProject.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	projects := make([]*ProjectResponse, len(results))
	for i, p := range results {
		projects[i] = toProjectResponse(p)
	}

	c.JSON(http.StatusOK, ListAllResponse{
		Projects: projects,
	})
}
