package http

import (
	"cloudflared-tunnel/internal/config"
	credentialApi "cloudflared-tunnel/internal/module/credential/ui/api"
	userApi "cloudflared-tunnel/internal/module/user/ui/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, userRouter *userApi.Router, credentialRouter *credentialApi.Router) *gin.Engine {
	if cfg.App.Env == "dev" {
		gin.ForceConsoleColor()
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	r.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	userRouter.SetupRoutes(r)
	credentialRouter.SetupRoutes(r)

	return r
}
