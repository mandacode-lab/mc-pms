package clientmgmt

import (
	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
)

func toClientAppInfo(ca *entity.ClientApp, serviceID vo.ServicePublicID) in.ClientAppInfo {
	return in.ClientAppInfo{
		ServiceID:   serviceID,
		ClientAppID: ca.PublicID(),
		Name:        ca.Name(),
		Desc:        ca.Description().String(),
		IsActive:    ca.IsActive(),
		CreatedAt:   ca.CreatedAt(),
		UpdatedAt:   ca.UpdatedAt(),
	}
}

func toClientAppInfos(clientApps []*entity.ClientApp, serviceID vo.ServicePublicID) []in.ClientAppInfo {
	clientAppInfos := make([]in.ClientAppInfo, len(clientApps))
	for i, ca := range clientApps {
		clientAppInfos[i] = toClientAppInfo(ca, serviceID)
	}
	return clientAppInfos
}
