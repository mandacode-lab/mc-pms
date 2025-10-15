package clientmgmt

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/middleware"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type Handler struct {
	permission *middleware.PermissionMiddleware
	clientMgmt in.ClientMgmtUsecase
}

func NewHandler(clientAppMgmt in.ClientMgmtUsecase) *Handler {
	return &Handler{
		clientMgmt: clientAppMgmt,
	}
}

func (h *Handler) SetPermissionMiddleware(permission *middleware.PermissionMiddleware) {
	h.permission = permission
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	read := h.permission.RequirePermission(out.ResourceClientApp, out.ActionRead)
	write := h.permission.RequirePermission(out.ResourceClientApp, out.ActionWrite)

	rg.GET("", read, h.ListClientApps)
	rg.GET("/", read, h.ListClientApps)
	rg.POST("", write, h.CreateClientApp)
	rg.POST("/", write, h.CreateClientApp)
	rg.PUT("/:id", write, h.UpdateClientApp)
	rg.DELETE("/:id", write, h.DeleteClientApp)
	rg.POST("/:id/refresh-secret", write, h.RefreshSecret)
}
