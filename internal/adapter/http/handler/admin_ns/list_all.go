package admin_ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListAllResponse struct {
	Namespaces []*NamespaceResponse `json:"namespaces"`
}

// ListAll godoc
// @Summary List all namespaces
// @Description List all namespaces without filters
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} ListAllResponse "Namespaces listed successfully"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /namespaces/all [get]
func (h *Handler) ListAll(c *gin.Context) {
	results, err := h.adminNS.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	namespaces := make([]*NamespaceResponse, len(results))
	for i, ns := range results {
		namespaces[i] = toNamespaceResponse(ns)
	}

	c.JSON(http.StatusOK, ListAllResponse{
		Namespaces: namespaces,
	})
}
