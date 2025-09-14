package service

import (
	"time"

	"github.com/google/uuid"
	serviceval "github.com/mandacode-com/serengeti/internal/domain/service/value"
	"github.com/mandacode-com/serengeti/internal/domain/shared"
)

type Service struct {
	id          serviceval.ID
	publicID    serviceval.PublicID
	name        serviceval.Name
	description *string
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
	events      []shared.DomainEvent
}

func NewService(
	id serviceval.ID,
	publicID serviceval.PublicID,
	name serviceval.Name,
	description *string,
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
		events:      make([]shared.DomainEvent, 0),
	}
}

func DraftService(name serviceval.Name, description *string) *Service {
	now := time.Now().UTC()
	id := serviceval.NewID(0) // ID will be set by the database
	publicID := serviceval.NewPublicID(uuid.New())

	s := NewService(
		id,
		publicID,
		name,
		description,
		true,
		now,
		now,
	)

	s.raise(NewServiceCreatedEvent(s.publicID.String(), s.name.Value()))
	return s
}

// Getter methods

func (s *Service) ID() serviceval.ID {
	return s.id
}

func (s *Service) PublicID() serviceval.PublicID {
	return s.publicID
}

func (s *Service) Name() serviceval.Name {
	return s.name
}

func (s *Service) Description() *string {
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

func (s *Service) UpdateName(name serviceval.Name) {
	if s.name.Value() != name.Value() {
		s.name = name
		s.updatedAt = time.Now().UTC()
		s.raise(NewServiceUpdatedEvent(s.publicID.String(), "name_changed"))
	}
}

func (s *Service) UpdateDescription(description *string) {
	s.description = description
	s.updatedAt = time.Now().UTC()
	s.raise(NewServiceUpdatedEvent(s.publicID.String(), "description_changed"))
}

func (s *Service) Activate() {
	if !s.isActive {
		s.isActive = true
		s.updatedAt = time.Now().UTC()
		s.raise(NewServiceUpdatedEvent(s.publicID.String(), "activated"))
	}
}

func (s *Service) Deactivate() {
	if s.isActive {
		s.isActive = false
		s.updatedAt = time.Now().UTC()
		s.raise(NewServiceUpdatedEvent(s.publicID.String(), "deactivated"))
	}
}

// Event methods
func (s *Service) raise(e shared.DomainEvent) {
	s.events = append(s.events, e)
}

func (s *Service) PullEvents() []shared.DomainEvent {
	events := s.events
	s.events = nil
	return events
}
