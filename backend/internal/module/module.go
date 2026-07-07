package module

import (
	"cloudflared-tunnel/internal/module/credential"
	"cloudflared-tunnel/internal/module/tunnel"
	"cloudflared-tunnel/internal/module/user"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"Module",
	user.Module,
	credential.Module,
	tunnel.Module,
)
