package clientapp

import (
	"time"

	"github.com/google/uuid"
	clientappval "github.com/mandacode-com/mandacode-ssam/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
)

type ClientApp struct {
	id          clientappval.ID
	publicID    clientappval.PublicID
	serviceID   serviceval.ID
	name        string
	description *string
	secretHash  []byte
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
}

func NewClientApp(
	id clientappval.ID,
	publicID clientappval.PublicID,
	serviceID serviceval.ID,
	name string,
	description *string,
	secretHash []byte,
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
	}
}

func DraftClientApp(serviceID serviceval.ID, name string, description *string, secretHash []byte) *ClientApp {
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

	return ca
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

func (ca *ClientApp) SecretHash() []byte {
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
	}
}

func (ca *ClientApp) UpdateDescription(description *string) {
	ca.description = description
	ca.updatedAt = time.Now().UTC()
}

func (ca *ClientApp) RegenerateSecret(secretHash []byte) {
	ca.secretHash = secretHash
	ca.updatedAt = time.Now().UTC()
}

func (ca *ClientApp) Activate() {
	if !ca.isActive {
		ca.isActive = true
		ca.updatedAt = time.Now().UTC()
	}
}

func (ca *ClientApp) Deactivate() {
	if ca.isActive {
		ca.isActive = false
		ca.updatedAt = time.Now().UTC()
	}
}
