package persistence

import (
	"context"
	"fmt"

	"github.com/mandacode-lab/mc-pms/ent"
	entNamespace "github.com/mandacode-lab/mc-pms/ent/namespace"
	entProject "github.com/mandacode-lab/mc-pms/ent/project"
	"github.com/mandacode-lab/mc-pms/internal/domain/ns"
	"github.com/mandacode-lab/mc-pms/internal/domain/project"
	"github.com/mandacode-lab/mc-pms/internal/port/repo"
)

type ProjectRepository struct {
	client *ent.Client
}

func NewProjectRepository(client *ent.Client) *ProjectRepository {
	return &ProjectRepository{
		client: client,
	}
}

// Upsert creates or updates a project
func (r *ProjectRepository) Upsert(ctx context.Context, proj *project.Project) error {
	// First, get or create the namespace relationship
	_, err := r.client.Namespace.Get(ctx, proj.NamespaceID().String())
	if err != nil {
		return fmt.Errorf("namespace not found: %w", err)
	}

	// Create or update the project with namespace edge
	err = r.client.Project.
		Create().
		SetID(proj.ID().String()).
		SetName(proj.Name().String()).
		SetDescription(proj.Description().String()).
		SetCreatedAt(proj.CreatedAt()).
		SetUpdatedAt(proj.UpdatedAt()).
		SetNamespaceID(proj.NamespaceID().String()).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)

	return err
}

// Delete removes a project by ID
func (r *ProjectRepository) Delete(ctx context.Context, id project.ProjectID) error {
	return r.client.Project.
		DeleteOneID(id.String()).
		Exec(ctx)
}

// GetByID retrieves a project by ID
func (r *ProjectRepository) GetByID(ctx context.Context, id project.ProjectID) (*project.Project, error) {
	entProj, err := r.client.Project.
		Query().
		Where(entProject.IDEQ(id.String())).
		WithNamespace().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return toProjectDomain(entProj), nil
}

// Find searches for projects matching the filter
func (r *ProjectRepository) Find(ctx context.Context, filter repo.ProjectFilter) ([]*project.Project, error) {
	query := r.client.Project.Query().WithNamespace()

	if filter.NameContains != "" {
		query = query.Where(entProject.NameContains(filter.NameContains))
	}

	if filter.NamespaceID.String() != "" {
		query = query.Where(entProject.HasNamespaceWith(entNamespace.IDEQ(filter.NamespaceID.String())))
	}

	entProjects, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find projects: %w", err)
	}

	projects := make([]*project.Project, len(entProjects))
	for i, entProj := range entProjects {
		projects[i] = toProjectDomain(entProj)
	}

	return projects, nil
}

// ListAll retrieves all projects
func (r *ProjectRepository) ListAll(ctx context.Context) ([]*project.Project, error) {
	entProjects, err := r.client.Project.Query().WithNamespace().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list all projects: %w", err)
	}

	projects := make([]*project.Project, len(entProjects))
	for i, entProj := range entProjects {
		projects[i] = toProjectDomain(entProj)
	}

	return projects, nil
}

// toProjectDomain converts ent.Project to domain.Project
func toProjectDomain(entProj *ent.Project) *project.Project {
	var namespaceID string
	if entProj.Edges.Namespace != nil {
		namespaceID = entProj.Edges.Namespace.ID
	}

	return project.NewProject(
		project.NewProjectID(entProj.ID),
		project.NewProjectName(entProj.Name),
		project.NewProjectDescription(entProj.Description),
		entProj.CreatedAt,
		entProj.UpdatedAt,
		ns.NewNamespaceID(namespaceID),
	)
}
