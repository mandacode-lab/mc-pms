package admin_project

import (
	"context"

	"github.com/mandacode-com/mandacode-pms/internal/domain/ns"
	"github.com/mandacode-com/mandacode-pms/internal/port/app"
	"github.com/mandacode-com/mandacode-pms/internal/port/repo"
)

func (a *Application) Find(ctx context.Context, input app.FindProjectInput) ([]*app.AdminProjectResult, error) {
	projects, err := a.projectRepo.Find(ctx, repo.ProjectFilter{
		NameContains: input.NameContains,
		NamespaceID:  ns.NewNamespaceID(input.NamespaceID),
	})
	if err != nil {
		return nil, err
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
