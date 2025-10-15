package svcclient

import (
	"time"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
)

type SvcClient struct {
	id          ID
	publicID    PublicID
	name        Name
	description Description
	secretHash  SecretHash
	serviceID   service.ID
	isActive    bool
	createdAt   time.Time
	updatedAt   time.Time
}

func NewSvcClient(
	id ID,
	publicID PublicID,
	name Name,
	description Description,
	secretHash SecretHash,
	serviceID service.ID,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
) *SvcClient {
	return &SvcClient{
		id:          id,
		publicID:    publicID,
		name:        name,
		description: description,
		secretHash:  secretHash,
		serviceID:   serviceID,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// Getter methods

func (sc *SvcClient) ID() ID {
	return sc.id
}

func (sc *SvcClient) PublicID() PublicID {
	return sc.publicID
}

func (sc *SvcClient) Name() Name {
	return sc.name
}

func (sc *SvcClient) Description() Description {
	return sc.description
}

func (sc *SvcClient) SecretHash() SecretHash {
	return sc.secretHash
}

func (sc *SvcClient) ServiceID() service.ID {
	return sc.serviceID
}

func (sc *SvcClient) IsActive() bool {
	return sc.isActive
}

func (sc *SvcClient) CreatedAt() time.Time {
	return sc.createdAt
}

func (sc *SvcClient) UpdatedAt() time.Time {
	return sc.updatedAt
}
