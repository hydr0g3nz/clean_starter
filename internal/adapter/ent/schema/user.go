package schema

import (
	"time"

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
		field.Int("id").
			Positive().
			Unique(),

		field.String("email").
			NotEmpty().
			Unique(),

		field.String("password_hash").
			NotEmpty().
			Sensitive(),

		field.Enum("role").
			Values("candidate", "company_hr", "admin").
			Comment("User role: candidate, company_hr, or admin"),

		field.Bool("is_active").
			Default(true).
			Comment("Whether the user account is active"),

		field.Bool("email_verified").
			Default(false).
			Comment("Whether the user's email has been verified"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("When the user was created"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("When the user was last updated"),

		field.Time("last_login_at").
			Optional().
			Nillable().
			Comment("When the user last logged in"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		// Unique index on email (already handled by field.Unique())
		index.Fields("email").Unique(),

		// Index on role for faster queries
		index.Fields("role"),

		// Index on is_active for faster queries
		index.Fields("is_active"),

		// Composite index for common queries
		index.Fields("role", "is_active"),
	}
}
