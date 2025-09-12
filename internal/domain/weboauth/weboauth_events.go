package weboauth

import (
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
)

const (
	WebOAuthCreatedEventType = "weboauth.created"
	WebOAuthUpdatedEventType = "weboauth.updated"
	WebOAuthDeletedEventType = "weboauth.deleted"
)

type WebOAuthCreatedEvent struct {
	shared.BaseEvent
	Provider      string `json:"provider"`
	OAuthClientID string `json:"oauth_client_id"`
}

func NewWebOAuthCreatedEvent(aggregateID, provider, oauthClientID string) WebOAuthCreatedEvent {
	return WebOAuthCreatedEvent{
		BaseEvent:     shared.NewBaseEvent(aggregateID, WebOAuthCreatedEventType),
		Provider:      provider,
		OAuthClientID: oauthClientID,
	}
}

type WebOAuthUpdatedEvent struct {
	shared.BaseEvent
	UpdateType string `json:"update_type"`
}

func NewWebOAuthUpdatedEvent(aggregateID, updateType string) WebOAuthUpdatedEvent {
	return WebOAuthUpdatedEvent{
		BaseEvent:  shared.NewBaseEvent(aggregateID, WebOAuthUpdatedEventType),
		UpdateType: updateType,
	}
}

type WebOAuthDeletedEvent struct {
	shared.BaseEvent
}

func NewWebOAuthDeletedEvent(aggregateID string) WebOAuthDeletedEvent {
	return WebOAuthDeletedEvent{
		BaseEvent: shared.NewBaseEvent(aggregateID, WebOAuthDeletedEventType),
	}
}