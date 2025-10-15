package service

import (
	"time"
)

// Service
type Service struct {
	id          ID
	publicID    PublicID
	name        Name
	description Description
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
}

func NewService(
	id ID,
	publicID PublicID,
	name Name,
	description Description,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *Service {
	return &Service{
		id:          id,
		publicID:    publicID,
		name:        name,
		description: description,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Getter methods

func (s *Service) ID() ID {
	return s.id
}

func (s *Service) PublicID() PublicID {
	return s.publicID
}

func (s *Service) Name() Name {
	return s.name
}

func (s *Service) Description() Description {
	return s.description
}

func (s *Service) IsActive() bool {
	return s.isActive
}

func (s *Service) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Service) UpdatedAt() time.Time {
	return s.updatedAt
}
