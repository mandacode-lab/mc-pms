package clientmgmt

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
	_ "github.com/mandacode-com/merr/middleware"
)

type ListClientsResponse struct {
	Clients []ClientResponse `json:"client_apps"`
}

func toListClientAppsResponse(result *in.ListClientsResult) *ListClientsResponse {
	clients := make([]ClientResponse, len(result.Clients))
	for i, c := range result.Clients {
		clients[i] = toClientResponse(c)
	}
	return &ListClientsResponse{
		Clients: clients,
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

	// Parse and validate service_id query parameter
	var serviceID *service.PublicID
	serviceIDStr := c.Query("service_id")
	if serviceIDStr != "" {
		sid, err := service.NewPublicID(serviceIDStr)
		if err != nil {
			err := merr.New(merr.ErrBadRequest, "invalid service ID", err)
			_ = c.Error(err)
			return
		}
		serviceID = &sid
	}

	input := &in.ListClientsInput{
		ServiceID: serviceID,
	}

	result, err := h.clientMgmt.ListClients(ctx, input)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response := toListClientAppsResponse(result)
	c.JSON(http.StatusOK, response)
}
