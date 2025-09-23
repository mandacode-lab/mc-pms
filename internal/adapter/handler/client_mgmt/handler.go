package clientmgmt

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/middleware"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type Handler struct {
	permission    *middleware.PermissionMiddleware
	clientAppMgmt in.ClientAppMgmtUsecase
}

func NewHandler(clientAppMgmt in.ClientAppMgmtUsecase) *Handler {
	return &Handler{
		clientAppMgmt: clientAppMgmt,
	}
}

func (h *Handler) SetPermissionMiddleware(permission *middleware.PermissionMiddleware) {
	h.permission = permission
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/",
		h.permission.RequirePermission(out.ResourceClientApp, out.ActionWrite),
		h.CreateClientApp)
	rg.PUT("/:id",
		h.permission.RequirePermission(out.ResourceClientApp, out.ActionWrite),
		h.UpdateClientApp)
	rg.POST("/:id/refresh-secret",
		h.permission.RequirePermission(out.ResourceClientApp, out.ActionWrite),
		h.RefreshSecret)
	rg.GET("/",
		h.permission.RequirePermission(out.ResourceClientApp, out.ActionRead),
		h.ListClientApps)
}
