package clientapp_mgmt

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) ListClientApps(ctx context.Context, req *in.ListClientAppsRequest) (*in.ListClientAppsResponse, error) {
	// Validate service exists
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, req.ServiceID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Find all client apps for the service
	clientApps, err := u.clientAppQueryRepo.FindByServiceID(ctx, service.ID())
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Convert to result models
	clientAppInfos := toClientAppInfos(clientApps, service.PublicID())

	return &in.ListClientAppsResponse{
		ServiceID:  service.PublicID(),
		ClientApps: clientAppInfos,
	}, nil
}
