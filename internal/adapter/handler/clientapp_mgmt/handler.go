package clientapp_mgmt

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

type Handler struct {
	clientAppMgmt in.ClientAppMgmtUsecase
}

func NewHandler(clientAppMgmt in.ClientAppMgmtUsecase) *Handler {
	return &Handler{
		clientAppMgmt: clientAppMgmt,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.CreateClientApp)
	rg.PUT("/:id", h.UpdateClientApp)
	rg.POST("/:id/refresh-secret", h.RefreshSecret)
	rg.GET("/", h.ListClientApps)
}
