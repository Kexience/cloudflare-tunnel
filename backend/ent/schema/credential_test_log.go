package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type CredentialTestLog struct {
	ent.Schema
}

func (CredentialTestLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:   "credential_test_logs",
			Options: "COMMENT='凭证测试日志表'",
		},
		entsql.WithComments(true),
	}
}

func (CredentialTestLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Positive().Immutable().Comment("日志ID"),
		field.String("status").NotEmpty().Comment("测试结果: success/failed"),
		field.String("error_message").Optional().Comment("错误信息"),
		field.Time("tested_at").Default(time.Now).Immutable().Comment("测试时间"),
	}
}

func (CredentialTestLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("credential", Credential.Type).
			Ref("test_logs").
			Unique().
			Required().
			Comment("关联凭证"),
	}
}
