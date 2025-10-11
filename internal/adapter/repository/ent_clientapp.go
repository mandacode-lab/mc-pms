package repository

import (
	"context"
	"errors"

	"github.com/mandacode-com/mandacode-ssam/ent"
	entclientapp "github.com/mandacode-com/mandacode-ssam/ent/clientapp"
	entservice "github.com/mandacode-com/mandacode-ssam/ent/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
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

	publicIDUUID, err := clientEntity.PublicID().UUID()
	if err != nil {
		return nil, err
	}

	entClient, err := builder.
		SetPublicID(publicIDUUID).
		SetServiceID(clientEntity.ServiceID().Int64()).
		SetName(clientEntity.Name()).
		SetNillableDescription(clientEntity.Description().Ptr()).
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
		builder = entTx.ClientApp.UpdateOneID(clientEntity.ID().Int64())
	} else {
		builder = r.client.ClientApp.UpdateOneID(clientEntity.ID().Int64())
	}

	_, err := builder.
		SetName(clientEntity.Name()).
		SetNillableDescription(clientEntity.Description().Ptr()).
		SetSecretHash(clientEntity.SecretHash()).
		SetIsActive(clientEntity.IsActive()).
		Save(ctx)

	return err
}

func (r *EntClientAppRepository) FindByID(ctx context.Context, id vo.ClientAppID) (*entity.ClientApp, error) {
	entClient, err := r.client.ClientApp.Get(ctx, id.Int64())
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
		return entTx.ClientApp.DeleteOneID(clientEntity.ID().Int64()).Exec(ctx)
	}
	return r.client.ClientApp.DeleteOneID(clientEntity.ID().Int64()).Exec(ctx)
}

func (r *EntClientAppRepository) toDomain(entClient *ent.ClientApp) (*entity.ClientApp, error) {
	id, err := vo.NewClientAppID(entClient.ID)
	if err != nil {
		return nil, err
	}
	publicID, err := vo.NewClientAppPublicID(entClient.PublicID.String())
	if err != nil {
		return nil, err
	}
	serviceID, err := vo.NewServiceID(entClient.ServiceID)
	if err != nil {
		return nil, err
	}

	description, err := vo.NewClientAppDescription(entClient.Description)
	if err != nil {
		return nil, err
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
		Where(entclientapp.ID(id.Int64())).
		WithService().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entClient)
}

func (r *EntClientAppQueryRepository) FindByPublicID(ctx context.Context, publicID vo.ClientAppPublicID) (*entity.ClientApp, error) {
	publicIDUUID, err := publicID.UUID()
	if err != nil {
		return nil, err
	}
	entClient, err := r.client.ClientApp.Query().
		Where(entclientapp.PublicID(publicIDUUID)).
		WithService().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entClient)
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
	return r.toDomain(entClient)
}

func (r *EntClientAppQueryRepository) List(ctx context.Context, filter *out.ClientAppListFilter, options *out.ClientAppListOptions) ([]*entity.ClientApp, int, error) {
	query := r.client.ClientApp.Query().WithService()

	if filter != nil {
		if filter.ServiceID != nil {
			query = query.Where(entclientapp.HasServiceWith(entservice.ID(filter.ServiceID.Int64())))
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
		client, err := r.toDomain(entClient)
		if err != nil {
			return nil, 0, err
		}
		clients = append(clients, client)
	}

	return clients, count, nil
}

func (r *EntClientAppQueryRepository) toDomain(entClient *ent.ClientApp) (*entity.ClientApp, error) {
	id, err := vo.NewClientAppID(entClient.ID)
	if err != nil {
		return nil, err
	}
	publicID, err := vo.NewClientAppPublicID(entClient.PublicID.String())
	if err != nil {
		return nil, err
	}
	serviceID, err := vo.NewServiceID(entClient.Edges.Service.ID)
	if err != nil {
		return nil, err
	}

	description, err := vo.NewClientAppDescription(entClient.Description)
	if err != nil {
		return nil, err
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
	), nil
}
