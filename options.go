package core

import (
	"context"
	"os"
	"strings"

	"git.code.oa.com/qdgo/core/config/env"
	"git.code.oa.com/qdgo/core/errs"
	"git.code.oa.com/qdgo/core/log"
	"git.code.oa.com/qdgo/core/stat/tracing"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"github.com/micro/go-micro/v2/web"
	"github.com/micro/go-plugins/registry/consul/v2"
	opt "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

func ServiceLogOptions() []micro.Option {
	return []micro.Option{
		micro.AfterStart(func() error {
			log.Infof("%s started, using config(%s)", env.AppName, config.Bytes())
			return nil
		}),
		micro.AfterStop(func() error {
			log.Infof("%s stopped", env.AppName)
			return log.Sync()
		}),
		micro.WrapHandler(func(hf server.HandlerFunc) server.HandlerFunc {
			return func(ctx context.Context, req server.Request, rsp interface{}) error {
				err := hf(ctx, req, rsp)
				log.Debugf("handle method(%s) body(%s) rsp(%+v) error(%v)", req.Method(), req.Body(), rsp, err)
				return err
			}
		}),
	}
}

func WebLogOptions() []web.Option {
	return []web.Option{
		web.AfterStart(func() error {
			log.Infof("%s started, using config(%s)", env.AppName, config.Bytes())
			return nil
		}),
		web.AfterStop(func() error {
			log.Infof("%s stopped", env.AppName)
			return log.Sync()
		}),
	}
}

func ServiceTracingOptions() []micro.Option {
	tracing.MustInitTracer()
	ot := opentracing.GlobalTracer()
	return []micro.Option{
		// for service
		micro.WrapHandler(opt.NewHandlerWrapper(ot)),
		// for client
		micro.WrapClient(opt.NewClientWrapper(ot)),
		// for client call
		micro.WrapCall(opt.NewCallWrapper(ot)),
		// for broker subscriber
		micro.WrapSubscriber(opt.NewSubscriberWrapper(ot)),
		// close tracing after service stops
		micro.AfterStop(tracing.CloseTracer),
	}
}

func WebTracingOptions() []web.Option {
	tracing.MustInitTracer()
	return []web.Option{
		web.AfterStop(tracing.CloseTracer),
	}
}

func ServiceErrorOption() micro.Option {
	return micro.WrapHandler(func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			err := fn(ctx, req, rsp)
			if err == nil {
				return nil
			}
			return errs.FromError(err)
		}
	})
}

func RegistryOption() registry.Registry {
	var addr string
	var opts []registry.Option
	// 可以通过k8s内部路由的环境变量选择直连consul server (默认连本地consul agent)
	e := os.Getenv("CONSUL_SERVER_PORT_8500_TCP")
	if e != "" {
		parts := strings.Split(e, "://")
		if len(parts) == 1 {
			addr = parts[0]
		} else {
			addr = parts[1]
		}
		addr = strings.Trim(addr, "/")
		opts = append(opts, registry.Addrs(addr))
	}
	return consul.NewRegistry(opts...)
}
