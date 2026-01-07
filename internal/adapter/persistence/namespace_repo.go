package persistence

import (
	"context"
	"fmt"

	"github.com/mandacode-lab/mc-pms/ent"
	entNamespace "github.com/mandacode-lab/mc-pms/ent/namespace"
	"github.com/mandacode-lab/mc-pms/internal/domain/ns"
	"github.com/mandacode-lab/mc-pms/internal/port/repo"
)

type NamespaceRepository struct {
	client *ent.Client
}

func NewNamespaceRepository(client *ent.Client) *NamespaceRepository {
	return &NamespaceRepository{
		client: client,
	}
}

// Upsert creates or updates a namespace
func (r *NamespaceRepository) Upsert(ctx context.Context, namespace *ns.Namespace) error {
	return r.client.Namespace.
		Create().
		SetID(namespace.ID().String()).
		SetName(namespace.Name().String()).
		SetDescription(namespace.Description().String()).
		SetCreatedAt(namespace.CreatedAt()).
		SetUpdatedAt(namespace.UpdatedAt()).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
}

// Delete removes a namespace by ID
func (r *NamespaceRepository) Delete(ctx context.Context, id ns.NamespaceID) error {
	return r.client.Namespace.
		DeleteOneID(id.String()).
		Exec(ctx)
}

// GetByID retrieves a namespace by ID
func (r *NamespaceRepository) GetByID(ctx context.Context, id ns.NamespaceID) (*ns.Namespace, error) {
	entNs, err := r.client.Namespace.Get(ctx, id.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get namespace: %w", err)
	}

	return toNamespaceDomain(entNs), nil
}

// Find searches for namespaces matching the filter
func (r *NamespaceRepository) Find(ctx context.Context, filter repo.NamespaceFilter) ([]*ns.Namespace, error) {
	query := r.client.Namespace.Query()

	if filter.NameContains != "" {
		query = query.Where(entNamespace.NameContains(filter.NameContains))
	}

	entNamespaces, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find namespaces: %w", err)
	}

	namespaces := make([]*ns.Namespace, len(entNamespaces))
	for i, entNs := range entNamespaces {
		namespaces[i] = toNamespaceDomain(entNs)
	}

	return namespaces, nil
}

// ListAll retrieves all namespaces
func (r *NamespaceRepository) ListAll(ctx context.Context) ([]*ns.Namespace, error) {
	entNamespaces, err := r.client.Namespace.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all namespaces: %w", err)
	}

	namespaces := make([]*ns.Namespace, len(entNamespaces))
	for i, entNs := range entNamespaces {
		namespaces[i] = toNamespaceDomain(entNs)
	}

	return namespaces, nil
}

// toNamespaceDomain converts ent.Namespace to domain.Namespace
func toNamespaceDomain(entNs *ent.Namespace) *ns.Namespace {
	return ns.NewNamespace(
		ns.NewNamespaceID(entNs.ID),
		ns.NewNamespaceName(entNs.Name),
		ns.NewNamespaceDescription(entNs.Description),
		entNs.CreatedAt,
		entNs.UpdatedAt,
	)
}
