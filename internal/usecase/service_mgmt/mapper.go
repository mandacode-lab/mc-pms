package servicemgmt

import (
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/mandacode-ssam/pkg/utils"
)

func toServiceInfo(s *service.Service) in.ServiceInfo {
	return in.ServiceInfo{
		ServiceID: s.PublicID(),
		Name:      s.Name().Value(),
		Desc:      utils.StringValue(s.Description()),
		IsActive:  s.IsActive(),
		CreatedAt: s.CreatedAt(),
		UpdatedAt: s.UpdatedAt(),
	}
}
