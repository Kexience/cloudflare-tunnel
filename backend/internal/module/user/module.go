package user

import (
	"cloudflared-tunnel/internal/module/user/repo"
	"cloudflared-tunnel/internal/module/user/svc"
	"cloudflared-tunnel/internal/module/user/ui/api"
	"cloudflared-tunnel/internal/module/user/ui/api/ctrl"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"user",
	fx.Provide(
		repo.NewUserRepo,
		svc.NewUserSvc,
		ctrl.NewCtrl,
		api.NewRouter,
	),
)
