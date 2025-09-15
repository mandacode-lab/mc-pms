package client_access

import (
	"github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/pkg/utils"
)

func toClientAppInfo(ca *clientapp.ClientApp, serviceID serviceval.PublicID) in.ClientAppInfo {
	return in.ClientAppInfo{
		ServiceID:   serviceID.String(),
		ClientAppID: ca.PublicID().String(),
		Name:        ca.Name(),
		Desc:        utils.StringValue(ca.Description()),
		IsActive:    ca.IsActive(),
		CreatedAt:   ca.CreatedAt(),
		UpdatedAt:   ca.UpdatedAt(),
	}
}

func toServiceInfo(s *service.Service) in.ServiceInfo {
	return in.ServiceInfo{
		ServiceID: s.PublicID().String(),
		Name:      s.Name().Value(),
		Desc:      utils.StringValue(s.Description()),
		IsActive:  s.IsActive(),
		CreatedAt: s.CreatedAt(),
		UpdatedAt: s.UpdatedAt(),
	}
}
