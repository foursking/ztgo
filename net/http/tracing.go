package http

import (
	"context"
	"net/http"

	"github.com/foursking/ztgo/log"
	qdmd "github.com/foursking/ztgo/metadata"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// Tracing is gin server tracing middleware
func Tracing() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var span opentracing.Span
		carrier := opentracing.HTTPHeadersCarrier(ctx.Request.Header)
		parent, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, carrier)
		if err == nil {
			span = opentracing.StartSpan(ctx.Request.URL.Path, opentracing.ChildOf(parent))
		} else {
			span = opentracing.StartSpan(ctx.Request.URL.Path)
		}
		ext.SpanKind.Set(span, qdmd.Server)
		ext.Component.Set(span, qdmd.HTTP)
		ext.HTTPMethod.Set(span, ctx.Request.Method)
		ext.HTTPUrl.Set(span, ctx.Request.URL.String())
		if caller, ok := metadata.Get(ctx.Request.Context(), qdmd.Caller); ok {
			span.SetTag(qdmd.Caller, caller)
		} else {
			span.SetTag(qdmd.Caller, qdmd.Unknown)
		}

		ctx.Request = ctx.Request.WithContext(opentracing.ContextWithSpan(ctx.Request.Context(), span))
		ctx.Next()

		ext.HTTPStatusCode.Set(span, uint16(ctx.Writer.Status()))
		span.Finish()
	}
}

func newClientSpan(ctx context.Context, req *http.Request) (opentracing.Span, context.Context) {
	// gin Context 是自己实现的，它的 Value() 方法实现和标准库不同，需要取其内部ctx传递
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ctx = ginCtx.Request.Context()
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, req.URL.Host+req.URL.Path)
	ext.SpanKind.Set(span, qdmd.Client)
	ext.Component.Set(span, qdmd.HTTP)
	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(span, req.URL.String())
	err := span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		log.Errorf("http client: trace inject http headers error(%v)", err)
	}
	return span, ctx
}
