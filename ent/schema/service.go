package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Service holds the schema definition for the Service entity.
type Service struct {
	ent.Schema
}

// Fields of the Service.
func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique(),
		field.UUID("public_id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.String("name").
			NotEmpty(),
		field.String("description").
			Optional(),
		field.Bool("is_active").
			Default(true),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Service.
func (Service) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("client_apps", ServiceClient.Type),
	}
}

// Indexes of the Service.
func (Service) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("public_id").
			Unique(),
		index.Fields("name").
			Unique(),
	}
}
