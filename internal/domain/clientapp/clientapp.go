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

func DraftClientApp(serviceID serviceval.ID, name string, description *string, plainSecret []byte, secretHash clientappval.SecretHash) (*ClientApp, []byte) {
	now := time.Now().UTC()
	id := clientappval.NewID(0) // ID will be set by the database
	publicID := clientappval.NewPublicID(uuid.New())

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
	return ca, plainSecret
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

func (ca *ClientApp) RegenerateSecret(plainSecret []byte, secretHash clientappval.SecretHash) []byte {
	ca.secretHash = secretHash
	ca.updatedAt = time.Now().UTC()
	ca.raise(NewClientAppUpdatedEvent(ca.publicID.String(), "secret_regenerated"))

	return plainSecret
}

func (ca *ClientApp) VerifySecret(plainSecret string, hash []byte) bool {
	return ca.secretHash.VerifySecret(plainSecret, hash)
}

func (ca *ClientApp) VerifySecretBytes(plainSecret []byte, hash []byte) bool {
	return ca.secretHash.VerifySecretBytes(plainSecret, hash)
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
