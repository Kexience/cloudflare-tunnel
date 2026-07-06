package api

import (
	"cloudflared-tunnel/internal/middleware"
	"cloudflared-tunnel/internal/module/user/ui/api/ctrl"
	"cloudflared-tunnel/pkg/core"

	"github.com/gin-gonic/gin"
)

type Router struct {
	ctrl *ctrl.Ctrl
	jwt  *core.JWT
}

func NewRouter(ctrl *ctrl.Ctrl, jwt *core.JWT) *Router {
	return &Router{ctrl: ctrl, jwt: jwt}
}

func (r *Router) SetupRoutes(g *gin.Engine) {
	v1 := g.Group("/v1")
	{
		// 公开路由
		user := v1.Group("/user")
		{
			user.POST("/register", r.ctrl.CreateUser)
		}

		// 需要鉴权的路由
		authorized := v1.Group("/user")
		authorized.Use(middleware.Auth(r.jwt))
		{
			authorized.GET("/me", r.ctrl.GetCurrentUser)
		}
	}
}
