package admin_ns

import (
	"context"

	"github.com/mandacode-com/mandacode-project/internal/domain/ns"
	"github.com/mandacode-com/mandacode-project/internal/port/app"
)

func (a *Application) Update(ctx context.Context, input app.UpdateNamespaceInput) (*app.AdminNamespaceResult, error) {
	namespace, err := a.nsRepo.GetByID(ctx, ns.NewNamespaceID(input.NamespaceID))
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		namespace.UpdateName(ns.NewNamespaceName(*input.Name))
	}

	if input.Description != nil {
		namespace.UpdateDescription(ns.NewNamespaceDescription(*input.Description))
	}

	if err := a.nsRepo.Upsert(ctx, namespace); err != nil {
		return nil, err
	}

	return &app.AdminNamespaceResult{
		NamespaceID: namespace.ID().String(),
		Name:        namespace.Name().String(),
		Description: namespace.Description().String(),
		CreatedAt:   namespace.CreatedAt(),
		UpdatedAt:   namespace.UpdatedAt(),
	}, nil
}
