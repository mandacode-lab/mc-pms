package user_mgmt

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
)

type Handler struct {
	userMgmt in.UserMgmtUsecase
}

func NewHandler(userMgmt in.UserMgmtUsecase) *Handler {
	return &Handler{
		userMgmt: userMgmt,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/", h.FindUserInfo)
}
