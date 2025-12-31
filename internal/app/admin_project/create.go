package admin_project

import (
	"context"
	"time"

	"github.com/mandacode-com/mandacode-project/internal/domain/ns"
	"github.com/mandacode-com/mandacode-project/internal/domain/project"
	"github.com/mandacode-com/mandacode-project/internal/port/app"
)

func (a *Application) Create(ctx context.Context, input app.CreateProjectInput) (*app.AdminProjectResult, error) {
	id, err := a.projectIDGenerator.Generate()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	proj := project.NewProject(
		project.NewProjectID(id),
		project.NewProjectName(input.Name),
		project.NewProjectDescription(input.Description),
		now,
		now,
		ns.NewNamespaceID(input.NamespaceID),
	)

	if err := a.projectRepo.Upsert(ctx, proj); err != nil {
		return nil, err
	}

	return &app.AdminProjectResult{
		ProjectID:   proj.ID().String(),
		Name:        proj.Name().String(),
		Description: proj.Description().String(),
		CreatedAt:   proj.CreatedAt(),
		UpdatedAt:   proj.UpdatedAt(),
	}, nil
}
