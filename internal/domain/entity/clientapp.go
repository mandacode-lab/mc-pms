package entity

import (
	"time"

	"github.com/google/uuid"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
)

type ClientApp struct {
	id          vo.ClientAppID
	publicID    vo.ClientAppPublicID
	serviceID   vo.ServiceID
	name        string
	description vo.ClientAppDescription
	secretHash  []byte
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
}

func NewClientApp(
	id vo.ClientAppID,
	publicID vo.ClientAppPublicID,
	serviceID vo.ServiceID,
	name string,
	description vo.ClientAppDescription,
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

func DraftClientApp(serviceID vo.ServiceID, name string, description vo.ClientAppDescription, secretHash []byte) (*ClientApp, error) {
	now := time.Now().UTC()

	// ID 1 is temporary, will be replaced by DB
	id, err := vo.NewClientAppID(1)
	if err != nil {
		return nil, err
	}

	publicID, err := vo.NewClientAppPublicID(uuid.New().String())
	if err != nil {
		return nil, err
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

	return ca, nil
}

// Getter methods

func (ca *ClientApp) ID() vo.ClientAppID {
	return ca.id
}

func (ca *ClientApp) PublicID() vo.ClientAppPublicID {
	return ca.publicID
}

func (ca *ClientApp) ServiceID() vo.ServiceID {
	return ca.serviceID
}

func (ca *ClientApp) Name() string {
	return ca.name
}

func (ca *ClientApp) Description() vo.ClientAppDescription {
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

func (ca *ClientApp) UpdateDescription(description vo.ClientAppDescription) {
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
