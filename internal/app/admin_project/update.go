package admin_project

import (
	"context"

	"github.com/mandacode-com/mandacode-project/internal/domain/project"
	"github.com/mandacode-com/mandacode-project/internal/port/app"
)

func (a *Application) Update(ctx context.Context, input app.UpdateProjectInput) (*app.AdminProjectResult, error) {
	proj, err := a.projectRepo.GetByID(ctx, project.NewProjectID(input.ProjectID))
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		proj.UpdateName(project.NewProjectName(*input.Name))
	}

	if input.Description != nil {
		proj.UpdateDescription(project.NewProjectDescription(*input.Description))
	}

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
