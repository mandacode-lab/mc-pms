package project

import "time"

type Project struct {
	id          ProjectID
	name        ProjectName
	description ProjectDescription
	createdAt   time.Time
	updatedAt   time.Time
}

func NewProject(
	id ProjectID,
	name ProjectName,
	description ProjectDescription,
	createdAt time.Time,
	updatedAt time.Time,
) *Project {
	return &Project{
		id:          id,
		name:        name,
		description: description,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (p *Project) ID() ProjectID                   { return p.id }
func (p *Project) Name() ProjectName               { return p.name }
func (p *Project) Description() ProjectDescription { return p.description }
func (p *Project) CreatedAt() time.Time            { return p.createdAt }
func (p *Project) UpdatedAt() time.Time            { return p.updatedAt }
