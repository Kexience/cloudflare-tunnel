package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.Any("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	return r
}
