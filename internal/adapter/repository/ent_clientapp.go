package repository

import (
	"context"

	"github.com/mandacode-com/serengeti/ent"
	entclientapp "github.com/mandacode-com/serengeti/ent/clientapp"
	entservice "github.com/mandacode-com/serengeti/ent/service"
	"github.com/mandacode-com/serengeti/internal/domain/clientapp"
	clientappval "github.com/mandacode-com/serengeti/internal/domain/clientapp/value"
	serviceval "github.com/mandacode-com/serengeti/internal/domain/service/value"
	"github.com/mandacode-com/serengeti/internal/port/out"
)

type EntClientAppRepository struct {
	client *ent.Client
}

func NewEntClientAppRepository(client *ent.Client) out.ClientAppRepository {
	return &EntClientAppRepository{
		client: client,
	}
}

func (r *EntClientAppRepository) Create(ctx context.Context, tx out.Tx, clientEntity *clientapp.ClientApp) (*clientapp.ClientApp, error) {
	var builder *ent.ClientAppCreate

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return nil, err
		}
		builder = entTx.ClientApp.Create()
	} else {
		builder = r.client.ClientApp.Create()
	}

	entClient, err := builder.
		SetPublicID(clientEntity.PublicID().Value()).
		SetServiceID(clientEntity.ServiceID().Value()).
		SetName(clientEntity.Name()).
		SetNillableDescription(clientEntity.Description()).
		SetSecretHash(clientEntity.SecretHash().Value()).
		SetIsActive(clientEntity.IsActive()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomain(entClient), nil
}

func (r *EntClientAppRepository) Update(ctx context.Context, tx out.Tx, clientEntity *clientapp.ClientApp) error {
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
		SetSecretHash(clientEntity.SecretHash().Value()).
		SetIsActive(clientEntity.IsActive()).
		Save(ctx)

	return err
}

func (r *EntClientAppRepository) FindByID(ctx context.Context, id clientappval.ID) (*clientapp.ClientApp, error) {
	entClient, err := r.client.ClientApp.Get(ctx, id.Value())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entClient), nil
}

func (r *EntClientAppRepository) Delete(ctx context.Context, tx out.Tx, clientEntity *clientapp.ClientApp) error {
	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		return entTx.ClientApp.DeleteOneID(clientEntity.ID().Value()).Exec(ctx)
	}
	return r.client.ClientApp.DeleteOneID(clientEntity.ID().Value()).Exec(ctx)
}

func (r *EntClientAppRepository) toDomain(entClient *ent.ClientApp) *clientapp.ClientApp {
	id := clientappval.NewID(entClient.ID)
	publicID := clientappval.NewPublicID(entClient.PublicID)
	serviceID := serviceval.NewID(entClient.Edges.Service.ID)
	secretHash := clientappval.NewSecretHash(entClient.SecretHash)

	return clientapp.NewClientApp(
		id,
		publicID,
		serviceID,
		entClient.Name,
		&entClient.Description,
		secretHash,
		entClient.IsActive,
		entClient.CreatedAt,
		entClient.UpdatedAt,
	)
}

type EntClientAppQueryRepository struct {
	client *ent.Client
}

func NewEntClientAppQueryRepository(client *ent.Client) out.ClientAppQueryRepository {
	return &EntClientAppQueryRepository{
		client: client,
	}
}

func (r *EntClientAppQueryRepository) FindByID(ctx context.Context, id clientappval.ID) (*clientapp.ClientApp, error) {
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

func (r *EntClientAppQueryRepository) FindByPublicID(ctx context.Context, publicID clientappval.PublicID) (*clientapp.ClientApp, error) {
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

func (r *EntClientAppQueryRepository) FindByName(ctx context.Context, name string) (*clientapp.ClientApp, error) {
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

func (r *EntClientAppQueryRepository) FindByServiceID(ctx context.Context, serviceID serviceval.ID) ([]*clientapp.ClientApp, error) {
	entClients, err := r.client.ClientApp.Query().
		Where(entclientapp.HasServiceWith(entservice.ID(serviceID.Value()))).
		WithService().
		All(ctx)
	if err != nil {
		return nil, err
	}

	clients := make([]*clientapp.ClientApp, 0, len(entClients))
	for _, entClient := range entClients {
		clients = append(clients, r.toDomain(entClient))
	}

	return clients, nil
}

func (r *EntClientAppQueryRepository) List(ctx context.Context, filter *out.ClientAppListFilter, options *out.ClientAppListOptions) ([]*clientapp.ClientApp, int, error) {
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
	clients := make([]*clientapp.ClientApp, 0, count)
	for _, entClient := range entClients {
		clients = append(clients, r.toDomain(entClient))
	}

	return clients, count, nil
}

func (r *EntClientAppQueryRepository) toDomain(entClient *ent.ClientApp) *clientapp.ClientApp {
	id := clientappval.NewID(entClient.ID)
	publicID := clientappval.NewPublicID(entClient.PublicID)
	serviceID := serviceval.NewID(entClient.Edges.Service.ID)
	secretHash := clientappval.NewSecretHash(entClient.SecretHash)

	var description *string
	if entClient.Description != "" {
		description = &entClient.Description
	}

	return clientapp.NewClientApp(
		id,
		publicID,
		serviceID,
		entClient.Name,
		description,
		secretHash,
		entClient.IsActive,
		entClient.CreatedAt,
		entClient.UpdatedAt,
	)
}

