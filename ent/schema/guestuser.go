package schema

import (
	"entgo.io/ent/schema/edge"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// GuestUser holds the schema definition for the GuestUser entity.
type GuestUser struct {
	ent.Schema
}

// Fields of the GuestUser.
func (GuestUser) Fields() []ent.Field {
	return []ent.Field{
		field.String("token").
			NotEmpty().
			Unique(),
		field.Int("version_number"),
		field.String("operation_system"),
		field.String("user_agent"),
		field.String("display_details"),
		field.Time("created_at").
			Default(time.Now()),
		field.Time("deleted_at").Optional().Nillable(),
	}
}

// Edges of the GuestUser.
func (GuestUser) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("followed", User.Type),
		//edge.To("scanned", Story.Type),
	}
}
