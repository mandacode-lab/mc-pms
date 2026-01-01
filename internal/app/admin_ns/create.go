package admin_ns

import (
	"context"
	"time"

	"github.com/mandacode-lab/mc-pms/internal/domain/ns"
	"github.com/mandacode-lab/mc-pms/internal/port/app"
	"github.com/go-mandacode/merr"
)

func (a *Application) Create(ctx context.Context, input app.CreateNamespaceInput) (*app.AdminNamespaceResult, error) {
	id, err := a.nsIDGenerator.Generate()
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, "Failed to generate namespace ID", err)
	}

	now := time.Now()
	namespace := ns.NewNamespace(
		ns.NewNamespaceID(id),
		ns.NewNamespaceName(input.Name),
		ns.NewNamespaceDescription(input.Description),
		now,
		now,
	)

	if err := a.nsRepo.Upsert(ctx, namespace); err != nil {
		return nil, merr.New(merr.ErrInternalServerError, "Failed to create namespace", err)
	}

	return &app.AdminNamespaceResult{
		NamespaceID: namespace.ID().String(),
		Name:        namespace.Name().String(),
		Description: namespace.Description().String(),
		CreatedAt:   namespace.CreatedAt(),
		UpdatedAt:   namespace.UpdatedAt(),
	}, nil
}
