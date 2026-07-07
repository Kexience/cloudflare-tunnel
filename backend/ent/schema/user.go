package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:   "users",
			Options: "COMMENT='用户表'",
		},
		entsql.WithComments(true),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Positive().Immutable().Comment("用户ID"),
		field.String("nickname").NotEmpty().Unique().Comment("昵称"),
		field.String("username").NotEmpty().Unique().Comment("用户名"),
		field.String("password").NotEmpty().Comment("密码"),
		field.String("email").NotEmpty().Unique().Comment("邮箱"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username").Unique(),
		index.Fields("email").Unique(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("credentials", Credential.Type),
	}
}
