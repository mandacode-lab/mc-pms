package client

import (
	"time"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
)

type Client struct {
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
) *Client {
	return &Client{
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

func (sc *Client) ID() ID {
	return sc.id
}

func (sc *Client) PublicID() PublicID {
	return sc.publicID
}

func (sc *Client) Name() Name {
	return sc.name
}

func (sc *Client) Description() Description {
	return sc.description
}

func (sc *Client) SecretHash() SecretHash {
	return sc.secretHash
}

func (sc *Client) ServiceID() service.ID {
	return sc.serviceID
}

func (sc *Client) IsActive() bool {
	return sc.isActive
}

func (sc *Client) CreatedAt() time.Time {
	return sc.createdAt
}

func (sc *Client) UpdatedAt() time.Time {
	return sc.updatedAt
}
