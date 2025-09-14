package repository

import (
	"context"

	"github.com/mandacode-com/serengeti/ent"
	entclientapp "github.com/mandacode-com/serengeti/ent/clientapp"
	entweboauth "github.com/mandacode-com/serengeti/ent/weboauth"
	clientappval "github.com/mandacode-com/serengeti/internal/domain/clientapp/value"
	"github.com/mandacode-com/serengeti/internal/domain/shared"
	"github.com/mandacode-com/serengeti/internal/domain/weboauth"
	weboauthval "github.com/mandacode-com/serengeti/internal/domain/weboauth/value"
	"github.com/mandacode-com/serengeti/internal/port/out"
	"github.com/mandacode-com/serengeti/pkg/utils"
)

type EntWebOAuthRepository struct {
	client *ent.Client
}

func NewEntWebOAuthRepository(client *ent.Client) out.WebOAuthRepository {
	return &EntWebOAuthRepository{
		client: client,
	}
}

func (r *EntWebOAuthRepository) Create(ctx context.Context, tx out.Tx, oauthEntity *weboauth.WebOAuth) (*weboauth.WebOAuth, error) {
	var builder *ent.WebOAuthCreate

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return nil, err
		}
		builder = entTx.WebOAuth.Create()
	} else {
		builder = r.client.WebOAuth.Create()
	}

	entOAuth, err := builder.
		SetClientAppID(oauthEntity.ClientAppID().Value()).
		SetProvider(oauthEntity.Provider()).
		SetOauthClientID(oauthEntity.OAuthClientID()).
		SetOauthSecretCt(oauthEntity.OAuthSecretCT()).
		SetOauthSecretNonce(oauthEntity.OAuthSecretNonce()).
		SetDekWrapped(oauthEntity.DEKWrapped()).
		SetDekNonce(oauthEntity.DEKNonce()).
		SetDekRotatedAt(oauthEntity.DEKRotatedAt()).
		SetNillableRedirectURI(utils.StringNil(oauthEntity.RedirectURI())).
		SetScopes(oauthEntity.Scopes()).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomain(entOAuth), nil
}

func (r *EntWebOAuthRepository) Update(ctx context.Context, tx out.Tx, oauthEntity *weboauth.WebOAuth) error {
	var builder *ent.WebOAuthUpdateOne

	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		builder = entTx.WebOAuth.UpdateOneID(oauthEntity.ID().Value())
	} else {
		builder = r.client.WebOAuth.UpdateOneID(oauthEntity.ID().Value())
	}

	_, err := builder.
		SetProvider(oauthEntity.Provider()).
		SetOauthClientID(oauthEntity.OAuthClientID()).
		SetOauthSecretCt(oauthEntity.OAuthSecretCT()).
		SetOauthSecretNonce(oauthEntity.OAuthSecretNonce()).
		SetDekWrapped(oauthEntity.DEKWrapped()).
		SetDekNonce(oauthEntity.DEKNonce()).
		SetDekRotatedAt(oauthEntity.DEKRotatedAt()).
		SetNillableRedirectURI(utils.StringNil(oauthEntity.RedirectURI())).
		SetScopes(oauthEntity.Scopes()).
		Save(ctx)

	return err
}

func (r *EntWebOAuthRepository) FindByID(ctx context.Context, id weboauthval.ID) (*weboauth.WebOAuth, error) {
	entOAuth, err := r.client.WebOAuth.Get(ctx, id.Value())
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entOAuth), nil
}

func (r *EntWebOAuthRepository) Delete(ctx context.Context, tx out.Tx, oauthEntity *weboauth.WebOAuth) error {
	if tx != nil {
		entTx, err := asEntTx(tx)
		if err != nil {
			return err
		}
		return entTx.WebOAuth.DeleteOneID(oauthEntity.ID().Value()).Exec(ctx)
	}
	return r.client.WebOAuth.DeleteOneID(oauthEntity.ID().Value()).Exec(ctx)
}

func (r *EntWebOAuthRepository) toDomain(entOAuth *ent.WebOAuth) *weboauth.WebOAuth {
	id := weboauthval.NewID(entOAuth.ID)
	clientAppID := clientappval.NewID(entOAuth.Edges.ClientApp.ID)
	provider := shared.Provider(entOAuth.Provider)

	return weboauth.NewWebOAuth(
		id,
		clientAppID,
		provider,
		entOAuth.OauthClientID,
		entOAuth.OauthSecretCt,
		entOAuth.OauthSecretNonce,
		entOAuth.DekWrapped,
		entOAuth.DekNonce,
		entOAuth.DekRotatedAt,
		entOAuth.RedirectURI,
		entOAuth.Scopes,
		entOAuth.CreatedAt,
		entOAuth.UpdatedAt,
	)
}

type EntWebOAuthQueryRepository struct {
	client *ent.Client
}

func NewEntWebOAuthQueryRepository(client *ent.Client) out.WebOAuthQueryRepository {
	return &EntWebOAuthQueryRepository{
		client: client,
	}
}

func (r *EntWebOAuthQueryRepository) FindByID(ctx context.Context, id weboauthval.ID) (*weboauth.WebOAuth, error) {
	entOAuth, err := r.client.WebOAuth.Query().
		Where(entweboauth.ID(id.Value())).
		WithClientApp().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entOAuth), nil
}

func (r *EntWebOAuthQueryRepository) FindByProvider(ctx context.Context, clientAppID clientappval.ID, provider shared.Provider) (*weboauth.WebOAuth, error) {
	entOAuth, err := r.client.WebOAuth.Query().
		Where(
			entweboauth.HasClientAppWith(entclientapp.ID(clientAppID.Value())),
			entweboauth.ProviderEQ(provider),
		).
		WithClientApp().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomain(entOAuth), nil
}

func (r *EntWebOAuthQueryRepository) FindByClientAppID(ctx context.Context, clientAppID clientappval.ID) ([]*weboauth.WebOAuth, error) {
	entOAuths, err := r.client.WebOAuth.Query().
		Where(entweboauth.HasClientAppWith(entclientapp.ID(clientAppID.Value()))).
		WithClientApp().
		All(ctx)
	if err != nil {
		return nil, err
	}

	oauths := make([]*weboauth.WebOAuth, 0, len(entOAuths))
	for _, entOAuth := range entOAuths {
		oauths = append(oauths, r.toDomain(entOAuth))
	}

	return oauths, nil
}

func (r *EntWebOAuthQueryRepository) List(ctx context.Context, filter *out.WebOAuthListFilter, options *out.WebOAuthListOptions) ([]*weboauth.WebOAuth, int, error) {
	query := r.client.WebOAuth.Query().WithClientApp()

	if filter != nil {
		if filter.ClientAppID != nil {
			query = query.Where(entweboauth.HasClientAppWith(entclientapp.ID(filter.ClientAppID.Value())))
		}
		if filter.Provider != nil {
			query = query.Where(entweboauth.ProviderEQ(*filter.Provider))
		}
	}

	if options != nil {
		switch options.Order {
		case out.WebOAuthListOrderProviderAsc:
			query = query.Order(ent.Asc(entweboauth.FieldProvider))
		case out.WebOAuthListOrderProviderDesc:
			query = query.Order(ent.Desc(entweboauth.FieldProvider))
		}

		if options.Limit > 0 {
			query = query.Limit(options.Limit)
		}
		if options.Offset > 0 {
			query = query.Offset(options.Offset)
		}
	}

	entOAuths, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	count := len(entOAuths)
	oauths := make([]*weboauth.WebOAuth, 0, count)
	for _, entOAuth := range entOAuths {
		oauths = append(oauths, r.toDomain(entOAuth))
	}

	return oauths, count, nil
}

func (r *EntWebOAuthQueryRepository) toDomain(entOAuth *ent.WebOAuth) *weboauth.WebOAuth {
	id := weboauthval.NewID(entOAuth.ID)
	clientAppID := clientappval.NewID(entOAuth.Edges.ClientApp.ID)
	provider := shared.Provider(entOAuth.Provider)

	return weboauth.NewWebOAuth(
		id,
		clientAppID,
		provider,
		entOAuth.OauthClientID,
		entOAuth.OauthSecretCt,
		entOAuth.OauthSecretNonce,
		entOAuth.DekWrapped,
		entOAuth.DekNonce,
		entOAuth.DekRotatedAt,
		entOAuth.RedirectURI,
		entOAuth.Scopes,
		entOAuth.CreatedAt,
		entOAuth.UpdatedAt,
	)
}
