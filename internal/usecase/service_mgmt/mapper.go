package service_mgmt

import (
	"github.com/mandacode-com/serengeti/internal/domain/service"
	"github.com/mandacode-com/serengeti/internal/port/in"
	"github.com/mandacode-com/serengeti/pkg/utils"
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
