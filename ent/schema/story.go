package schema

import (
	"regexp"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Story holds the schema definition for the Story entity.
type Story struct {
	ent.Schema
}

// Fields of the Story.
func (Story) Fields() []ent.Field {
	return []ent.Field{
		field.String("storyName").
			NotEmpty(),
		field.String("backgroundColor").
			Match(regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}){1,2}$|^rgb\(\d{1,3},\s?\d{1,3},\s?\d{1,3}\)$`)).
			NotEmpty(),
		field.String("backgroundImage").
			NotEmpty(),
		field.Bool("isShareable").
			Default(true),
		field.String("attachedFile").
			Optional(),
		field.String("externalWebLink").
			Optional(),
		field.Time("createdAt").
			Default(time.Now),
		field.Time("updatedAt").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("deletedAt").
			Nillable().
			Optional(),
		field.Time("fromTime"),
		field.Time("toTime"),
		field.Uint("scanCount").
			Optional().
			Default(1),
		field.Enum("status").
			Values("private", "public").
			Default("private"),
	}
}

// Edges of the Story.
func (Story) Edges() []ent.Edge {
	return []ent.Edge{
		//edge.From("scannedby", GuestUser.Type).Ref("scanned"),
		edge.From("postedby", User.Type).Ref("posted").Unique(),
	}
}
