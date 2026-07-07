package tunnel

import (
	"cloudflared-tunnel/internal/module/tunnel/svc"
	"cloudflared-tunnel/internal/module/tunnel/ui/api"
	"cloudflared-tunnel/internal/module/tunnel/ui/api/ctrl"
	"cloudflared-tunnel/pkg/cloudflare"

	"go.uber.org/fx"
)

// Module 隧道管理模块
var Module = fx.Module(
	"tunnel",
	fx.Provide(
		cloudflare.NewTunnelClient,
		cloudflare.NewDNSClient,
		fx.Annotate(
			svc.NewSvc,
			fx.As(new(svc.TunnelSvc)),
			fx.As(new(svc.DNSsvc)),
		),
		ctrl.NewCtrl,
		api.NewRouter,
	),
)
