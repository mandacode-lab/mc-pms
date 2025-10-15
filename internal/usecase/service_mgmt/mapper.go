package servicemgmt

import (
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
)

func toServiceInfo(s *service.Service) in.ServiceView {
	return in.ServiceView{
		ServiceID: s.PublicID(),
		Name:      s.Name(),
		Desc:      s.Description(),
		IsActive:  s.IsActive(),
		CreatedAt: s.CreatedAt(),
		UpdatedAt: s.UpdatedAt(),
	}
}
