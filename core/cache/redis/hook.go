package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/opentracing/opentracing-go"
)

type metricCtxKey struct{}

type clientHook struct{}

func (h *clientHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	var key string
	if len(cmd.Args()) > 1 {
		key, _ = cmd.Args()[1].(string)
	}
	// for metrics
	ctx = context.WithValue(ctx, metricCtxKey{}, time.Now())
	// for tracing
	op := "Redis: " + cmd.Name()
	span, ctx := opentracing.StartSpanFromContext(ctx, op, opentracing.Tags{"key": key})
	ctx = opentracing.ContextWithSpan(ctx, span)
	return ctx, nil
}

func (h *clientHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	var key string
	if len(cmd.Args()) > 1 {
		key, _ = cmd.Args()[1].(string)
	}
	// for metrics
	if t, ok := ctx.Value(metricCtxKey{}).(time.Time); ok {
		_metricReqDur.Observe(int64(time.Since(t)/time.Millisecond), cmd.Name(), key)
	}
	if cmd.Err() == nil {
		_metricHits.Inc(cmd.Name(), key)
	} else {
		if cmd.Err() == redis.Nil {
			_metricMisses.Inc(cmd.Name(), key)
		} else {
			_metricReqErr.Inc(cmd.Name(), key, cmd.Err().Error())
		}
	}
	// for tracing
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.Finish()
	}
	return nil
}

func (h *clientHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	// TODO add metrics & tracing
	return ctx, nil
}

func (h *clientHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	// TODO add metrics & tracing
	return nil
}
