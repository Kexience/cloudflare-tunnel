package infra

import (
	"cloudflared-tunnel/internal/infra/db"
	"cloudflared-tunnel/internal/infra/logger"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"infra",
	fx.Provide(
		logger.NewLogger,
		db.NewClient,
	),
)
