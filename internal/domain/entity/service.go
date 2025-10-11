package entity

import (
	"time"

	"github.com/google/uuid"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
)

type Service struct {
	id          vo.ServiceID
	publicID    vo.ServicePublicID
	name        vo.ServiceName
	description vo.ServiceDescription
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
}

func NewService(
	id vo.ServiceID,
	publicID vo.ServicePublicID,
	name vo.ServiceName,
	description vo.ServiceDescription,
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

func DraftService(name vo.ServiceName, description vo.ServiceDescription) (*Service, error) {
	now := time.Now().UTC()

	// ID 0 is a special case for draft entities (will be set by DB)
	id, err := vo.NewServiceID(1) // Use 1 temporarily, will be replaced by DB
	if err != nil {
		return nil, err
	}

	publicID, err := vo.NewServicePublicID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	s := NewService(
		id,
		publicID,
		name,
		description,
		true,
		now,
		now,
	)

	return s, nil
}

// Getter methods

func (s *Service) ID() vo.ServiceID {
	return s.id
}

func (s *Service) PublicID() vo.ServicePublicID {
	return s.publicID
}

func (s *Service) Name() vo.ServiceName {
	return s.name
}

func (s *Service) Description() vo.ServiceDescription {
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

// Business methods

func (s *Service) UpdateName(name vo.ServiceName) {
	if s.name.String() != name.String() {
		s.name = name
		s.updatedAt = time.Now().UTC()
	}
}

func (s *Service) UpdateDescription(description vo.ServiceDescription) {
	s.description = description
	s.updatedAt = time.Now().UTC()
}

func (s *Service) Activate() {
	if !s.isActive {
		s.isActive = true
		s.updatedAt = time.Now().UTC()
	}
}

func (s *Service) Deactivate() {
	if s.isActive {
		s.isActive = false
		s.updatedAt = time.Now().UTC()
	}
}
