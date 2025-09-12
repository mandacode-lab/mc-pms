package serviceclient

import (
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
)

const (
	ServiceClientCreatedEventType = "serviceclient.created"
	ServiceClientUpdatedEventType = "serviceclient.updated"
	ServiceClientDeletedEventType = "serviceclient.deleted"
)

type ServiceClientCreatedEvent struct {
	shared.BaseEvent
	ClientName string `json:"client_name"`
}

func NewServiceClientCreatedEvent(aggregateID, clientName string) ServiceClientCreatedEvent {
	return ServiceClientCreatedEvent{
		BaseEvent:  shared.NewBaseEvent(aggregateID, ServiceClientCreatedEventType),
		ClientName: clientName,
	}
}

type ServiceClientUpdatedEvent struct {
	shared.BaseEvent
	UpdateType string `json:"update_type"`
}

func NewServiceClientUpdatedEvent(aggregateID, updateType string) ServiceClientUpdatedEvent {
	return ServiceClientUpdatedEvent{
		BaseEvent:  shared.NewBaseEvent(aggregateID, ServiceClientUpdatedEventType),
		UpdateType: updateType,
	}
}

type ServiceClientDeletedEvent struct {
	shared.BaseEvent
}

func NewServiceClientDeletedEvent(aggregateID string) ServiceClientDeletedEvent {
	return ServiceClientDeletedEvent{
		BaseEvent: shared.NewBaseEvent(aggregateID, ServiceClientDeletedEventType),
	}
}