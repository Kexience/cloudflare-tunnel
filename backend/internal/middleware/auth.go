package middleware

import (
	"errors"
	"strings"

	"cloudflared-tunnel/pkg/core"
	"cloudflared-tunnel/pkg/errno"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyUserID   = "user_id"
	ContextKeyUsername = "username"
)

func Auth(jwt *core.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			core.Fail(ctx, errno.ErrUnauthorized)
			ctx.Abort()
			return
		}

		tokenString, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found || tokenString == "" {
			core.Fail(ctx, errno.ErrUnauthorized)
			ctx.Abort()
			return
		}

		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			if errors.Is(err, core.ErrTokenExpired) {
				core.Fail(ctx, errno.ErrUnauthorized)
			} else {
				core.Fail(ctx, errno.ErrTokenInvalid)
			}
			ctx.Abort()
			return
		}

		ctx.Set(ContextKeyUserID, claims.UserID)
		ctx.Set(ContextKeyUsername, claims.Username)
		ctx.Next()
	}
}
