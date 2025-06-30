package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type Alpha struct {
	ent.Schema
}

func (Alpha) Fields() []ent.Field {
	return []ent.Field{
		field.String("Uuid").
			NotEmpty().
			StorageKey("Uuid"),

		field.String("FirstInsert").
			StorageKey("FirstInsert").
			SchemaType(map[string]string{"mysql": "datetime(6)"}),

		field.String("LastUpdate").
			StorageKey("LastUpdate").
			SchemaType(map[string]string{"mysql": "datetime(6)"}),

		field.String("Animal").
			NotEmpty().
			StorageKey("Animal"),

		field.String("BigNumber").
			StorageKey("BigNumber"),

		field.String("TestField").
			Optional().
			StorageKey("test_field"),
	}
}
