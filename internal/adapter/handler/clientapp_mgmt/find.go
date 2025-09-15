package clientapp_mgmt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
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
		ServiceID:   info.ServiceID,
		ClientAppID: info.ClientAppID,
		Name:        info.Name,
		Description: info.Desc,
		IsActive:    info.IsActive,
		CreatedAt:   info.CreatedAt,
		UpdatedAt:   info.UpdatedAt,
	}
}

func toListClientAppsResponse(result *in.ListClientAppsResponse) *ListClientAppsResponse {
	clientApps := make([]ClientAppResponse, len(result.ClientApps))
	for i, info := range result.ClientApps {
		clientApps[i] = toClientAppResponse(info)
	}
	return &ListClientAppsResponse{
		ServiceID:  result.ServiceID,
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

	usecaseReq := &in.ListClientAppsRequest{
		ServiceID: serviceID,
	}

	result, err := h.clientAppMgmt.ListClientApps(ctx, usecaseReq)
	if err != nil {
		c.Error(err)
		return
	}

	response := toListClientAppsResponse(result)
	c.JSON(http.StatusOK, response)
}
