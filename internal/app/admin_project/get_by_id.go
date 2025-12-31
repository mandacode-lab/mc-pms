package admin_project

import (
	"context"

	"github.com/mandacode-com/mandacode-project/internal/domain/project"
	"github.com/mandacode-com/mandacode-project/internal/port/app"
)

func (a *Application) GetByID(ctx context.Context, projectID string) (*app.AdminProjectResult, error) {
	proj, err := a.projectRepo.GetByID(ctx, project.NewProjectID(projectID))
	if err != nil {
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
