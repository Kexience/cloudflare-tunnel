package credential

import (
	"cloudflared-tunnel/internal/config"
	"cloudflared-tunnel/internal/module/credential/repo"
	"cloudflared-tunnel/internal/module/credential/svc"
	"cloudflared-tunnel/internal/module/credential/ui/api"
	"cloudflared-tunnel/internal/module/credential/ui/api/ctrl"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"credential",
	fx.Provide(
		repo.NewCredentialRepo,
		func(cfg *config.Config) ([]byte, error) {
			return []byte(cfg.Credential.Secret), nil
		},
		svc.NewCredentialSvc,
		ctrl.NewCtrl,
		api.NewRouter,
	),
)
