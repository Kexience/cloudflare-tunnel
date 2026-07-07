package middleware

import (
	"net/http"
	"sync"
	"time"

	"cloudflared-tunnel/pkg/errno"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	count    int
	lastSeen time.Time
}

type RateLimiter struct {
	visitors sync.Map
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limit:  limit,
		window: window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window * 2)
	defer ticker.Stop()
	for range ticker.C {
		rl.visitors.Range(func(key, value any) bool {
			v := value.(*visitor)
			if time.Since(v.lastSeen) > rl.window*2 {
				rl.visitors.Delete(key)
			}
			return true
		})
	}
}

func (rl *RateLimiter) isAllowed(key string) bool {
	val, _ := rl.visitors.LoadOrStore(key, &visitor{count: 0, lastSeen: time.Now()})
	v := val.(*visitor)

	now := time.Now()
	if now.Sub(v.lastSeen) > rl.window {
		v.count = 0
		v.lastSeen = now
	}

	if v.count >= rl.limit {
		return false
	}

	v.count++
	v.lastSeen = now
	return true
}

// RateLimit 返回基于客户端 IP 的速率限制中间件
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	rl := NewRateLimiter(limit, window)
	return func(c *gin.Context) {
		key := c.ClientIP()
		if !rl.isAllowed(key) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    errno.ErrParam.Code,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
