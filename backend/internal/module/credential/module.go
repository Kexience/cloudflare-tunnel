package credential

import (
	"encoding/hex"

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
			// 尝试 hex 解码
			secretBytes, err := hex.DecodeString(cfg.Credential.Secret)
			if err != nil {
				// 不是 hex 字符串，直接使用原始字符串
				return types.CredentialSecret(cfg.Credential.Secret), nil
			}
			return types.CredentialSecret(secretBytes), nil
		},
		cloudflare.NewValidator,
		svc.NewCredentialSvc,
		ctrl.NewCtrl,
		api.NewRouter,
	),
)
