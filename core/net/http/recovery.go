package http

import (
	"fmt"
	"net/http/httputil"
	"os"
	"runtime"

	"github.com/foursking/ztgo/log"

	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			var rawReq []byte
			if err := recover(); err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				if ctx.Request != nil {
					rawReq, _ = httputil.DumpRequest(ctx.Request, false)
				}
				msg := fmt.Sprintf("http call panic: %s\n%v\n%s\n", string(rawReq), err, buf)
				_, _ = fmt.Fprintf(os.Stderr, msg)
				log.Error(msg)
				ctx.AbortWithStatus(500)
			}
		}()
		ctx.Next()
	}
}
