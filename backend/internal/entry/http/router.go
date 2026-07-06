package http

import (
	"cloudflared-tunnel/internal/config"
	"cloudflared-tunnel/internal/module/user/ui/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config, userRepo *api.Router) *gin.Engine {
	if cfg.App.Env == "dev" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.Default()

	r.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	userRepo.SetupRoutes(r)

	return r
}
