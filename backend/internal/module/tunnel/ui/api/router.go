package api

import (
	"cloudflared-tunnel/internal/middleware"
	"cloudflared-tunnel/internal/module/tunnel/ui/api/ctrl"
	"cloudflared-tunnel/pkg/core"

	"github.com/gin-gonic/gin"
)

// Router 隧道管理路由
type Router struct {
	ctrl *ctrl.Ctrl
	jwt  *core.JWT
}

// NewRouter 创建隧道管理路由
func NewRouter(ctrl *ctrl.Ctrl, jwt *core.JWT) *Router {
	return &Router{ctrl: ctrl, jwt: jwt}
}

// SetupRoutes 注册路由
func (r *Router) SetupRoutes(g *gin.Engine) {
	v1 := g.Group("/v1")

	// 隧道管理
	tunnels := v1.Group("/tunnels")
	tunnels.Use(middleware.Auth(r.jwt))
	{
		tunnels.GET("", r.ctrl.ListTunnels)
		tunnels.POST("", r.ctrl.CreateTunnel)
		tunnels.GET("/:id", r.ctrl.GetTunnel)
		tunnels.DELETE("/:id", r.ctrl.DeleteTunnel)
		tunnels.GET("/:id/token", r.ctrl.GetTunnelToken)
		tunnels.GET("/:id/connections", r.ctrl.ListTunnelConnections)
		tunnels.GET("/:id/config", r.ctrl.GetTunnelConfig)
		tunnels.PUT("/:id/config", r.ctrl.UpdateTunnelConfig)
		tunnels.POST("/:id/start", r.ctrl.StartTunnel)
		tunnels.POST("/:id/stop", r.ctrl.StopTunnel)
		tunnels.GET("/:id/status", r.ctrl.GetTunnelStatus)
	}

	// DNS 记录管理
	dns := v1.Group("/dns")
	dns.Use(middleware.Auth(r.jwt))
	{
		dns.GET("/records", r.ctrl.ListDNSRecords)
		dns.POST("/records", r.ctrl.CreateDNSRecord)
		dns.PUT("/records/:id", r.ctrl.UpdateDNSRecord)
		dns.DELETE("/records/:id", r.ctrl.DeleteDNSRecord)
	}
}
