package servicemgmt

import (
	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
)

func toServiceInfo(s *entity.Service) in.ServiceInfo {
	return in.ServiceInfo{
		ServiceID: s.PublicID(),
		Name:      s.Name().String(),
		Desc:      s.Description().String(),
		IsActive:  s.IsActive(),
		CreatedAt: s.CreatedAt(),
		UpdatedAt: s.UpdatedAt(),
	}
}
