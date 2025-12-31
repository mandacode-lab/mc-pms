package admin_ns

import (
	"context"

	"github.com/mandacode-com/mandacode-project/internal/domain/ns"
)

func (a *Application) Delete(ctx context.Context, namespaceID string) error {
	return a.nsRepo.Delete(ctx, ns.NewNamespaceID(namespaceID))
}
