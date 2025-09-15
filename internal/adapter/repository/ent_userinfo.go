package repository

import (
	"context"

	"github.com/mandacode-com/mandacode-ssam/ent"
	entuseridentity "github.com/mandacode-com/mandacode-ssam/ent/useridentity"
	entuserinfo "github.com/mandacode-com/mandacode-ssam/ent/userinfo"
	useridentityval "github.com/mandacode-com/mandacode-ssam/internal/domain/useridentity/value"
	"github.com/mandacode-com/mandacode-ssam/internal/domain/userinfo"
	userinfoval "github.com/mandacode-com/mandacode-ssam/internal/domain/userinfo/value"
	"github.com/mandacode-com/mandacode-ssam/internal/port/out"
)

type EntUserInfoRepository struct {
	client *ent.Client
}

func NewEntUserInfoRepository(client *ent.Client) out.UserInfoRepository {
	return &EntUserInfoRepository{
		client: client,
	}
}

func (r *EntUserInfoRepository) Create(ctx context.Context, tx out.Tx, infoEntity *userinfo.UserInfo) (*userinfo.UserInfo, error) {
	var builder *ent.UserInfoCreate

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return nil, err
		}
		builder = entTx.UserInfo.Create()
	} else {
		builder = r.client.UserInfo.Create()
	}

	entInfo, err := builder.
		SetPublicID(infoEntity.PublicID().Value()).
		SetUserIdentityID(infoEntity.UserIdentityID().Value()).
		SetNickname(infoEntity.Nickname()).
		SetEmail(infoEntity.Email()).
		SetRawData(infoEntity.RawData()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomain(entInfo), nil
}

func (r *EntUserInfoRepository) Update(ctx context.Context, tx out.Tx, infoEntity *userinfo.UserInfo) error {
	var builder *ent.UserInfoUpdateOne

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		builder = entTx.UserInfo.UpdateOneID(infoEntity.ID().Value())
	} else {
		builder = r.client.UserInfo.UpdateOneID(infoEntity.ID().Value())
	}

	_, err := builder.
		SetNickname(infoEntity.Nickname()).
		SetEmail(infoEntity.Email()).
		SetRawData(infoEntity.RawData()).
		Save(ctx)

	return err
}

func (r *EntUserInfoRepository) FindByID(ctx context.Context, id userinfoval.ID) (*userinfo.UserInfo, error) {
	entInfo, err := r.client.UserInfo.Get(ctx, id.Value())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entInfo), nil
}

func (r *EntUserInfoRepository) Delete(ctx context.Context, tx out.Tx, infoEntity *userinfo.UserInfo) error {
	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		return entTx.UserInfo.DeleteOneID(infoEntity.ID().Value()).Exec(ctx)
	}
	return r.client.UserInfo.DeleteOneID(infoEntity.ID().Value()).Exec(ctx)
}

func (r *EntUserInfoRepository) toDomain(entInfo *ent.UserInfo) *userinfo.UserInfo {
	id := userinfoval.NewID(entInfo.ID)
	publicID := userinfoval.NewPublicID(entInfo.PublicID)
	userIdentityID := useridentityval.NewID(entInfo.Edges.UserIdentity.ID)

	return userinfo.NewUserInfo(
		id,
		publicID,
		userIdentityID,
		entInfo.Nickname,
		entInfo.Email,
		entInfo.RawData,
		entInfo.CreatedAt,
		entInfo.UpdatedAt,
	)
}

type EntUserInfoQueryRepository struct {
	client *ent.Client
}

func NewEntUserInfoQueryRepository(client *ent.Client) out.UserInfoQueryRepository {
	return &EntUserInfoQueryRepository{
		client: client,
	}
}

func (r *EntUserInfoQueryRepository) FindByID(ctx context.Context, id userinfoval.ID) (*userinfo.UserInfo, error) {
	entInfo, err := r.client.UserInfo.Query().
		Where(entuserinfo.ID(id.Value())).
		WithUserIdentity().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entInfo), nil
}

func (r *EntUserInfoQueryRepository) FindByPublicID(ctx context.Context, publicID userinfoval.PublicID) (*userinfo.UserInfo, error) {
	entInfo, err := r.client.UserInfo.Query().
		Where(entuserinfo.PublicID(publicID.Value())).
		WithUserIdentity().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entInfo), nil
}

func (r *EntUserInfoQueryRepository) FindByUserIdentityID(ctx context.Context, userIdentityID useridentityval.ID) (*userinfo.UserInfo, error) {
	entInfo, err := r.client.UserInfo.Query().
		Where(entuserinfo.HasUserIdentityWith(entuseridentity.ID(userIdentityID.Value()))).
		WithUserIdentity().
		WithUserIdentity().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entInfo), nil
}

func (r *EntUserInfoQueryRepository) FindByEmail(ctx context.Context, email string) (*userinfo.UserInfo, error) {
	entInfo, err := r.client.UserInfo.Query().
		Where(entuserinfo.Email(email)).
		WithUserIdentity().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entInfo), nil
}

func (r *EntUserInfoQueryRepository) FindByNickname(ctx context.Context, nickname string) (*userinfo.UserInfo, error) {
	entInfo, err := r.client.UserInfo.Query().
		Where(entuserinfo.Nickname(nickname)).
		WithUserIdentity().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entInfo), nil
}

func (r *EntUserInfoQueryRepository) List(ctx context.Context, filter *out.UserInfoListFilter, options *out.UserInfoListOptions) ([]*userinfo.UserInfo, int, error) {
	query := r.client.UserInfo.Query().WithUserIdentity()

	if filter != nil {
		if filter.UserIdentityID != nil {
			query = query.Where(entuserinfo.HasUserIdentityWith(entuseridentity.ID(filter.UserIdentityID.Value())))
		}
		if filter.Email != nil {
			query = query.Where(entuserinfo.EmailContains(*filter.Email))
		}
		if filter.Nickname != nil {
			query = query.Where(entuserinfo.NicknameContains(*filter.Nickname))
		}
	}

	if options != nil {
		switch options.Order {
		case out.UserInfoListOrderNicknameAsc:
			query = query.Order(ent.Asc(entuserinfo.FieldNickname))
		case out.UserInfoListOrderNicknameDesc:
			query = query.Order(ent.Desc(entuserinfo.FieldNickname))
		}

		if options.Limit > 0 {
			query = query.Limit(options.Limit)
		}
		if options.Offset > 0 {
			query = query.Offset(options.Offset)
		}
	}

	entInfos, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	count := len(entInfos)
	infos := make([]*userinfo.UserInfo, 0, count)
	for _, entInfo := range entInfos {
		infos = append(infos, r.toDomain(entInfo))
	}

	return infos, count, nil
}

func (r *EntUserInfoQueryRepository) toDomain(entInfo *ent.UserInfo) *userinfo.UserInfo {
	id := userinfoval.NewID(entInfo.ID)
	publicID := userinfoval.NewPublicID(entInfo.PublicID)
	userIdentityID := useridentityval.NewID(entInfo.Edges.UserIdentity.ID)

	return userinfo.NewUserInfo(
		id,
		publicID,
		userIdentityID,
		entInfo.Nickname,
		entInfo.Email,
		entInfo.RawData,
		entInfo.CreatedAt,
		entInfo.UpdatedAt,
	)
}

