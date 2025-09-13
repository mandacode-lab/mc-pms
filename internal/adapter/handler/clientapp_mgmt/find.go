package clientapp_mgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

func (h *Handler) ListClientApps(c *gin.Context) {
	ctx := c.Request.Context()

	serviceID := c.Query("service_id")
	if serviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "service_id is required"})
		return
	}

	servicePublicID, err := serviceval.ParsePublicID(serviceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service_id format"})
		return
	}


	cmd := &in.ListClientAppsCommand{
		ServiceID: servicePublicID,
	}

	view, err := h.clientAppMgmt.ListClientApps(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, view)
}