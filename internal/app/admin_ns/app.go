package admin_ns

import (
	"github.com/mandacode-com/mandacode-pms/internal/port/rand"
	"github.com/mandacode-com/mandacode-pms/internal/port/repo"
)

type Application struct {
	nsRepo        repo.NamespaceRepo
	nsIDGenerator rand.IDGenerator
}

func NewApplication(
	nsRepo repo.NamespaceRepo,
	nsIDGenerator rand.IDGenerator,
) *Application {
	return &Application{
		nsRepo:        nsRepo,
		nsIDGenerator: nsIDGenerator,
	}
}
