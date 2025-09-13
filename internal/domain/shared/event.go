package shared

import (
	"time"

	"github.com/google/uuid"
)

type DomainEvent interface {
	AggregateID() string
	EventType() string
	OccurredAt() time.Time
	EventID() string
}

type BaseEvent struct {
	eventID     string
	aggregateID string
	eventType   string
	occurredAt  time.Time
}

func NewBaseEvent(aggregateID, eventType string) BaseEvent {
	return BaseEvent{
		eventID:     uuid.New().String(),
		aggregateID: aggregateID,
		eventType:   eventType,
		occurredAt:  time.Now().UTC(),
	}
}

func NewEventID() string {
	return uuid.New().String()
}


func (e BaseEvent) EventID() string {
	return e.eventID
}

func (e BaseEvent) AggregateID() string {
	return e.aggregateID
}

func (e BaseEvent) EventType() string {
	return e.eventType
}

func (e BaseEvent) OccurredAt() time.Time {
	return e.occurredAt
}
