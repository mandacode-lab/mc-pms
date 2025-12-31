package admin_project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
)

type FindResponse struct {
	Projects []*ProjectResponse `json:"projects"`
}

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
