package repo

import (
	"context"

	"github.com/mandacode-com/mandacode-pms/internal/domain/ns"
)

type NamespaceFilter struct {
	NameContains string
}

type NamespaceRepo interface {
	// Command
	Upsert(ctx context.Context, ns *ns.Namespace) error
	Delete(ctx context.Context, id ns.NamespaceID) error
	// Query
	GetByID(ctx context.Context, id ns.NamespaceID) (*ns.Namespace, error)
	Find(ctx context.Context, filter NamespaceFilter) ([]*ns.Namespace, error)
	ListAll(ctx context.Context) ([]*ns.Namespace, error)
}
