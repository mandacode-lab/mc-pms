package clientapp_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/serengeti/internal/adapter/handler/common"
	serviceval "github.com/mandacode-com/serengeti/internal/domain/service/value"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

type ClientAppResponse struct {
	ServiceID   string    `json:"service_id"`
	ClientAppID string    `json:"client_app_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListClientAppsResponse struct {
	ServiceID  string              `json:"service_id"`
	ClientApps []ClientAppResponse `json:"client_apps"`
}

func toClientAppResponse(info in.ClientAppInfo) ClientAppResponse {
	return ClientAppResponse{
		ServiceID:   info.ServiceID.String(),
		ClientAppID: info.ClientAppID.String(),
		Name:        info.Name,
		Description: info.Desc,
		IsActive:    info.IsActive,
		CreatedAt:   info.CreatedAt,
		UpdatedAt:   info.UpdatedAt,
	}
}

func toListClientAppsResponse(view *in.ListClientAppsView) *ListClientAppsResponse {
	clientApps := make([]ClientAppResponse, len(view.ClientApps))
	for i, info := range view.ClientApps {
		clientApps[i] = toClientAppResponse(info)
	}
	return &ListClientAppsResponse{
		ServiceID:  view.ServiceID.String(),
		ClientApps: clientApps,
	}
}

// ListClientApps lists all client applications for a service
// @Summary List client apps
// @Description List all OAuth client applications for a specific service
// @Tags client-apps
// @Produce json
// @Param service_id query string true "Service ID"
// @Success 200 {object} ListClientAppsResponse "List of client apps"
// @Failure 400 {object} merrmid.ErrorResponse "Invalid request"
// @Failure 500 {object} merrmid.ErrorResponse "Internal server error"
// @Router /client-apps [get]
func (h *Handler) ListClientApps(c *gin.Context) {
	ctx := c.Request.Context()

	serviceID := c.Query("service_id")
	if serviceID == "" {
		err := merr.New(merr.ErrBadRequest, "service_id is required", nil)
		c.Error(err)
		return
	}

	servicePublicID, err := serviceval.ParsePublicID(serviceID)
	if err != nil {
		err := merr.New(merr.ErrBadRequest, "invalid service_id format", err)
		c.Error(err)
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

	response := toListClientAppsResponse(view)
	c.JSON(http.StatusOK, response)
}
