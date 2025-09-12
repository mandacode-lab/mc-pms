package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
)

// WebOAuth holds the schema definition for the WebOAuth entity.
type WebOAuth struct {
	ent.Schema
}

// Fields of the WebOAuth.
func (WebOAuth) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique(),
		field.Int64("client_app_id").
			NonNegative(),
		field.Enum("provider").
			GoType(shared.Provider("")),
		field.String("oauth_client_id").
			NotEmpty(),
		field.Bytes("oauth_secret_ct").
			NotEmpty(),
		field.Bytes("oauth_secret_nonce").
			NotEmpty(),
		field.Bytes("dek_wrapped").
			NotEmpty(),
		field.Bytes("dek_nonce").
			NotEmpty(),
		field.Time("dek_rotated_at").
			Default(time.Now),
		field.String("redirect_uri").
			Optional(),
		field.Strings("scopes").
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the WebOAuth.
func (WebOAuth) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("client_app", ClientApp.Type).
			Ref("web_oauths").
			Field("client_app_id").
			Unique().
			Required(),
	}
}

// Indexes of the WebOAuth.
func (WebOAuth) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("client_app_id", "provider").
			Unique(),
	}
}

