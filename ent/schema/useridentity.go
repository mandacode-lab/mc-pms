package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"github.com/mandacode-com/serengeti-integrated/internal/domain/shared"
)

// UserIdentity holds the schema definition for the UserIdentity entity.
type UserIdentity struct {
	ent.Schema
}

// Fields of the UserIdentity.
func (UserIdentity) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique(),
		field.Int64("service_id"),

		field.UUID("public_id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),

		field.String("provider_id").
			NotEmpty(),
		field.Enum("provider").
			GoType(shared.Provider("")),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the UserIdentity.
func (UserIdentity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("service", Service.Type).
			Ref("user_identities").
			Field("service_id").
			Unique().
			Required(),
		edge.To("user_info", UserInfo.Type).
			Unique(),
	}
}

// Indexes of the UserIdentity.
func (UserIdentity) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("service_id", "provider", "provider_id").Unique(),
	}
}
