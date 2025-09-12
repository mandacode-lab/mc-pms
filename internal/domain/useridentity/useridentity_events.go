package useridentity

import (
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
)

const (
	UserIdentityCreatedEventType = "useridentity.created"
	UserIdentityUpdatedEventType = "useridentity.updated"
	UserIdentityDeletedEventType = "useridentity.deleted"
)

type UserIdentityCreatedEvent struct {
	shared.BaseEvent
	Provider   string `json:"provider"`
	ProviderID string `json:"provider_id"`
}

func NewUserIdentityCreatedEvent(aggregateID, provider, providerID string) UserIdentityCreatedEvent {
	return UserIdentityCreatedEvent{
		BaseEvent:  shared.NewBaseEvent(aggregateID, UserIdentityCreatedEventType),
		Provider:   provider,
		ProviderID: providerID,
	}
}

type UserIdentityUpdatedEvent struct {
	shared.BaseEvent
	UpdateType string `json:"update_type"`
}

func NewUserIdentityUpdatedEvent(aggregateID, updateType string) UserIdentityUpdatedEvent {
	return UserIdentityUpdatedEvent{
		BaseEvent:  shared.NewBaseEvent(aggregateID, UserIdentityUpdatedEventType),
		UpdateType: updateType,
	}
}

type UserIdentityDeletedEvent struct {
	shared.BaseEvent
}

func NewUserIdentityDeletedEvent(aggregateID string) UserIdentityDeletedEvent {
	return UserIdentityDeletedEvent{
		BaseEvent: shared.NewBaseEvent(aggregateID, UserIdentityDeletedEventType),
	}
}
