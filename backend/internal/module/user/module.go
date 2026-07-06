package user

import (
	"cloudflared-tunnel/internal/module/user/ui/api"
	"cloudflared-tunnel/internal/module/user/ui/api/ctrl"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		ctrl.NewCtrl,
		api.NewRouter,
	),
)
