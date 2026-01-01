package admin_ns

import (
	"context"

	"github.com/mandacode-lab/mc-pms/internal/domain/ns"
	"github.com/mandacode-lab/mc-pms/internal/port/app"
	"github.com/go-mandacode/merr"
)

func (a *Application) Update(ctx context.Context, input app.UpdateNamespaceInput) (*app.AdminNamespaceResult, error) {
	namespace, err := a.nsRepo.GetByID(ctx, ns.NewNamespaceID(input.NamespaceID))
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, "Namespace not found", err)
	}

	if input.Name != nil {
		namespace.UpdateName(ns.NewNamespaceName(*input.Name))
	}

	if input.Description != nil {
		namespace.UpdateDescription(ns.NewNamespaceDescription(*input.Description))
	}

	if err := a.nsRepo.Upsert(ctx, namespace); err != nil {
		return nil, merr.New(merr.ErrInternalServerError, "Failed to update namespace", err)
	}

	return &app.AdminNamespaceResult{
		NamespaceID: namespace.ID().String(),
		Name:        namespace.Name().String(),
		Description: namespace.Description().String(),
		CreatedAt:   namespace.CreatedAt(),
		UpdatedAt:   namespace.UpdatedAt(),
	}, nil
}
