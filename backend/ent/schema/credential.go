package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Credential struct {
	ent.Schema
}

func (Credential) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:   "credentials",
			Options: "COMMENT='Cloudflare 凭证表'",
		},
		entsql.WithComments(true),
	}
}

func (Credential) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Positive().Immutable().Comment("凭证ID"),
		field.String("name").NotEmpty().Comment("凭证名称"),
		field.String("api_token").NotEmpty().Comment("加密后的 API Token"),
		field.String("account_id").NotEmpty().Comment("Cloudflare Account ID"),
		field.Bool("is_default").Default(false).Comment("是否为默认凭证"),
		field.Time("created_at").Default(time.Now).Immutable().Comment("创建时间"),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).Comment("更新时间"),
	}
}

func (Credential) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("credentials").
			Unique().
			Required().
			Comment("所属用户"),
	}
}

func (Credential) Indexes() []ent.Index {
	return []ent.Index{
		index.Edges("owner"),
	}
}
