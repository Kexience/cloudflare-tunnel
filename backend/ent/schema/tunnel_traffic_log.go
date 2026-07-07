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

type TunnelTrafficLog struct {
	ent.Schema
}

func (TunnelTrafficLog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:   "tunnel_traffic_logs",
			Options: "COMMENT='隧道流量统计记录表'",
		},
		entsql.WithComments(true),
	}
}

func (TunnelTrafficLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique().Positive().Immutable().Comment("记录ID"),
		field.String("tunnel_id").NotEmpty().Comment("隧道ID"),
		field.Int64("bytes_in").Default(0).Comment("入站字节数"),
		field.Int64("bytes_out").Default(0).Comment("出站字节数"),
		field.Int64("total_requests").Default(0).Comment("总请求数"),
		field.Time("recorded_at").Default(time.Now).Immutable().Comment("记录时间"),
	}
}

func (TunnelTrafficLog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("traffic_logs").
			Unique().
			Required().
			Comment("所属用户"),
	}
}

func (TunnelTrafficLog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tunnel_id", "recorded_at"),
	}
}
