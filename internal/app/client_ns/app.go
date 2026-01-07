package client_ns

import (
	"github.com/mandacode-lab/mc-pms/internal/port/repo"
)

type Application struct {
	nsRepo repo.NamespaceRepo
}

func NewApplication(
	nsRepo repo.NamespaceRepo,
) *Application {
	return &Application{
		nsRepo: nsRepo,
	}
}
