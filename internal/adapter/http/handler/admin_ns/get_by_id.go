package admin_ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetByIDResponse struct {
	Namespace *NamespaceResponse `json:"namespace"`
}

// GetByID godoc
// @Summary Get namespace by ID
// @Description Get a namespace by its ID
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Namespace ID"
// @Success 200 {object} GetByIDResponse "Namespace found"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 404 {object} merrmid.ErrorResponse "Namespace not found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /namespaces/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	namespaceID := c.Param("id")
	if namespaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace_id is required"})
		return
	}

	result, err := h.adminNS.GetByID(c.Request.Context(), namespaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GetByIDResponse{
		Namespace: toNamespaceResponse(result),
	})
}
