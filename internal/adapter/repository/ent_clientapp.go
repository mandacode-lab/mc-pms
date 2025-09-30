package repository

import (
	"context"
	"errors"

	"github.com/mandacode-com/mandacode-ssam/ent"
	entclientapp "github.com/mandacode-com/mandacode-ssam/ent/clientapp"
	entservice "github.com/mandacode-com/mandacode-ssam/ent/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/mandacode-ssam/pkg/utils"
)

type EntClientAppRepository struct {
	client *ent.Client
}

func NewEntClientAppRepository(client *ent.Client) out.ClientAppRepository {
	return &EntClientAppRepository{
		client: client,
	}
}

func (r *EntClientAppRepository) Create(ctx context.Context, tx out.Tx, clientEntity *entity.ClientApp) (*entity.ClientApp, error) {
	var builder *ent.ClientAppCreate

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return nil, err
		}
		builder = entTx.ClientApp.Create()
	} else {
		return nil, errors.New("creating ClientApp requires a transaction")
	}

	entClient, err := builder.
		SetPublicID(clientEntity.PublicID().Value()).
		SetServiceID(clientEntity.ServiceID().Value()).
		SetName(clientEntity.Name()).
		SetNillableDescription(clientEntity.Description()).
		SetSecretHash(clientEntity.SecretHash()).
		SetIsActive(clientEntity.IsActive()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomain(entClient)
}

func (r *EntClientAppRepository) Update(ctx context.Context, tx out.Tx, clientEntity *entity.ClientApp) error {
	var builder *ent.ClientAppUpdateOne

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		builder = entTx.ClientApp.UpdateOneID(clientEntity.ID().Value())
	} else {
		builder = r.client.ClientApp.UpdateOneID(clientEntity.ID().Value())
	}

	_, err := builder.
		SetName(clientEntity.Name()).
		SetNillableDescription(clientEntity.Description()).
		SetSecretHash(clientEntity.SecretHash()).
		SetIsActive(clientEntity.IsActive()).
		Save(ctx)

	return err
}

func (r *EntClientAppRepository) FindByID(ctx context.Context, id vo.ClientAppID) (*entity.ClientApp, error) {
	entClient, err := r.client.ClientApp.Get(ctx, id.Value())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(entClient)
}

func (r *EntClientAppRepository) Delete(ctx context.Context, tx out.Tx, clientEntity *entity.ClientApp) error {
	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		return entTx.ClientApp.DeleteOneID(clientEntity.ID().Value()).Exec(ctx)
	}
	return r.client.ClientApp.DeleteOneID(clientEntity.ID().Value()).Exec(ctx)
}

func (r *EntClientAppRepository) toDomain(entClient *ent.ClientApp) (*entity.ClientApp, error) {
	id := vo.NewClientAppID(entClient.ID)
	publicID := vo.NewClientAppPublicID(entClient.PublicID)
	serviceID := vo.NewServiceID(entClient.ServiceID)
	return entity.NewClientApp(
		id,
		publicID,
		serviceID,
		entClient.Name,
		utils.StringNil(entClient.Description),
		entClient.SecretHash,
		entClient.IsActive,
		entClient.CreatedAt,
		entClient.UpdatedAt,
	), nil
}

type EntClientAppQueryRepository struct {
	client *ent.Client
}

func NewEntClientAppQueryRepository(client *ent.Client) out.ClientAppQueryRepository {
	return &EntClientAppQueryRepository{
		client: client,
	}
}

func (r *EntClientAppQueryRepository) FindByID(ctx context.Context, id vo.ClientAppID) (*entity.ClientApp, error) {
	entClient, err := r.client.ClientApp.Query().
		Where(entclientapp.ID(id.Value())).
		WithService().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entClient), nil
}

func (r *EntClientAppQueryRepository) FindByPublicID(ctx context.Context, publicID vo.ClientAppPublicID) (*entity.ClientApp, error) {
	entClient, err := r.client.ClientApp.Query().
		Where(entclientapp.PublicID(publicID.Value())).
		WithService().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entClient), nil
}

func (r *EntClientAppQueryRepository) FindByName(ctx context.Context, name string) (*entity.ClientApp, error) {
	entClient, err := r.client.ClientApp.Query().
		Where(entclientapp.Name(name)).
		WithService().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entClient), nil
}

func (r *EntClientAppQueryRepository) List(ctx context.Context, filter *out.ClientAppListFilter, options *out.ClientAppListOptions) ([]*entity.ClientApp, int, error) {
	query := r.client.ClientApp.Query().WithService()

	if filter != nil {
		if filter.ServiceID != nil {
			query = query.Where(entclientapp.HasServiceWith(entservice.ID(filter.ServiceID.Value())))
		}
		if filter.Name != nil {
			query = query.Where(entclientapp.NameContains(*filter.Name))
		}
		if filter.IsActive != nil {
			query = query.Where(entclientapp.IsActive(*filter.IsActive))
		}
	}

	if options != nil {
		switch options.Order {
		case out.ClientAppListOrderNameAsc:
			query = query.Order(ent.Asc(entclientapp.FieldName))
		case out.ClientAppListOrderNameDesc:
			query = query.Order(ent.Desc(entclientapp.FieldName))
		}

		if options.Limit > 0 {
			query = query.Limit(options.Limit)
		}
		if options.Offset > 0 {
			query = query.Offset(options.Offset)
		}
	}

	entClients, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	count := len(entClients)
	clients := make([]*entity.ClientApp, 0, count)
	for _, entClient := range entClients {
		clients = append(clients, r.toDomain(entClient))
	}

	return clients, count, nil
}

func (r *EntClientAppQueryRepository) toDomain(entClient *ent.ClientApp) *entity.ClientApp {
	id := vo.NewClientAppID(entClient.ID)
	publicID := vo.NewClientAppPublicID(entClient.PublicID)
	serviceID := vo.NewServiceID(entClient.Edges.Service.ID)
	var description *string
	if entClient.Description != "" {
		description = &entClient.Description
	}

	return entity.NewClientApp(
		id,
		publicID,
		serviceID,
		entClient.Name,
		description,
		entClient.SecretHash,
		entClient.IsActive,
		entClient.CreatedAt,
		entClient.UpdatedAt,
	)
}
