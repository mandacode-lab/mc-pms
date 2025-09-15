package clientapp

import (
	"github.com/mandacode-com/mandacode-service-hub/internal/domain/shared"
)

const (
	ClientAppCreatedEventType = "clientapp.created"
	ClientAppUpdatedEventType = "clientapp.updated"
	ClientAppDeletedEventType = "clientapp.deleted"
)

type ClientAppCreatedEvent struct {
	shared.BaseEvent
	ClientName string `json:"client_name"`
}

func NewClientAppCreatedEvent(aggregateID, clientName string) ClientAppCreatedEvent {
	return ClientAppCreatedEvent{
		BaseEvent:  shared.NewBaseEvent(aggregateID, ClientAppCreatedEventType),
		ClientName: clientName,
	}
}

type ClientAppUpdatedEvent struct {
	shared.BaseEvent
	UpdateType string `json:"update_type"`
}

func NewClientAppUpdatedEvent(aggregateID, updateType string) ClientAppUpdatedEvent {
	return ClientAppUpdatedEvent{
		BaseEvent:  shared.NewBaseEvent(aggregateID, ClientAppUpdatedEventType),
		UpdateType: updateType,
	}
}

type ClientAppDeletedEvent struct {
	shared.BaseEvent
}

func NewClientAppDeletedEvent(aggregateID string) ClientAppDeletedEvent {
	return ClientAppDeletedEvent{
		BaseEvent: shared.NewBaseEvent(aggregateID, ClientAppDeletedEventType),
	}
}
