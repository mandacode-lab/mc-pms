package admin_ns

import (
	"context"

	"github.com/mandacode-lab/mc-pms/internal/port/app"
	"github.com/go-mandacode/merr"
)

func (a *Application) ListAll(ctx context.Context) ([]*app.AdminNamespaceResult, error) {
	namespaces, err := a.nsRepo.ListAll(ctx)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, "Failed to list namespaces", err)
	}

	results := make([]*app.AdminNamespaceResult, len(namespaces))
	for i, namespace := range namespaces {
		results[i] = &app.AdminNamespaceResult{
			NamespaceID: namespace.ID().String(),
			Name:        namespace.Name().String(),
			Description: namespace.Description().String(),
			CreatedAt:   namespace.CreatedAt(),
			UpdatedAt:   namespace.UpdatedAt(),
		}
	}

	return results, nil
}
