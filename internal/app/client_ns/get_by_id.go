package client_ns

import (
	"context"

	"github.com/mandacode-com/mandacode-pms/internal/domain/ns"
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
	"github.com/mandacode-com/merr"
)

func (a *Application) GetByID(ctx context.Context, namespaceID string) (*app.ClientNamespaceResult, error) {
	namespace, err := a.nsRepo.GetByID(ctx, ns.NewNamespaceID(namespaceID))
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, "Namespace not found", err)
	}

	return &app.ClientNamespaceResult{
		NamespaceID: namespace.ID().String(),
		Name:        namespace.Name().String(),
		Description: namespace.Description().String(),
		CreatedAt:   namespace.CreatedAt(),
		UpdatedAt:   namespace.UpdatedAt(),
	}, nil
}
