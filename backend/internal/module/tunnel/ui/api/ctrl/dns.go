package ctrl

import (
	"cloudflared-tunnel/internal/middleware"
	v1 "cloudflared-tunnel/internal/module/tunnel/ui/api/req/v1"
	"cloudflared-tunnel/pkg/core"
	"cloudflared-tunnel/pkg/errno"

	"github.com/gin-gonic/gin"
)

// ListDNSRecords 查询 DNS 记录列表
func (c *Ctrl) ListDNSRecords(ctx *gin.Context) {
	var req v1.ListDNSRecordsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vos, err := c.dnsSvc.ListDNSRecords(userID, req.CredentialID, req.ZoneID, req.Name, req.Type)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vos)
}

// CreateDNSRecord 创建 DNS 记录
func (c *Ctrl) CreateDNSRecord(ctx *gin.Context) {
	var req v1.CreateDNSRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.dnsSvc.CreateDNSRecord(userID, req.CredentialID, req.ZoneID, req.Name, req.Content, req.Proxied, req.TTL)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

// UpdateDNSRecord 更新 DNS 记录
func (c *Ctrl) UpdateDNSRecord(ctx *gin.Context) {
	recordID := ctx.Param("id")
	if recordID == "" {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的记录ID"))
		return
	}

	var req v1.UpdateDNSRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.dnsSvc.UpdateDNSRecord(userID, req.CredentialID, req.ZoneID, recordID, req.Name, req.Content, req.Proxied, req.TTL)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

// DeleteDNSRecord 删除 DNS 记录
func (c *Ctrl) DeleteDNSRecord(ctx *gin.Context) {
	recordID := ctx.Param("id")
	if recordID == "" {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的记录ID"))
		return
	}

	var req v1.DeleteDNSRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	if err := c.dnsSvc.DeleteDNSRecord(userID, req.CredentialID, req.ZoneID, recordID); err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, nil)
}
