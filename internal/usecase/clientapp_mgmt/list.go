package clientapp_mgmt

import (
	"context"

	"github.com/mandacode-com/merr"
	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
)

func (u *Usecase) ListClientApps(ctx context.Context, cmd *in.ListClientAppsCommand) (*in.ListClientAppsResult, error) {
	// Validate service exists
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, cmd.ServiceID)
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, ErrServiceNotFoundMsg, err)
	}

	// Find all client apps for the service
	clientApps, err := u.clientAppQueryRepo.FindByServiceID(ctx, service.ID())
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, ErrInternalServerMsg, err)
	}

	// Convert to result models
	clientAppInfos := toClientAppInfos(clientApps, cmd.ServiceID)

	return &in.ListClientAppsResult{
		ServiceID:  cmd.ServiceID,
		ClientApps: clientAppInfos,
	}, nil
}
