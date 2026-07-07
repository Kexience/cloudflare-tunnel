package credential

import (
	"cloudflared-tunnel/internal/config"
	"cloudflared-tunnel/internal/module/credential/repo"
	"cloudflared-tunnel/internal/module/credential/svc"
	"cloudflared-tunnel/internal/module/credential/ui/api"
	"cloudflared-tunnel/internal/module/credential/ui/api/ctrl"
	"cloudflared-tunnel/internal/types"
	"cloudflared-tunnel/pkg/cloudflare"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"credential",
	fx.Provide(
		repo.NewCredentialRepo,
		func(cfg *config.Config) (types.CredentialSecret, error) {
			return types.CredentialSecret(cfg.Credential.Secret), nil
		},
		cloudflare.NewValidator,
		svc.NewCredentialSvc,
		ctrl.NewCtrl,
		api.NewRouter,
	),
)
