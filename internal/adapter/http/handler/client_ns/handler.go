package client_ns

import (
	"github.com/mandacode-com/mandacode-project/internal/port/app"
)

type Handler struct {
	clientNS app.ClientNamespace
}

func NewHandler(clientNS app.ClientNamespace) *Handler {
	return &Handler{
		clientNS: clientNS,
	}
}

func toNamespaceResponse(ns *app.ClientNamespaceResult) *NamespaceResponse {
	return &NamespaceResponse{
		NamespaceID: ns.NamespaceID,
		Name:        ns.Name,
		Description: ns.Description,
		CreatedAt:   ns.CreatedAt,
		UpdatedAt:   ns.UpdatedAt,
	}
}
