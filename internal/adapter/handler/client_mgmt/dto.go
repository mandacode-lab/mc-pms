package clientmgmt

import (
	"time"

	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
)

type ClientResponse struct {
	ServiceID   string    `json:"service_id"`
	ClientAppID string    `json:"client_app_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toClientResponse(view *in.ClientView) ClientResponse {
	return ClientResponse{
		ServiceID:   view.ServiceID.Value(),
		ClientAppID: view.ClientID.Value(),
		Name:        view.Name.Value(),
		Description: view.Desc.String(),
		IsActive:    view.IsActive,
		CreatedAt:   view.CreatedAt,
		UpdatedAt:   view.UpdatedAt,
	}
}
