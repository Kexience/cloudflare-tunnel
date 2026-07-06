package api

import (
	"cloudflared-tunnel/internal/module/user/ui/api/ctrl"

	"github.com/gin-gonic/gin"
)

type Router struct {
	ctrl *ctrl.Ctrl
}

func NewRouter(ctrl *ctrl.Ctrl) *Router {
	return &Router{ctrl: ctrl}
}

func (r *Router) SetupRoutes(g *gin.Engine) {
	user := g.Group("/v1/user")
	{
		user.POST("/", r.ctrl.CreateUser)
	}
}
