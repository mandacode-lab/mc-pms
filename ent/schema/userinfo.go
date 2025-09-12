package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// UserInfo holds the schema definition for the UserInfo entity.
type UserInfo struct {
	ent.Schema
}

// Fields of the UserInfo.
func (UserInfo) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique(),
		field.String("nickname").
			NotEmpty(),
		field.String("email").
			Optional(),
		field.Bytes("raw_data").
			Optional(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the UserInfo.
func (UserInfo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user_identity", UserIdentity.Type).
			Ref("user_info").
			Unique().
			Required(),
	}
}
