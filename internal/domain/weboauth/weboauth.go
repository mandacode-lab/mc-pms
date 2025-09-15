package weboauth

import (
	"time"

	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/shared"
	weboauthval "github.com/mandacode-com/mandacode-ssam/internal/domain/weboauth/value"
)

type WebOAuth struct {
	id               weboauthval.ID
	clientAppID      clientappval.ID
	provider         shared.Provider
	oauthClientID    string
	oauthSecretCT    []byte
	oauthSecretNonce []byte
	dekWrapped       []byte
	dekNonce         []byte
	dekRotatedAt     time.Time
	redirectURI      string
	scopes           []string
	createdAt        time.Time
	updatedAt        time.Time
	events           []shared.DomainEvent
}

func NewWebOAuth(
	id weboauthval.ID,
	clientAppID clientappval.ID,
	provider shared.Provider,
	oauthClientID string,
	oauthSecretCT []byte,
	oauthSecretNonce []byte,
	dekWrapped []byte,
	dekNonce []byte,
	dekRotatedAt time.Time,
	redirectURI string,
	scopes []string,
	createdAt time.Time,
	updatedAt time.Time,
) *WebOAuth {
	return &WebOAuth{
		id:               id,
		clientAppID:      clientAppID,
		provider:         provider,
		oauthClientID:    oauthClientID,
		oauthSecretCT:    oauthSecretCT,
		oauthSecretNonce: oauthSecretNonce,
		dekWrapped:       dekWrapped,
		dekNonce:         dekNonce,
		dekRotatedAt:     dekRotatedAt,
		redirectURI:      redirectURI,
		scopes:           scopes,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
		events:           make([]shared.DomainEvent, 0),
	}
}

func DraftWebOAuth(
	clientAppID clientappval.ID,
	provider shared.Provider,
	oauthClientID string,
	oauthSecretCT []byte,
	oauthSecretNonce []byte,
	dekWrapped []byte,
	dekNonce []byte,
	redirectURI string,
	scopes []string,
) *WebOAuth {
	now := time.Now().UTC()
	id := weboauthval.NewID(0) // ID will be set by the database

	wo := NewWebOAuth(
		id,
		clientAppID,
		provider,
		oauthClientID,
		oauthSecretCT,
		oauthSecretNonce,
		dekWrapped,
		dekNonce,
		now,
		redirectURI,
		scopes,
		now,
		now,
	)

	wo.raise(NewWebOAuthCreatedEvent(wo.id.String(), string(provider), oauthClientID))
	return wo
}

// Getter methods

func (wo *WebOAuth) ID() weboauthval.ID {
	return wo.id
}

func (wo *WebOAuth) ClientAppID() clientappval.ID {
	return wo.clientAppID
}

func (wo *WebOAuth) Provider() shared.Provider {
	return wo.provider
}

func (wo *WebOAuth) OAuthClientID() string {
	return wo.oauthClientID
}

func (wo *WebOAuth) OAuthSecretCT() []byte {
	return wo.oauthSecretCT
}

func (wo *WebOAuth) OAuthSecretNonce() []byte {
	return wo.oauthSecretNonce
}

func (wo *WebOAuth) DEKWrapped() []byte {
	return wo.dekWrapped
}

func (wo *WebOAuth) DEKNonce() []byte {
	return wo.dekNonce
}

func (wo *WebOAuth) DEKRotatedAt() time.Time {
	return wo.dekRotatedAt
}

func (wo *WebOAuth) RedirectURI() string {
	return wo.redirectURI
}

func (wo *WebOAuth) Scopes() []string {
	return wo.scopes
}

func (wo *WebOAuth) CreatedAt() time.Time {
	return wo.createdAt
}

func (wo *WebOAuth) UpdatedAt() time.Time {
	return wo.updatedAt
}

// Business methods

func (wo *WebOAuth) UpdateClientID(clientID string) {
	if wo.oauthClientID != clientID {
		wo.oauthClientID = clientID
		wo.updatedAt = time.Now().UTC()
		wo.raise(NewWebOAuthUpdatedEvent(wo.id.String(), "client_id_changed"))
	}
}

func (wo *WebOAuth) UpdateSecret(secretCT, secretNonce []byte) {
	wo.oauthSecretCT = secretCT
	wo.oauthSecretNonce = secretNonce
	wo.updatedAt = time.Now().UTC()
	wo.raise(NewWebOAuthUpdatedEvent(wo.id.String(), "secret_updated"))
}

func (wo *WebOAuth) UpdateDEK(dekWrapped, dekNonce []byte) {
	wo.dekWrapped = dekWrapped
	wo.dekNonce = dekNonce
	wo.dekRotatedAt = time.Now().UTC()
	wo.updatedAt = time.Now().UTC()
	wo.raise(NewWebOAuthUpdatedEvent(wo.id.String(), "dek_rotated"))
}

func (wo *WebOAuth) UpdateRedirectURI(redirectURI string) {
	wo.redirectURI = redirectURI
	wo.updatedAt = time.Now().UTC()
	wo.raise(NewWebOAuthUpdatedEvent(wo.id.String(), "redirect_uris_updated"))
}

func (wo *WebOAuth) UpdateScopes(scopes []string) {
	wo.scopes = scopes
	wo.updatedAt = time.Now().UTC()
	wo.raise(NewWebOAuthUpdatedEvent(wo.id.String(), "scopes_updated"))
}

// Event methods

func (wo *WebOAuth) raise(e shared.DomainEvent) {
	wo.events = append(wo.events, e)
}

func (wo *WebOAuth) PullEvents() []shared.DomainEvent {
	events := wo.events
	wo.events = nil
	return events
}
