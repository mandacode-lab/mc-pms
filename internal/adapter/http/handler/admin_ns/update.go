package admin_ns

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
)

type UpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type UpdateResponse struct {
	Namespace *NamespaceResponse `json:"namespace"`
}

// Update godoc
// @Summary Update namespace
// @Description Update an existing namespace
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Namespace ID"
// @Param request body UpdateRequest true "Update namespace request"
// @Success 200 {object} UpdateResponse "Namespace updated successfully"
// @Failure 400 {object} merrmid.ErrorResponse "Bad request"
// @Failure 404 {object} merrmid.ErrorResponse "Namespace not found"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /namespaces/{id} [post]
func (h *Handler) Update(c *gin.Context) {
	namespaceID := c.Param("id")
	if namespaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "namespace_id is required"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.adminNS.Update(c.Request.Context(), app.UpdateNamespaceInput{
		NamespaceID: namespaceID,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, UpdateResponse{
		Namespace: toNamespaceResponse(result),
	})
}
