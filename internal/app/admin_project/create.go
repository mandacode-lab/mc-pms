package admin_project

import (
	"context"
	"time"

	"github.com/mandacode-com/mandacode-pms/internal/domain/ns"
	"github.com/mandacode-com/mandacode-pms/internal/domain/project"
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
	"github.com/mandacode-com/merr"
)

func (a *Application) Create(ctx context.Context, input app.CreateProjectInput) (*app.AdminProjectResult, error) {
	id, err := a.projectIDGenerator.Generate()
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, "Failed to generate project ID", err)
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
		return nil, merr.New(merr.ErrInternalServerError, "Failed to create project", err)
	}

	return &app.AdminProjectResult{
		ProjectID:   proj.ID().String(),
		Name:        proj.Name().String(),
		Description: proj.Description().String(),
		CreatedAt:   proj.CreatedAt(),
		UpdatedAt:   proj.UpdatedAt(),
	}, nil
}
