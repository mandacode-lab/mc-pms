package clientapp_mgmt

import (
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/clientapp"
	serviceval "github.com/mandacode-com/mandacode-service-hub/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-service-hub/internal/port/in"
	"github.com/mandacode-com/mandacode-service-hub/pkg/utils"
)

func toClientAppInfo(ca *clientapp.ClientApp, serviceID serviceval.PublicID) in.ClientAppInfo {
	return in.ClientAppInfo{
		ServiceID:   serviceID,
		ClientAppID: ca.PublicID(),
		Name:        ca.Name(),
		Desc:        utils.StringValue(ca.Description()),
		IsActive:    ca.IsActive(),
		CreatedAt:   ca.CreatedAt(),
		UpdatedAt:   ca.UpdatedAt(),
	}
}

func toClientAppInfos(clientApps []*clientapp.ClientApp, serviceID serviceval.PublicID) []in.ClientAppInfo {
	clientAppInfos := make([]in.ClientAppInfo, len(clientApps))
	for i, ca := range clientApps {
		clientAppInfos[i] = toClientAppInfo(ca, serviceID)
	}
	return clientAppInfos
}
