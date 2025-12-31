package admin_ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-project/internal/port/app"
)

type FindResponse struct {
	Namespaces []*NamespaceResponse `json:"namespaces"`
}

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
