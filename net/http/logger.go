package http

import (
	"strconv"
	"strings"
	"time"

	"git.code.oa.com/qdgo/core/log"
	"git.code.oa.com/qdgo/core/metadata"

	"github.com/gin-gonic/gin"
)

const defaultMaxMemory = 32 << 20 // 32 M

// 不需要记录 prom 和 log 的 path
var ignorePaths = map[string]bool{
	"/metrics":      true,
	"/health/check": true,
}

// Logger is middleware for prom and log
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			ip     string
			caller string
			now    = time.Now()
			req    = ctx.Request
			path   = req.URL.Path
		)

		ctype := req.Header.Get("Content-Type")
		switch {
		case strings.Contains(ctype, "multipart/form-data"):
			_ = req.ParseMultipartForm(defaultMaxMemory)
		default:
			_ = req.ParseForm()
		}

		ctx.Next()

		if _, ok := ignorePaths[path]; ok {
			return
		}
		ip = ctx.GetString(metadata.RemoteIP)
		caller = ctx.GetString(metadata.Caller)
		if caller == "" {
			caller = metadata.Unknown
		}
		status := ctx.Writer.Status()
		dt := time.Since(now)
		code := ctx.GetInt(metadata.ErrCode)
		_metricServerReqCodeTotal.Inc(path[1:], caller, req.Method, strconv.FormatInt(int64(code), 10))
		_metricServerReqDur.Observe(int64(dt/time.Millisecond), path[1:], caller, req.Method)
		log.Infof("method(%s) ip(%s) caller(%s) path(%s) params(%s) status(%d) ecode(%d) ua(%s) ts(%f)",
			req.Method, ip, caller, path, req.Form.Encode(), status, code, req.UserAgent(), dt.Seconds())
	}
}
