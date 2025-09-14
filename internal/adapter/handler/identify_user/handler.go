package identify_user

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/serengeti/internal/port/in"
)

type Handler struct {
	identifyUser in.IdentifyUserUsecase
}

func NewHandler(identifyUser in.IdentifyUserUsecase) *Handler {
	return &Handler{
		identifyUser: identifyUser,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/code", h.IdentifyByCode)
	rg.GET("/token", h.IdentifyByToken)
	rg.GET("/auth-url", h.GetAuthURL)
}
