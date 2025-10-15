package repository

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/ent"
	entservice "github.com/mandacode-com/mandacode-ssam/ent/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/tx"
)

type EntServiceRepository struct {
	client *ent.Client
}

func NewEntServiceRepository(client *ent.Client) service.ServiceRepository {
	return &EntServiceRepository{
		client: client,
	}
}

func (r *EntServiceRepository) Create(ctx context.Context, tx tx.Tx, serviceEntity *service.Service) (*service.Service, error) {
	var builder *ent.ServiceCreate

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return nil, err
		}
		builder = entTx.Service.Create()
	} else {
		builder = r.client.Service.Create()
	}

	entService, err := builder.
		SetPublicID(string(serviceEntity.PublicID().Value())).
		SetName(serviceEntity.Name().Value()).
		SetNillableDescription(serviceEntity.Description().Value()).
		SetIsActive(serviceEntity.IsActive()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomain(entService)
}

func (r *EntServiceRepository) Update(ctx context.Context, tx tx.Tx, serviceEntity *service.Service) error {
	var builder *ent.ServiceUpdateOne

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		builder = entTx.Service.UpdateOneID(serviceEntity.ID().Value())
	} else {
		builder = r.client.Service.UpdateOneID(serviceEntity.ID().Value())
	}

	_, err := builder.
		SetName(serviceEntity.Name().Value()).
		SetNillableDescription(serviceEntity.Description().Value()).
		SetIsActive(serviceEntity.IsActive()).
		Save(ctx)

	return err
}

func (r *EntServiceRepository) FindByID(ctx context.Context, id service.ID) (*service.Service, error) {
	entService, err := r.client.Service.Get(ctx, id.Value())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entService)
}

func (r *EntServiceRepository) Delete(ctx context.Context, tx tx.Tx, serviceEntity *service.Service) error {
	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		return entTx.Service.DeleteOneID(serviceEntity.ID().Value()).Exec(ctx)
	}
	return r.client.Service.DeleteOneID(serviceEntity.ID().Value()).Exec(ctx)
}

func (r *EntServiceRepository) toDomain(entService *ent.Service) (*service.Service, error) {
	id, err := service.NewID(entService.ID)
	if err != nil {
		return nil, err
	}
	publicID, err := service.NewPublicID(entService.PublicID)
	if err != nil {
		return nil, err
	}
	name, err := service.NewName(entService.Name)
	if err != nil {
		return nil, err
	}

	description, err := service.NewDescription(entService.Description)
	if err != nil {
		return nil, err
	}

	return service.NewService(
		id,
		publicID,
		name,
		description,
		entService.IsActive,
		entService.CreatedAt,
		entService.UpdatedAt,
	), nil
}

type EntServiceQueryRepository struct {
	client *ent.Client
}

func NewEntServiceQueryRepository(client *ent.Client) service.ServiceQueryRepository {
	return &EntServiceQueryRepository{
		client: client,
	}
}

func (r *EntServiceQueryRepository) FindByID(ctx context.Context, id service.ID) (*service.Service, error) {
	entService, err := r.client.Service.Get(ctx, id.Value())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entService)
}

func (r *EntServiceQueryRepository) FindByPublicID(ctx context.Context, publicID service.PublicID) (*service.Service, error) {
	entService, err := r.client.Service.Query().
		Where(entservice.PublicID(publicID.Value())).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entService)
}

func (r *EntServiceQueryRepository) FindByName(ctx context.Context, name service.Name) (*service.Service, error) {
	entService, err := r.client.Service.Query().
		Where(entservice.Name(name.Value())).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entService)
}

func (r *EntServiceQueryRepository) List(ctx context.Context, filter *service.ServiceListFilter, options *service.ServiceListOptions) ([]*service.Service, int, error) {
	query := r.client.Service.Query()

	if filter != nil {
		if filter.Name != nil {
			query = query.Where(entservice.NameContains(*filter.Name))
		}
		if filter.IsActive != nil {
			query = query.Where(entservice.IsActive(*filter.IsActive))
		}
	}

	if options != nil {
		switch options.Order {
		case service.ServiceListOrderNameAsc:
			query = query.Order(ent.Asc(entservice.FieldName))
		case service.ServiceListOrderNameDesc:
			query = query.Order(ent.Desc(entservice.FieldName))
		}

		if options.Limit > 0 {
			query = query.Limit(options.Limit)
		}
		if options.Offset > 0 {
			query = query.Offset(options.Offset)
		}
	}

	entServices, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	count := len(entServices)
	services := make([]*service.Service, 0, count)
	for i := range entServices {
		service, err := r.toDomain(entServices[i])
		if err != nil {
			return nil, 0, err
		}
		services = append(services, service)
	}

	return services, count, nil
}

func (r *EntServiceQueryRepository) toDomain(entService *ent.Service) (*service.Service, error) {
	id, err := service.NewID(entService.ID)
	if err != nil {
		return nil, err
	}
	publicID, err := service.NewPublicID(entService.PublicID)
	if err != nil {
		return nil, err
	}
	name, err := service.NewName(entService.Name)
	if err != nil {
		return nil, err
	}

	description, err := service.NewDescription(entService.Description)
	if err != nil {
		return nil, err
	}

	return service.NewService(
		id,
		publicID,
		name,
		description,
		entService.IsActive,
		entService.CreatedAt,
		entService.UpdatedAt,
	), nil
}
