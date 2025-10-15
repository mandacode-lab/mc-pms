package repository

import (
	"context"
	"errors"

	"github.com/mandacode-com/mandacode-ssam/ent"
	entservice "github.com/mandacode-com/mandacode-ssam/ent/service"
	"github.com/mandacode-com/mandacode-ssam/ent/svcclient"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/client"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
)

type EntClientAppRepository struct {
	client *ent.Client
}

func NewEntClientAppRepository(client *ent.Client) client.ClientRepository {
	return &EntClientAppRepository{
		client: client,
	}
}

func (r *EntClientAppRepository) Create(ctx context.Context, tx tx.Tx, clientModel *client.Client) (*client.Client, error) {
	var builder *ent.SvcClientCreate

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return nil, err
		}
		builder = entTx.SvcClient.Create()
	} else {
		return nil, errors.New("creating ClientApp requires a transaction")
	}

	entClient, err := builder.
		SetPublicID(clientModel.PublicID().Value()).
		SetServiceID(clientModel.ServiceID().Value()).
		SetName(clientModel.Name().Value()).
		SetNillableDescription(clientModel.Description().Value()).
		SetSecretHash(clientModel.SecretHash()).
		SetIsActive(clientModel.IsActive()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomain(entClient)
}

func (r *EntClientAppRepository) Update(ctx context.Context, tx tx.Tx, clientEntity *client.Client) error {
	var builder *ent.SvcClientUpdateOne

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		builder = entTx.SvcClient.UpdateOneID(clientEntity.ID().Value())
	} else {
		builder = r.client.SvcClient.UpdateOneID(clientEntity.ID().Value())
	}

	_, err := builder.
		SetName(clientEntity.Name().Value()).
		SetNillableDescription(clientEntity.Description().Value()).
		SetSecretHash(clientEntity.SecretHash()).
		SetIsActive(clientEntity.IsActive()).
		Save(ctx)

	return err
}

func (r *EntClientAppRepository) FindByID(ctx context.Context, id client.ID) (*client.Client, error) {
	entClient, err := r.client.SvcClient.Get(ctx, id.Value())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(entClient)
}

func (r *EntClientAppRepository) Delete(ctx context.Context, tx tx.Tx, clientEntity *client.Client) error {
	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		return entTx.SvcClient.DeleteOneID(clientEntity.ID().Value()).Exec(ctx)
	}
	return r.client.SvcClient.DeleteOneID(clientEntity.ID().Value()).Exec(ctx)
}

func (r *EntClientAppRepository) toDomain(entClient *ent.SvcClient) (*client.Client, error) {
	id, err := client.NewID(entClient.ID)
	if err != nil {
		return nil, err
	}
	publicID, err := client.NewPublicID(entClient.PublicID)
	if err != nil {
		return nil, err
	}
	name, err := client.NewName(entClient.Name)
	if err != nil {
		return nil, err
	}
	description, err := client.NewDescription(entClient.Description)
	if err != nil {
		return nil, err
	}
	serviceID, err := service.NewID(entClient.ServiceID)
	if err != nil {
		return nil, err
	}
	secretHash, err := client.NewSecretHash(entClient.SecretHash)
	if err != nil {
		return nil, err
	}

	return client.NewClient(
		id,
		publicID,
		name,
		description,
		secretHash,
		serviceID,
		entClient.IsActive,
		entClient.CreatedAt,
		entClient.UpdatedAt,
	), nil
}

type EntClientAppQueryRepository struct {
	client *ent.Client
}

func NewEntClientAppQueryRepository(client *ent.Client) client.ClientQueryRepository {
	return &EntClientAppQueryRepository{
		client: client,
	}
}

func (r *EntClientAppQueryRepository) FindByID(ctx context.Context, id client.ID) (*client.Client, error) {
	entClient, err := r.client.SvcClient.Query().
		Where(svcclient.ID(id.Value())).
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

func (r *EntClientAppQueryRepository) FindByPublicID(ctx context.Context, publicID client.PublicID) (*client.Client, error) {
	entClient, err := r.client.SvcClient.Query().
		Where(svcclient.PublicID(publicID.Value())).
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

func (r *EntClientAppQueryRepository) FindByName(ctx context.Context, name client.Name) (*client.Client, error) {
	entClient, err := r.client.SvcClient.Query().
		Where(svcclient.Name(name.Value())).
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

func (r *EntClientAppQueryRepository) List(ctx context.Context, filter *client.ClientListFilter, options *client.ClientListOptions) ([]*client.Client, int, error) {
	query := r.client.SvcClient.Query().WithService()

	if filter != nil {
		if filter.ServiceID != nil {
			query = query.Where(svcclient.HasServiceWith(entservice.ID(filter.ServiceID.Value())))
		}
		if filter.NameContains != nil {
			query = query.Where(svcclient.NameContains(*filter.NameContains))
		}
		if filter.IsActive != nil {
			query = query.Where(svcclient.IsActive(*filter.IsActive))
		}
	}

	if options != nil {
		switch options.Order {
		case client.ClientListOrderNameAsc:
			query = query.Order(ent.Asc(svcclient.FieldName))
		case client.ClientListOrderNameDesc:
			query = query.Order(ent.Desc(svcclient.FieldName))
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
	clients := make([]*client.Client, 0, count)
	for _, entClient := range entClients {
		client, err := r.toDomain(entClient)
		if err != nil {
			return nil, 0, err
		}
		clients = append(clients, client)
	}

	return clients, count, nil
}

func (r *EntClientAppQueryRepository) toDomain(entClient *ent.SvcClient) (*client.Client, error) {
	id, err := client.NewID(entClient.ID)
	if err != nil {
		return nil, err
	}
	publicID, err := client.NewPublicID(entClient.PublicID)
	if err != nil {
		return nil, err
	}
	name, err := client.NewName(entClient.Name)
	if err != nil {
		return nil, err
	}
	description, err := client.NewDescription(entClient.Description)
	if err != nil {
		return nil, err
	}
	secretHash, err := client.NewSecretHash(entClient.SecretHash)
	if err != nil {
		return nil, err
	}
	serviceID, err := service.NewID(entClient.Edges.Service.ID)
	if err != nil {
		return nil, err
	}

	return client.NewClient(
		id,
		publicID,
		name,
		description,
		secretHash,
		serviceID,
		entClient.IsActive,
		entClient.CreatedAt,
		entClient.UpdatedAt,
	), nil
}
