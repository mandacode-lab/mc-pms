package clientmgmt

import (
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
)

func toClientView(ca *client.Client, serviceID service.PublicID) *in.ClientView {
	return &in.ClientView{
		ServiceID: serviceID,
		ClientID:  ca.PublicID(),
		Name:      ca.Name(),
		Desc:      ca.Description(),
		IsActive:  ca.IsActive(),
		CreatedAt: ca.CreatedAt(),
		UpdatedAt: ca.UpdatedAt(),
	}
}
