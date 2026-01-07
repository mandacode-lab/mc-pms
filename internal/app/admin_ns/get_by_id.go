package admin_ns

import (
	"context"

	"github.com/mandacode-lab/mc-pms/internal/domain/ns"
	"github.com/mandacode-lab/mc-pms/internal/port/app"
	"github.com/go-mandacode/merr"
)

func (a *Application) GetByID(ctx context.Context, namespaceID string) (*app.AdminNamespaceResult, error) {
	namespace, err := a.nsRepo.GetByID(ctx, ns.NewNamespaceID(namespaceID))
	if err != nil {
		return nil, merr.New(merr.ErrNotFound, "Namespace not found", err)
	}

	return &app.AdminNamespaceResult{
		NamespaceID: namespace.ID().String(),
		Name:        namespace.Name().String(),
		Description: namespace.Description().String(),
		CreatedAt:   namespace.CreatedAt(),
		UpdatedAt:   namespace.UpdatedAt(),
	}, nil
}
