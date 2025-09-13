package service_mgmt

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

type Handler struct {
	serviceMgmt in.ServiceMgmtUsecase
}

func NewHandler(serviceMgmt in.ServiceMgmtUsecase) *Handler {
	return &Handler{
		serviceMgmt: serviceMgmt,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.CreateService)
	rg.PUT("/:id", h.UpdateService)
}