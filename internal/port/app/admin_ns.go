package app

import (
	"context"
	"time"
)

type AdminNamespaceResult struct {
	NamespaceID string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateNamespaceInput struct {
	Name        string
	Description string
}

type UpdateNamespaceInput struct {
	NamespaceID string
	Name        *string
	Description *string
}

type FindNamespaceInput struct {
	NameContains string
}

type AdminNamespace interface {
	// Command
	Create(ctx context.Context, input CreateNamespaceInput) (*AdminNamespaceResult, error)
	Update(ctx context.Context, input UpdateNamespaceInput) (*AdminNamespaceResult, error)
	Delete(ctx context.Context, namespaceID string) error
	// Query
	GetByID(ctx context.Context, namespaceID string) (*AdminNamespaceResult, error)
	Find(ctx context.Context, input FindNamespaceInput) ([]*AdminNamespaceResult, error)
	ListAll(ctx context.Context) ([]*AdminNamespaceResult, error)
}
