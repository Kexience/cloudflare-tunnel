package ctrl

import (
	v1 "cloudflared-tunnel/internal/module/user/ui/api/req/v1"
	"cloudflared-tunnel/pkg/core"
	"cloudflared-tunnel/pkg/errno"

	"github.com/gin-gonic/gin"
)

func (c *Ctrl) CreateUser(ctx *gin.Context) {
	var req v1.CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
		return
	}

	if _, err := c.UserSvc.Register(req.Nickname, req.Username, req.Password, req.Email); err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, nil)
}

func (c *Ctrl) LoginUser(ctx *gin.Context) {
	var req v1.Login
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.Fail(ctx, errno.ErrParam)
	}

	u, err := c.UserSvc.Login(req.Username, req.Password)
	if err != nil {
		core.Fail(ctx, err)
		return
	}
	core.OK(ctx, u)
}
