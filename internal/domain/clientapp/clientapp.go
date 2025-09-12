package clientapp

import (
	"time"

	"github.com/google/uuid"
	clientappval "github.com/mandacode-com/serengeti-integrated/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
)

type ClientApp struct {
	id          clientappval.ID
	publicID    clientappval.PublicID
	serviceID   serviceval.ID
	name        string
	description *string
	secretHash  clientappval.SecretHash
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
	events      []shared.DomainEvent
}

func NewClientApp(
	id clientappval.ID,
	publicID clientappval.PublicID,
	serviceID serviceval.ID,
	name string,
	description *string,
	secretHash clientappval.SecretHash,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *ClientApp {
	return &ClientApp{
		id:          id,
		publicID:    publicID,
		serviceID:   serviceID,
		name:        name,
		description: description,
		secretHash:  secretHash,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		events:      make([]shared.DomainEvent, 0),
	}
}

func DraftClientApp(serviceID serviceval.ID, name string, description *string) (*ClientApp, string, error) {
	now := time.Now().UTC()
	id := clientappval.NewID(0) // ID will be set by the database
	publicID := clientappval.NewPublicID(uuid.New())

	plainSecret, secretHash, err := clientappval.GenerateSecret()
	if err != nil {
		return nil, "", err
	}

	ca := NewClientApp(
		id,
		publicID,
		serviceID,
		name,
		description,
		secretHash,
		true,
		now,
		now,
	)

	ca.raise(NewClientAppCreatedEvent(ca.publicID.String(), ca.name))
	return ca, plainSecret, nil
}

// Getter methods

func (ca *ClientApp) ID() clientappval.ID {
	return ca.id
}

func (ca *ClientApp) PublicID() clientappval.PublicID {
	return ca.publicID
}

func (ca *ClientApp) ServiceID() serviceval.ID {
	return ca.serviceID
}

func (ca *ClientApp) Name() string {
	return ca.name
}

func (ca *ClientApp) Description() *string {
	return ca.description
}

func (ca *ClientApp) SecretHash() clientappval.SecretHash {
	return ca.secretHash
}

func (ca *ClientApp) IsActive() bool {
	return ca.isActive
}

func (ca *ClientApp) CreatedAt() time.Time {
	return ca.createdAt
}

func (ca *ClientApp) UpdatedAt() time.Time {
	return ca.updatedAt
}

// Business methods

func (ca *ClientApp) UpdateName(name string) {
	if ca.name != name {
		ca.name = name
		ca.updatedAt = time.Now().UTC()
		ca.raise(NewClientAppUpdatedEvent(ca.publicID.String(), "name_changed"))
	}
}

func (ca *ClientApp) UpdateDescription(description *string) {
	ca.description = description
	ca.updatedAt = time.Now().UTC()
	ca.raise(NewClientAppUpdatedEvent(ca.publicID.String(), "description_changed"))
}

func (ca *ClientApp) RegenerateSecret() (string, error) {
	plainSecret, secretHash, err := clientappval.GenerateSecret()
	if err != nil {
		return "", err
	}

	ca.secretHash = secretHash
	ca.updatedAt = time.Now().UTC()
	ca.raise(NewClientAppUpdatedEvent(ca.publicID.String(), "secret_regenerated"))

	return plainSecret, nil
}

func (ca *ClientApp) VerifySecret(plainSecret string) bool {
	return ca.secretHash.VerifySecret(plainSecret)
}

func (ca *ClientApp) Activate() {
	if !ca.isActive {
		ca.isActive = true
		ca.updatedAt = time.Now().UTC()
		ca.raise(NewClientAppUpdatedEvent(ca.publicID.String(), "activated"))
	}
}

func (ca *ClientApp) Deactivate() {
	if ca.isActive {
		ca.isActive = false
		ca.updatedAt = time.Now().UTC()
		ca.raise(NewClientAppUpdatedEvent(ca.publicID.String(), "deactivated"))
	}
}

// Event methods

func (ca *ClientApp) raise(e shared.DomainEvent) {
	ca.events = append(ca.events, e)
}

func (ca *ClientApp) PullEvents() []shared.DomainEvent {
	events := ca.events
	ca.events = nil
	return events
}
