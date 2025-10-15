package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// ServiceClient holds the schema definition for the ServiceClient entity.
type ServiceClient struct {
	ent.Schema
}

// Fields of the ClientApp.
func (ServiceClient) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique(),
		field.Int64("service_id"),
		field.UUID("public_id", uuid.UUID{}).
			Default(uuid.New).
			Unique(),
		field.Bytes("secret_hash").
			NotEmpty(),
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

// Edges of the ClientApp.
func (ServiceClient) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("service", Service.Type).
			Ref("client_apps").
			Field("service_id").
			Unique().
			Required(),
	}
}

// Indexes of the ClientApp.
func (ServiceClient) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("public_id").
			Unique(),
		index.Fields("service_id", "name").
			Unique(),
	}
}
