package client_project

import (
	"github.com/mandacode-lab/mc-pms/internal/port/repo"
)

type Application struct {
	projectRepo repo.ProjectRepo
}

func NewApplication(
	projectRepo repo.ProjectRepo,
) *Application {
	return &Application{
		projectRepo: projectRepo,
	}
}
