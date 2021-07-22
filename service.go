package core

import (
	"context"
	"os"

	"git.code.oa.com/qdgo/core/log"
	qdgrpc "git.code.oa.com/qdgo/core/net/grpc"
	"git.code.oa.com/qdgo/core/net/ip"
	"git.code.oa.com/qdgo/core/stat/metric"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/debug/profile/http"
	"github.com/micro/go-micro/v2/web"
)

func init() {
	// 获取真实对外网卡IP用于registry
	if os.Getenv("MICRO_SERVER_ADVERTISE") == "" {
		if ips := ip.ExternalIP(); len(ips) > 0 {
			_ = os.Setenv("MICRO_SERVER_ADVERTISE", ips[0]+os.Getenv("MICRO_SERVER_ADDRESS"))
		}
	}
}

// NewGrpcService creates qdgo grpc service
func NewGrpcService(opts ...micro.Option) micro.Service {
	srv := micro.NewService(
		micro.Registry(RegistryOption()),
		micro.Context(context.Background()),
		micro.Flags(DefaultFlags...),
		micro.Action(cliAction),
		micro.WrapHandler(qdgrpc.NewServerMetricHandlerWrapper()),
		ServiceErrorOption(),
	)
	srv.Init()
	opts = append(opts, ServiceLogOptions()...)
	opts = append(opts, ServiceTracingOptions()...)
	srv.Init(opts...)
	metric.BootPrometheus()
	return srv
}

// NewWebService creates qdgo web service
func NewWebService(opts ...web.Option) web.Service {
	profile := http.NewProfile()
	srv := web.NewService(
		web.Registry(RegistryOption()),
		web.Context(context.Background()),
		web.Flags(DefaultFlags...),
		web.Action(func(ctx *cli.Context) { _ = cliAction(ctx) }),
		web.BeforeStart(func() error { return profile.Start() }),
		web.AfterStop(func() error { return profile.Stop() }),
	)
	_ = srv.Init()
	opts = append(opts, WebLogOptions()...)
	opts = append(opts, WebTracingOptions()...)
	if err := srv.Init(opts...); err != nil {
		log.Fatalf("NewWeb srv Init(%+v) error(%v)", opts, err)
	}
	metric.BootPrometheus()
	return srv
}

// NewUDPService creates qdgo udp service
func NewUDPService(opts ...micro.Option) micro.Service {
	srv := micro.NewService(
		micro.Context(context.Background()),
		micro.Flags(DefaultFlags...),
		micro.Action(cliAction),
	)
	srv.Init(opts...)
	srv.Init(
		micro.Registry(RegistryOption()),
		micro.WrapHandler(qdgrpc.NewServerMetricHandlerWrapper()),
		ServiceErrorOption(),
	)
	srv.Init(ServiceLogOptions()...)
	srv.Init(ServiceTracingOptions()...)
	metric.BootPrometheus()
	return srv
}
