package http

import (
	"cloudflared-tunnel/internal/module/user/ui/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(userRepo *api.Router) *gin.Engine {
	r := gin.Default()

	r.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	userRepo.SetupRoutes(r)

	return r
}
