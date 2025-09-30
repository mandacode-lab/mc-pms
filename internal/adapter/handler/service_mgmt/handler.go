package servicemgmt

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/middleware"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type Handler struct {
	permission  *middleware.PermissionMiddleware
	serviceMgmt in.ServiceMgmtUsecase
}

func NewHandler(serviceMgmt in.ServiceMgmtUsecase) *Handler {
	return &Handler{
		serviceMgmt: serviceMgmt,
	}
}

func (h *Handler) SetPermissionMiddleware(permission *middleware.PermissionMiddleware) {
	h.permission = permission
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	read := h.permission.RequirePermission(out.ResourceService, out.ActionRead)
	write := h.permission.RequirePermission(out.ResourceService, out.ActionWrite)

	rg.GET("", read, h.ListServices)
	rg.GET("/", read, h.ListServices)
	rg.POST("", write, h.CreateService)
	rg.POST("/", write, h.CreateService)
	rg.PUT("/:id", write, h.UpdateService)
}
