package out

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/internal/domain/useridentity"
	useridentityval "github.com/mandacode-com/mandacode-ssam/internal/domain/useridentity/value"
	serviceval "github.com/mandacode-com/mandacode-ssam/internal/domain/service/value"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/shared"
)

type UserIdentityListFilter struct {
	ServiceID    *serviceval.ID
	ProviderID   *string
	Provider     *shared.Provider
}

type UserIdentityListOrder string

const (
	UserIdentityListOrderCreatedAtAsc  UserIdentityListOrder = "created_at_asc"
	UserIdentityListOrderCreatedAtDesc UserIdentityListOrder = "created_at_desc"
)

type UserIdentityListOptions struct {
	Limit  int
	Offset int
	Order  UserIdentityListOrder
}


type UserIdentityRepository interface {
	Create(ctx context.Context, tx Tx, identity *useridentity.UserIdentity) (*useridentity.UserIdentity, error)
	Update(ctx context.Context, tx Tx, identity *useridentity.UserIdentity) error
	FindByID(ctx context.Context, id useridentityval.ID) (*useridentity.UserIdentity, error)
	Delete(ctx context.Context, tx Tx, identity *useridentity.UserIdentity) error
}

type UserIdentityQueryRepository interface {
	FindByID(ctx context.Context, id useridentityval.ID) (*useridentity.UserIdentity, error)
	FindByPublicID(ctx context.Context, publicID useridentityval.PublicID) (*useridentity.UserIdentity, error)
	FindByProviderID(ctx context.Context, serviceID serviceval.ID, providerID string) (*useridentity.UserIdentity, error)
	FindByServiceID(ctx context.Context, serviceID serviceval.ID) ([]*useridentity.UserIdentity, error)
	List(ctx context.Context, filter *UserIdentityListFilter, options *UserIdentityListOptions) ([]*useridentity.UserIdentity, int, error)
}