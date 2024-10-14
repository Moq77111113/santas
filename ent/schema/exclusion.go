package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Exclusion holds the schema definition for the Exclusion entity.
type Exclusion struct {
	ent.Schema
}

// Fields of the Exclusion.
func (Exclusion) Fields() []ent.Field {
	return []ent.Field{
		field.Int("group_id"),
		field.Int("member_id"),
		field.Int("exclude_id"),
	}
}

// Edges of the Exclusion.
func (Exclusion) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("group", Group.Type).Field("group_id").Unique().Required(),
		edge.To("member", Member.Type).Field("member_id").Unique().Required(),
		edge.To("exclude", Member.Type).Field("exclude_id").Unique().Required(),
	}
}

// Indexes of the Exclusion.
func (Exclusion) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("group_id", "member_id", "exclude_id").Unique(),
	}
}
