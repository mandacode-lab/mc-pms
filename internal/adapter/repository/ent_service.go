package repository

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/ent"
	entservice "github.com/mandacode-com/mandacode-ssam/ent/service"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/entity"
	vo "github.com/mandacode-com/mandacode-ssam/internal/domain/value_object"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
	"github.com/mandacode-com/mandacode-ssam/pkg/utils"
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

	entService, err := builder.
		SetPublicID(serviceEntity.PublicID().Value()).
		SetName(serviceEntity.Name().Value()).
		SetNillableDescription(serviceEntity.Description()).
		SetIsActive(serviceEntity.IsActive()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomain(entService), nil
}

func (r *EntServiceRepository) Update(ctx context.Context, tx out.Tx, serviceEntity *entity.Service) error {
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
		SetNillableDescription(serviceEntity.Description()).
		SetIsActive(serviceEntity.IsActive()).
		Save(ctx)

	return err
}

func (r *EntServiceRepository) FindByID(ctx context.Context, id vo.ServiceID) (*entity.Service, error) {
	entService, err := r.client.Service.Get(ctx, id.Value())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entService), nil
}

func (r *EntServiceRepository) Delete(ctx context.Context, tx out.Tx, serviceEntity *entity.Service) error {
	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		return entTx.Service.DeleteOneID(serviceEntity.ID().Value()).Exec(ctx)
	}
	return r.client.Service.DeleteOneID(serviceEntity.ID().Value()).Exec(ctx)
}

func (r *EntServiceRepository) toDomain(entService *ent.Service) *entity.Service {
	id := vo.NewServiceID(entService.ID)
	publicID := vo.NewServicePublicID(entService.PublicID)
	name, _ := vo.NewServiceName(entService.Name)

	var description *string
	if entService.Description != "" {
		description = &entService.Description
	}

	return entity.NewService(
		id,
		publicID,
		name,
		description,
		entService.IsActive,
		entService.CreatedAt,
		entService.UpdatedAt,
	)
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
	entService, err := r.client.Service.Get(ctx, id.Value())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entService)
}

func (r *EntServiceQueryRepository) FindByPublicID(ctx context.Context, publicID vo.ServicePublicID) (*entity.Service, error) {
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

func (r *EntServiceQueryRepository) FindByName(ctx context.Context, name vo.ServiceName) (*entity.Service, error) {
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
	id := vo.NewServiceID(entService.ID)
	publicID := vo.NewServicePublicID(entService.PublicID)
	name, err := vo.NewServiceName(entService.Name)
	if err != nil {
		return nil, err
	}

	description := utils.StringNil(entService.Description)

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
