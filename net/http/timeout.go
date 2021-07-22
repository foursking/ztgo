package http

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func Timeout(dur time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		newctx, cancel := context.WithTimeout(ctx.Request.Context(), dur)
		ctx.Request = ctx.Request.WithContext(newctx)
		ctx.Next()
		cancel()
	}
}
