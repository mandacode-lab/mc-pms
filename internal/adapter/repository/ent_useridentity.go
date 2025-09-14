package repository

import (
	"context"

	"github.com/mandacode-com/serengeti-integrated/ent"
	entservice "github.com/mandacode-com/serengeti-integrated/ent/service"
	entuseridentity "github.com/mandacode-com/serengeti-integrated/ent/useridentity"
	serviceval "github.com/mandacode-com/serengeti-integrated/internal/domain/service/value"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/useridentity"
	useridentityval "github.com/mandacode-com/serengeti-integrated/internal/domain/useridentity/value"
	"github.com/mandacode-com/serengeti-integrated/internal/port/out"
)

type EntUserIdentityRepository struct {
	client *ent.Client
}

func NewEntUserIdentityRepository(client *ent.Client) out.UserIdentityRepository {
	return &EntUserIdentityRepository{
		client: client,
	}
}

func (r *EntUserIdentityRepository) Create(ctx context.Context, tx out.Tx, identityEntity *useridentity.UserIdentity) (*useridentity.UserIdentity, error) {
	var builder *ent.UserIdentityCreate

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return nil, err
		}
		builder = entTx.UserIdentity.Create()
	} else {
		builder = r.client.UserIdentity.Create()
	}

	entIdentity, err := builder.
		SetPublicID(identityEntity.PublicID().Value()).
		SetServiceID(identityEntity.ServiceID().Value()).
		SetProviderID(identityEntity.ProviderID()).
		SetProvider(identityEntity.Provider()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomain(entIdentity), nil
}

func (r *EntUserIdentityRepository) Update(ctx context.Context, tx out.Tx, identityEntity *useridentity.UserIdentity) error {
	var builder *ent.UserIdentityUpdateOne

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		builder = entTx.UserIdentity.UpdateOneID(identityEntity.ID().Value())
	} else {
		builder = r.client.UserIdentity.UpdateOneID(identityEntity.ID().Value())
	}

	_, err := builder.
		SetProviderID(identityEntity.ProviderID()).
		SetProvider(identityEntity.Provider()).
		Save(ctx)

	return err
}

func (r *EntUserIdentityRepository) FindByID(ctx context.Context, id useridentityval.ID) (*useridentity.UserIdentity, error) {
	entIdentity, err := r.client.UserIdentity.Query().Where(entuseridentity.ID(id.Value())).WithService().Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entIdentity), nil
}

func (r *EntUserIdentityRepository) Delete(ctx context.Context, tx out.Tx, identityEntity *useridentity.UserIdentity) error {
	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		return entTx.UserIdentity.DeleteOneID(identityEntity.ID().Value()).Exec(ctx)
	}
	return r.client.UserIdentity.DeleteOneID(identityEntity.ID().Value()).Exec(ctx)
}

func (r *EntUserIdentityRepository) toDomain(entIdentity *ent.UserIdentity) *useridentity.UserIdentity {
	id := useridentityval.NewID(entIdentity.ID)
	publicID := useridentityval.NewPublicID(entIdentity.PublicID)
	// Access service ID through the loaded edge
	serviceID := serviceval.NewID(entIdentity.Edges.Service.ID)

	return useridentity.NewUserIdentity(
		id,
		publicID,
		serviceID,
		entIdentity.ProviderID,
		entIdentity.Provider,
		entIdentity.CreatedAt,
		entIdentity.UpdatedAt,
	)
}

type EntUserIdentityQueryRepository struct {
	client *ent.Client
}

func NewEntUserIdentityQueryRepository(client *ent.Client) out.UserIdentityQueryRepository {
	return &EntUserIdentityQueryRepository{
		client: client,
	}
}

func (r *EntUserIdentityQueryRepository) FindByID(ctx context.Context, id useridentityval.ID) (*useridentity.UserIdentity, error) {
	entIdentity, err := r.client.UserIdentity.Get(ctx, id.Value())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entIdentity), nil
}

func (r *EntUserIdentityQueryRepository) FindByPublicID(ctx context.Context, publicID useridentityval.PublicID) (*useridentity.UserIdentity, error) {
	entIdentity, err := r.client.UserIdentity.Query().
		Where(entuseridentity.PublicID(publicID.Value())).
		WithService().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entIdentity), nil
}

func (r *EntUserIdentityQueryRepository) FindByProviderID(ctx context.Context, serviceID serviceval.ID, providerID string) (*useridentity.UserIdentity, error) {
	entIdentity, err := r.client.UserIdentity.Query().
		Where(
			entuseridentity.HasServiceWith(entservice.ID(serviceID.Value())),
			entuseridentity.ProviderID(providerID),
		).
		WithService().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entIdentity), nil
}

func (r *EntUserIdentityQueryRepository) FindByServiceID(ctx context.Context, serviceID serviceval.ID) ([]*useridentity.UserIdentity, error) {
	entIdentities, err := r.client.UserIdentity.Query().
		Where(entuseridentity.HasServiceWith(entservice.ID(serviceID.Value()))).
		WithService().
		All(ctx)
	if err != nil {
		return nil, err
	}

	identities := make([]*useridentity.UserIdentity, 0, len(entIdentities))
	for _, entIdentity := range entIdentities {
		identities = append(identities, r.toDomain(entIdentity))
	}

	return identities, nil
}

func (r *EntUserIdentityQueryRepository) List(ctx context.Context, filter *out.UserIdentityListFilter, options *out.UserIdentityListOptions) ([]*useridentity.UserIdentity, int, error) {
	query := r.client.UserIdentity.Query().WithService()

	if filter != nil {
		if filter.ServiceID != nil {
			query = query.Where(entuseridentity.HasServiceWith(entservice.ID(filter.ServiceID.Value())))
		}
		if filter.ProviderID != nil {
			query = query.Where(entuseridentity.ProviderIDContains(*filter.ProviderID))
		}
		if filter.Provider != nil {
			query = query.Where(entuseridentity.ProviderEQ(*filter.Provider))
		}
	}

	if options != nil {
		switch options.Order {
		case out.UserIdentityListOrderCreatedAtAsc:
			query = query.Order(ent.Asc(entuseridentity.FieldCreatedAt))
		case out.UserIdentityListOrderCreatedAtDesc:
			query = query.Order(ent.Desc(entuseridentity.FieldCreatedAt))
		}

		if options.Limit > 0 {
			query = query.Limit(options.Limit)
		}
		if options.Offset > 0 {
			query = query.Offset(options.Offset)
		}
	}

	entIdentities, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	count := len(entIdentities)
	identities := make([]*useridentity.UserIdentity, 0, count)
	for _, entIdentity := range entIdentities {
		identities = append(identities, r.toDomain(entIdentity))
	}

	return identities, count, nil
}

func (r *EntUserIdentityQueryRepository) toDomain(entIdentity *ent.UserIdentity) *useridentity.UserIdentity {
	id := useridentityval.NewID(entIdentity.ID)
	publicID := useridentityval.NewPublicID(entIdentity.PublicID)
	// Access service ID through the loaded edge
	serviceID := serviceval.NewID(entIdentity.Edges.Service.ID)

	return useridentity.NewUserIdentity(
		id,
		publicID,
		serviceID,
		entIdentity.ProviderID,
		entIdentity.Provider,
		entIdentity.CreatedAt,
		entIdentity.UpdatedAt,
	)
}

