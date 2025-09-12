package service

import (
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
)

const (
	ServiceCreatedEventType = "service.created"
	ServiceUpdatedEventType = "service.updated"
	ServiceDeletedEventType = "service.deleted"
)

type ServiceCreatedEvent struct {
	shared.BaseEvent
	ServiceName string `json:"service_name"`
}

func NewServiceCreatedEvent(aggregateID, serviceName string) ServiceCreatedEvent {
	return ServiceCreatedEvent{
		BaseEvent:   shared.NewBaseEvent(aggregateID, ServiceCreatedEventType),
		ServiceName: serviceName,
	}
}

type ServiceUpdatedEvent struct {
	shared.BaseEvent
	UpdateType string `json:"update_type"`
}

func NewServiceUpdatedEvent(aggregateID, updateType string) ServiceUpdatedEvent {
	return ServiceUpdatedEvent{
		BaseEvent:  shared.NewBaseEvent(aggregateID, ServiceUpdatedEventType),
		UpdateType: updateType,
	}
}

type ServiceDeletedEvent struct {
	shared.BaseEvent
}

func NewServiceDeletedEvent(aggregateID string) ServiceDeletedEvent {
	return ServiceDeletedEvent{
		BaseEvent: shared.NewBaseEvent(aggregateID, ServiceDeletedEventType),
	}
}