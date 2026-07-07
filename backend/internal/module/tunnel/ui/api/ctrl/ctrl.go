package ctrl

import (
	"cloudflared-tunnel/internal/module/tunnel/svc"
)

// Ctrl 隧道管理控制器
type Ctrl struct {
	tunnelSvc svc.TunnelSvc
	dnsSvc    svc.DNSsvc
}

// NewCtrl 创建隧道管理控制器
func NewCtrl(tunnelSvc svc.TunnelSvc, dnsSvc svc.DNSsvc) *Ctrl {
	return &Ctrl{tunnelSvc: tunnelSvc, dnsSvc: dnsSvc}
}
