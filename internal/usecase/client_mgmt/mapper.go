package clientmgmt

import (
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
)

func toClientInfo(ca *client.Client) in.MgmtClientInfo {
	return in.MgmtClientInfo{
		ClientID:  ca.PublicID(),
		Name:      ca.Name(),
		Desc:      ca.Description(),
		IsActive:  ca.IsActive(),
		CreatedAt: ca.CreatedAt(),
		UpdatedAt: ca.UpdatedAt(),
	}
}

func toClientInfos(clientApps []*client.Client) []in.MgmtClientInfo {
	clientAppInfos := make([]in.MgmtClientInfo, len(clientApps))
	for i, ca := range clientApps {
		clientAppInfos[i] = toClientInfo(ca)
	}
	return clientAppInfos
}
