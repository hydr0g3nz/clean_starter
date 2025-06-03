package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("transaction_id").MaxLen(100).Comment("TransactionID จากระบบ TGW format: Ydddhhxxxxxx"),
		field.String("source_transaction_id").MaxLen(100).Comment("Transaction ID ของระบบต้นทาง"),
		field.String("terminal_id").MaxLen(50).Comment("terminal number"),
	}
}

// Edges of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		// Composite unique index on all three fields
		index.Fields(
			"transaction_id",
			"source_transaction_id",
			"terminal_id",
		).Unique(),
	}
}
