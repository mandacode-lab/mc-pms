package weboauth_mgmt

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/serengeti-integrated/internal/port/in"
)

type Handler struct {
	webOAuthMgmt in.WebOAuthMgmtUsecase
}

func NewHandler(webOAuthMgmt in.WebOAuthMgmtUsecase) *Handler {
	return &Handler{
		webOAuthMgmt: webOAuthMgmt,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.RegisterWebOAuth)
	rg.GET("/", h.ReadWebOAuth)
}