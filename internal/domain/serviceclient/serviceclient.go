package serviceclient

import (
	"time"

	"github.com/google/uuid"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
	serviceclientval "github.com/mandacode-com/serengeti-integrated/internal/domain/serviceclient/value"
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
)

type ServiceClient struct {
	id          serviceclientval.ID
	publicID    serviceclientval.PublicID
	serviceID   serviceval.ID
	name        string
	description *string
	secretHash  serviceclientval.SecretHash
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
	events      []shared.DomainEvent
}

func NewServiceClient(
	id serviceclientval.ID,
	publicID serviceclientval.PublicID,
	serviceID serviceval.ID,
	name string,
	description *string,
	secretHash serviceclientval.SecretHash,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *ServiceClient {
	return &ServiceClient{
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

func DraftServiceClient(serviceID serviceval.ID, name string, description *string) (*ServiceClient, string, error) {
	now := time.Now().UTC()
	id := serviceclientval.NewID(0) // ID will be set by the database
	publicID := serviceclientval.NewPublicID(uuid.New())

	plainSecret, secretHash, err := serviceclientval.GenerateSecret()
	if err != nil {
		return nil, "", err
	}

	sc := NewServiceClient(
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

	sc.raise(NewServiceClientCreatedEvent(sc.publicID.String(), sc.name))
	return sc, plainSecret, nil
}

// Getter methods

func (sc *ServiceClient) ID() serviceclientval.ID {
	return sc.id
}

func (sc *ServiceClient) PublicID() serviceclientval.PublicID {
	return sc.publicID
}

func (sc *ServiceClient) ServiceID() serviceval.ID {
	return sc.serviceID
}

func (sc *ServiceClient) Name() string {
	return sc.name
}

func (sc *ServiceClient) Description() *string {
	return sc.description
}

func (sc *ServiceClient) SecretHash() serviceclientval.SecretHash {
	return sc.secretHash
}

func (sc *ServiceClient) IsActive() bool {
	return sc.isActive
}

func (sc *ServiceClient) CreatedAt() time.Time {
	return sc.createdAt
}

func (sc *ServiceClient) UpdatedAt() time.Time {
	return sc.updatedAt
}

// Business methods

func (sc *ServiceClient) UpdateName(name string) {
	if sc.name != name {
		sc.name = name
		sc.updatedAt = time.Now().UTC()
		sc.raise(NewServiceClientUpdatedEvent(sc.publicID.String(), "name_changed"))
	}
}

func (sc *ServiceClient) UpdateDescription(description *string) {
	sc.description = description
	sc.updatedAt = time.Now().UTC()
	sc.raise(NewServiceClientUpdatedEvent(sc.publicID.String(), "description_changed"))
}

func (sc *ServiceClient) RegenerateSecret() (string, error) {
	plainSecret, secretHash, err := serviceclientval.GenerateSecret()
	if err != nil {
		return "", err
	}

	sc.secretHash = secretHash
	sc.updatedAt = time.Now().UTC()
	sc.raise(NewServiceClientUpdatedEvent(sc.publicID.String(), "secret_regenerated"))

	return plainSecret, nil
}

func (sc *ServiceClient) VerifySecret(plainSecret string) bool {
	return sc.secretHash.VerifySecret(plainSecret)
}

func (sc *ServiceClient) Activate() {
	if !sc.isActive {
		sc.isActive = true
		sc.updatedAt = time.Now().UTC()
		sc.raise(NewServiceClientUpdatedEvent(sc.publicID.String(), "activated"))
	}
}

func (sc *ServiceClient) Deactivate() {
	if sc.isActive {
		sc.isActive = false
		sc.updatedAt = time.Now().UTC()
		sc.raise(NewServiceClientUpdatedEvent(sc.publicID.String(), "deactivated"))
	}
}

// Event methods

func (sc *ServiceClient) raise(e shared.DomainEvent) {
	sc.events = append(sc.events, e)
}

func (sc *ServiceClient) PullEvents() []shared.DomainEvent {
	events := sc.events
	sc.events = nil
	return events
}