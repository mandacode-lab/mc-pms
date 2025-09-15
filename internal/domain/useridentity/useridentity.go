package useridentity

import (
	"time"

	"github.com/google/uuid"
	serviceval "github.com/mandacode-com/mandacode-service-hub/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/shared"
	useridentityval "github.com/mandacode-com/mandacode-service-hub/internal/domain/useridentity/value"
)

type UserIdentity struct {
	id         useridentityval.ID
	publicID   useridentityval.PublicID
	serviceID  serviceval.ID
	providerID string
	provider   shared.Provider
	createdAt  time.Time
	updatedAt  time.Time
	events     []shared.DomainEvent
}

func NewUserIdentity(
	id useridentityval.ID,
	publicID useridentityval.PublicID,
	serviceID serviceval.ID,
	providerID string,
	provider shared.Provider,
	createdAt time.Time,
	updatedAt time.Time,
) *UserIdentity {
	return &UserIdentity{
		id:         id,
		publicID:   publicID,
		serviceID:  serviceID,
		providerID: providerID,
		provider:   provider,
		createdAt:  createdAt,
		updatedAt:  updatedAt,
		events:     make([]shared.DomainEvent, 0),
	}
}

func DraftUserIdentity(
	serviceID serviceval.ID,
	providerID string,
	provider shared.Provider,
) *UserIdentity {
	now := time.Now().UTC()
	id := useridentityval.NewID(0) // ID will be set by the database
	publicID := useridentityval.NewPublicID(uuid.New())

	ui := NewUserIdentity(
		id,
		publicID,
		serviceID,
		providerID,
		provider,
		now,
		now,
	)

	ui.raise(NewUserIdentityCreatedEvent(ui.publicID.String(), string(provider), providerID))
	return ui
}

// Getter methods

func (ui *UserIdentity) ID() useridentityval.ID {
	return ui.id
}

func (ui *UserIdentity) PublicID() useridentityval.PublicID {
	return ui.publicID
}

func (ui *UserIdentity) ServiceID() serviceval.ID {
	return ui.serviceID
}

func (ui *UserIdentity) ProviderID() string {
	return ui.providerID
}

func (ui *UserIdentity) Provider() shared.Provider {
	return ui.provider
}

func (ui *UserIdentity) CreatedAt() time.Time {
	return ui.createdAt
}

func (ui *UserIdentity) UpdatedAt() time.Time {
	return ui.updatedAt
}

// Event methods

func (ui *UserIdentity) raise(e shared.DomainEvent) {
	ui.events = append(ui.events, e)
}

func (ui *UserIdentity) PullEvents() []shared.DomainEvent {
	events := ui.events
	ui.events = nil
	return events
}
