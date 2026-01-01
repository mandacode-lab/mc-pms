package admin_project

import (
	"github.com/mandacode-lab/mc-pms/internal/port/rand"
	"github.com/mandacode-lab/mc-pms/internal/port/repo"
)

type Application struct {
	projectRepo        repo.ProjectRepo
	projectIDGenerator rand.IDGenerator
}

func NewApplication(
	projectRepo repo.ProjectRepo,
	projectIDGenerator rand.IDGenerator,
) *Application {
	return &Application{
		projectRepo:        projectRepo,
		projectIDGenerator: projectIDGenerator,
	}
}
