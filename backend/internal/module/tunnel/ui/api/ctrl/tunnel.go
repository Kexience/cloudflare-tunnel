package ctrl

import (
	"cloudflared-tunnel/internal/middleware"
	v1 "cloudflared-tunnel/internal/module/tunnel/ui/api/req/v1"
	"cloudflared-tunnel/pkg/core"
	"cloudflared-tunnel/pkg/errno"

	"github.com/gin-gonic/gin"
)

// ListTunnels 查询隧道列表
func (c *Ctrl) ListTunnels(ctx *gin.Context) {
	var req struct {
		CredentialID int64 `form:"credential_id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vos, err := c.tunnelSvc.ListTunnels(userID, req.CredentialID)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vos)
}

// GetTunnel 查询隧道详情
func (c *Ctrl) GetTunnel(ctx *gin.Context) {
	tunnelID := ctx.Param("id")
	if tunnelID == "" {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的隧道ID"))
		return
	}

	var req struct {
		CredentialID int64 `form:"credential_id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.tunnelSvc.GetTunnel(userID, req.CredentialID, tunnelID)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

// CreateTunnel 创建隧道
func (c *Ctrl) CreateTunnel(ctx *gin.Context) {
	var req v1.CreateTunnelRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.tunnelSvc.CreateTunnel(userID, req.CredentialID, req.Name)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

// DeleteTunnel 删除隧道
func (c *Ctrl) DeleteTunnel(ctx *gin.Context) {
	tunnelID := ctx.Param("id")
	if tunnelID == "" {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的隧道ID"))
		return
	}

	var req struct {
		CredentialID int64 `form:"credential_id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	if err := c.tunnelSvc.DeleteTunnel(userID, req.CredentialID, tunnelID); err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, nil)
}

// GetTunnelToken 获取隧道连接 Token
func (c *Ctrl) GetTunnelToken(ctx *gin.Context) {
	tunnelID := ctx.Param("id")
	if tunnelID == "" {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的隧道ID"))
		return
	}

	var req struct {
		CredentialID int64 `form:"credential_id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.tunnelSvc.GetTunnelToken(userID, req.CredentialID, tunnelID)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

// ListTunnelConnections 查询隧道连接列表
func (c *Ctrl) ListTunnelConnections(ctx *gin.Context) {
	tunnelID := ctx.Param("id")
	if tunnelID == "" {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的隧道ID"))
		return
	}

	var req struct {
		CredentialID int64 `form:"credential_id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.tunnelSvc.ListTunnelConnections(userID, req.CredentialID, tunnelID)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

// GetTunnelConfig 获取隧道配置
func (c *Ctrl) GetTunnelConfig(ctx *gin.Context) {
	tunnelID := ctx.Param("id")
	if tunnelID == "" {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的隧道ID"))
		return
	}

	var req struct {
		CredentialID int64 `form:"credential_id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.tunnelSvc.GetTunnelConfig(userID, req.CredentialID, tunnelID)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

// UpdateTunnelConfig 更新隧道配置
func (c *Ctrl) UpdateTunnelConfig(ctx *gin.Context) {
	tunnelID := ctx.Param("id")
	if tunnelID == "" {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的隧道ID"))
		return
	}

	var req v1.UpdateTunnelConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.tunnelSvc.UpdateTunnelConfig(userID, req.CredentialID, tunnelID, req.Config)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

// StartTunnel 启动隧道
func (c *Ctrl) StartTunnel(ctx *gin.Context) {
	tunnelID := ctx.Param("id")
	if tunnelID == "" {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的隧道ID"))
		return
	}

	var req struct {
		CredentialID int64 `form:"credential_id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	if err := c.tunnelSvc.StartTunnel(userID, req.CredentialID, tunnelID); err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, nil)
}

// StopTunnel 停止隧道
func (c *Ctrl) StopTunnel(ctx *gin.Context) {
	tunnelID := ctx.Param("id")
	if tunnelID == "" {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的隧道ID"))
		return
	}

	var req struct {
		CredentialID int64 `form:"credential_id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	if err := c.tunnelSvc.StopTunnel(userID, req.CredentialID, tunnelID); err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, nil)
}

// GetTunnelStatus 获取隧道运行状态
func (c *Ctrl) GetTunnelStatus(ctx *gin.Context) {
	tunnelID := ctx.Param("id")
	if tunnelID == "" {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的隧道ID"))
		return
	}

	var req struct {
		CredentialID int64 `form:"credential_id" binding:"required"`
	}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.tunnelSvc.GetTunnelStatus(userID, req.CredentialID, tunnelID)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}
