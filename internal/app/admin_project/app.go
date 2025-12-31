package admin_project

import (
	"github.com/mandacode-com/mandacode-project/internal/port/rand"
	"github.com/mandacode-com/mandacode-project/internal/port/repo"
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
