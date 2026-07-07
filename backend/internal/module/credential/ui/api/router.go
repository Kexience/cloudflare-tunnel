package api

import (
	"cloudflared-tunnel/internal/middleware"
	"cloudflared-tunnel/internal/module/credential/ui/api/ctrl"
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
	authorized := v1.Group("/credentials")
	authorized.Use(middleware.Auth(r.jwt))
	{
		authorized.POST("", r.ctrl.CreateCredential)
		authorized.GET("", r.ctrl.GetCredentials)
		authorized.GET("/:id", r.ctrl.GetCredential)
		authorized.PUT("/:id", r.ctrl.UpdateCredential)
		authorized.DELETE("/:id", r.ctrl.DeleteCredential)
		authorized.PUT("/:id/default", r.ctrl.SetDefaultCredential)
	}
}
