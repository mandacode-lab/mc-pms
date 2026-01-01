package admin_ns

import (
	"context"

	"github.com/mandacode-lab/mc-pms/internal/domain/ns"
	"github.com/go-mandacode/merr"
)

func (a *Application) Delete(ctx context.Context, namespaceID string) error {
	if err := a.nsRepo.Delete(ctx, ns.NewNamespaceID(namespaceID)); err != nil {
		return merr.New(merr.ErrInternalServerError, "Failed to delete namespace", err)
	}
	return nil
}
