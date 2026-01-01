package admin_ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-lab/mc-pms/internal/port/app"
)

type FindResponse struct {
	Namespaces []*NamespaceResponse `json:"namespaces"`
}

// Find godoc
// @Summary Find namespaces
// @Description Find namespaces with optional filters
// @Tags admin
// @Accept json
// @Produce json
// @Param name query string false "Filter by name contains"
// @Success 200 {object} FindResponse "Namespaces found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /namespaces [get]
func (h *Handler) Find(c *gin.Context) {
	nameContains := c.Query("name")

	results, err := h.adminNS.Find(c.Request.Context(), app.FindNamespaceInput{
		NameContains: nameContains,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	namespaces := make([]*NamespaceResponse, len(results))
	for i, ns := range results {
		namespaces[i] = toNamespaceResponse(ns)
	}

	c.JSON(http.StatusOK, FindResponse{
		Namespaces: namespaces,
	})
}
