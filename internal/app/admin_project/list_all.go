package admin_project

import (
	"context"

	"github.com/mandacode-lab/mc-pms/internal/port/app"
	"github.com/go-mandacode/merr"
)

func (a *Application) ListAll(ctx context.Context) ([]*app.AdminProjectResult, error) {
	projects, err := a.projectRepo.ListAll(ctx)
	if err != nil {
		return nil, merr.New(merr.ErrInternalServerError, "Failed to list projects", err)
	}

	results := make([]*app.AdminProjectResult, len(projects))
	for i, proj := range projects {
		results[i] = &app.AdminProjectResult{
			ProjectID:   proj.ID().String(),
			Name:        proj.Name().String(),
			Description: proj.Description().String(),
			CreatedAt:   proj.CreatedAt(),
			UpdatedAt:   proj.UpdatedAt(),
		}
	}

	return results, nil
}
