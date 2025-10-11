package clientaccess

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/port/in"
	"github.com/mandacode-com/merr"
)

func (u *Usecase) CheckServiceExists(ctx context.Context, req *in.CheckServiceExistsRequest) (bool, error) {
	if req == nil {
		return false, merr.New(merr.ErrBadRequest, "invalid request", nil)
	}

	// Try to find service by public ID
	service, err := u.serviceQueryRepo.FindByPublicID(ctx, req.ServiceID)
	if err != nil {
		return false, merr.New(merr.ErrInternalServerError, "failed to check service existence", err)
	}

	if service == nil {
		return false, merr.New(merr.ErrNotFound, "service not found", nil)
	}

	return true, nil
}
