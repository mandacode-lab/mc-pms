package repository

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/ent"
	entservice "github.com/mandacode-com/mandacode-ssam/ent/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/vo"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type EntServiceRepository struct {
	client *ent.Client
}

func NewEntServiceRepository(client *ent.Client) out.ServiceRepository {
	return &EntServiceRepository{
		client: client,
	}
}

func (r *EntServiceRepository) Create(ctx context.Context, tx out.Tx, serviceEntity *entity.Service) (*entity.Service, error) {
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

	publicIDUUID, _ := serviceEntity.PublicID().UUID()
	entService, err := builder.
		SetPublicID(publicIDUUID).
		SetName(serviceEntity.Name().String()).
		SetNillableDescription(serviceEntity.Description().Ptr()).
		SetIsActive(serviceEntity.IsActive()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomain(entService)
}

func (r *EntServiceRepository) Update(ctx context.Context, tx out.Tx, serviceEntity *entity.Service) error {
	var builder *ent.ServiceUpdateOne

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		builder = entTx.Service.UpdateOneID(serviceEntity.ID().Int64())
	} else {
		builder = r.client.Service.UpdateOneID(serviceEntity.ID().Int64())
	}

	_, err := builder.
		SetName(serviceEntity.Name().String()).
		SetNillableDescription(serviceEntity.Description().Ptr()).
		SetIsActive(serviceEntity.IsActive()).
		Save(ctx)

	return err
}

func (r *EntServiceRepository) FindByID(ctx context.Context, id vo.ServiceID) (*entity.Service, error) {
	entService, err := r.client.Service.Get(ctx, id.Int64())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entService)
}

func (r *EntServiceRepository) Delete(ctx context.Context, tx out.Tx, serviceEntity *entity.Service) error {
	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		return entTx.Service.DeleteOneID(serviceEntity.ID().Int64()).Exec(ctx)
	}
	return r.client.Service.DeleteOneID(serviceEntity.ID().Int64()).Exec(ctx)
}

func (r *EntServiceRepository) toDomain(entService *ent.Service) (*entity.Service, error) {
	id, err := vo.NewServiceID(entService.ID)
	if err != nil {
		return nil, err
	}
	publicID, err := vo.NewServicePublicID(entService.PublicID.String())
	if err != nil {
		return nil, err
	}
	name, err := vo.NewServiceName(entService.Name)
	if err != nil {
		return nil, err
	}

	description, err := vo.NewServiceDescription(entService.Description)
	if err != nil {
		return nil, err
	}

	return entity.NewService(
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

func NewEntServiceQueryRepository(client *ent.Client) out.ServiceQueryRepository {
	return &EntServiceQueryRepository{
		client: client,
	}
}

func (r *EntServiceQueryRepository) FindByID(ctx context.Context, id vo.ServiceID) (*entity.Service, error) {
	entService, err := r.client.Service.Get(ctx, id.Int64())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entService)
}

func (r *EntServiceQueryRepository) FindByPublicID(ctx context.Context, publicID vo.ServicePublicID) (*entity.Service, error) {
	publicIDUUID, err := publicID.UUID()
	if err != nil {
		return nil, err
	}
	entService, err := r.client.Service.Query().
		Where(entservice.PublicID(publicIDUUID)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entService)
}

func (r *EntServiceQueryRepository) FindByName(ctx context.Context, name vo.ServiceName) (*entity.Service, error) {
	entService, err := r.client.Service.Query().
		Where(entservice.Name(name.String())).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entService)
}

func (r *EntServiceQueryRepository) List(ctx context.Context, filter *out.ServiceListFilter, options *out.ServiceListOptions) ([]*entity.Service, int, error) {
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
		case out.ServiceListOrderNameAsc:
			query = query.Order(ent.Asc(entservice.FieldName))
		case out.ServiceListOrderNameDesc:
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
	services := make([]*entity.Service, 0, count)
	for i := range entServices {
		service, err := r.toDomain(entServices[i])
		if err != nil {
			return nil, 0, err
		}
		services = append(services, service)
	}

	return services, count, nil
}

func (r *EntServiceQueryRepository) toDomain(entService *ent.Service) (*entity.Service, error) {
	id, err := vo.NewServiceID(entService.ID)
	if err != nil {
		return nil, err
	}
	publicID, err := vo.NewServicePublicID(entService.PublicID.String())
	if err != nil {
		return nil, err
	}
	name, err := vo.NewServiceName(entService.Name)
	if err != nil {
		return nil, err
	}

	description, err := vo.NewServiceDescription(entService.Description)
	if err != nil {
		return nil, err
	}

	return entity.NewService(
		id,
		publicID,
		name,
		description,
		entService.IsActive,
		entService.CreatedAt,
		entService.UpdatedAt,
	), nil
}
