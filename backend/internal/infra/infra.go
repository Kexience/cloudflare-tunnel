package infra

import (
	"cloudflared-tunnel/internal/config"
	"cloudflared-tunnel/internal/infra/db"
	"cloudflared-tunnel/internal/infra/logger"
	"cloudflared-tunnel/pkg/core"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"infra",
	fx.Provide(
		logger.NewLogger,
		db.NewClient,
		func(cfg *config.Config) *core.JWT {
			return core.NewJWT(cfg.JWT.Secret, cfg.JWT.ExpireHour)
		},
	),
)
