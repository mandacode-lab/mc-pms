package project

import (
	"time"

	"github.com/mandacode-lab/mc-pms/internal/domain/ns"
)

type Project struct {
	id          ProjectID
	name        ProjectName
	description ProjectDescription
	createdAt   time.Time
	updatedAt   time.Time
	// Relations
	namespaceID ns.NamespaceID
}

func NewProject(
	id ProjectID,
	name ProjectName,
	description ProjectDescription,
	createdAt time.Time,
	updatedAt time.Time,
	namespaceID ns.NamespaceID,
) *Project {
	return &Project{
		id:          id,
		name:        name,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		namespaceID: namespaceID,
	}
}

func (p *Project) ID() ProjectID                   { return p.id }
func (p *Project) Name() ProjectName               { return p.name }
func (p *Project) Description() ProjectDescription { return p.description }
func (p *Project) CreatedAt() time.Time            { return p.createdAt }
func (p *Project) UpdatedAt() time.Time            { return p.updatedAt }
func (p *Project) NamespaceID() ns.NamespaceID     { return p.namespaceID }
