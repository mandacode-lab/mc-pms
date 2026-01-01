package admin_ns

import (
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
)

type Handler struct {
	adminNS app.AdminNamespace
}

func NewHandler(adminNS app.AdminNamespace) *Handler {
	return &Handler{
		adminNS: adminNS,
	}
}

func toNamespaceResponse(ns *app.AdminNamespaceResult) *NamespaceResponse {
	return &NamespaceResponse{
		NamespaceID: ns.NamespaceID,
		Name:        ns.Name,
		Description: ns.Description,
		CreatedAt:   ns.CreatedAt,
		UpdatedAt:   ns.UpdatedAt,
	}
}
