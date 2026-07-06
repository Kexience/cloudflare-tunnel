package ctrl

import (
	"strconv"

	"cloudflared-tunnel/internal/middleware"
	"cloudflared-tunnel/pkg/core"
	"cloudflared-tunnel/pkg/errno"

	"github.com/gin-gonic/gin"
)

func (c *Ctrl) GetCurrentUser(ctx *gin.Context) {
	userID, exists := ctx.Get(middleware.ContextKeyUserID)
	if !exists {
		core.Fail(ctx, errno.ErrUnauthorized)
		return
	}

	user, err := c.UserSvc.GetUserByID(strconv.FormatInt(userID.(int64), 10))
	if err != nil {
		core.Fail(ctx, err)
		return
	}

	core.OK(ctx, user)
}
