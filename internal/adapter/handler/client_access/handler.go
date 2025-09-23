package client_access

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
)

type Handler struct {
	clientAccess in.ClientAccessUsecase
}

func NewHandler(clientAccess in.ClientAccessUsecase) *Handler {
	return &Handler{
		clientAccess: clientAccess,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/verify", h.VerifyClient)
}
