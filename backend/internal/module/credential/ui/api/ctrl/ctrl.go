package ctrl

import (
	v1 "cloudflared-tunnel/internal/module/credential/ui/api/req/v1"
	"cloudflared-tunnel/internal/middleware"
	"cloudflared-tunnel/internal/module/credential/svc"
	"cloudflared-tunnel/pkg/core"
	"cloudflared-tunnel/pkg/errno"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Ctrl struct {
	CredentialSvc svc.CredentialSvc
}

func NewCtrl(cs svc.CredentialSvc) *Ctrl {
	return &Ctrl{CredentialSvc: cs}
}

func (c *Ctrl) ValidateCredential(ctx *gin.Context) {
	var req v1.ValidateCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.CredentialSvc.ValidateCredential(userID, &req)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

func (c *Ctrl) CreateCredential(ctx *gin.Context) {
	var req v1.CreateCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.CredentialSvc.CreateCredential(userID, &req)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

func (c *Ctrl) GetCredential(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的凭证ID"))
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.CredentialSvc.GetCredential(userID, id)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

func (c *Ctrl) GetCredentials(ctx *gin.Context) {
	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vos, err := c.CredentialSvc.GetCredentials(userID)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vos)
}

func (c *Ctrl) UpdateCredential(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的凭证ID"))
		return
	}

	var req v1.UpdateCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.CredentialSvc.UpdateCredential(userID, id, &req)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

func (c *Ctrl) DeleteCredential(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的凭证ID"))
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	if err := c.CredentialSvc.DeleteCredential(userID, id); err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, nil)
}

func (c *Ctrl) SetDefaultCredential(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的凭证ID"))
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.CredentialSvc.SetDefaultCredential(userID, id)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}

func (c *Ctrl) GetTestLogs(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		core.Fail(ctx, errno.ErrParam.WithMessage("无效的凭证ID"))
		return
	}

	userID := ctx.GetInt64(middleware.ContextKeyUserID)
	vo, err := c.CredentialSvc.GetTestLogs(userID, id)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, vo)
}
