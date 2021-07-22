package core

import (
	qdgrpc "github.com/foursking/ztgo/net/grpc"
	"github.com/foursking/ztgo/stat"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-plugins/registry/consul/v2"
	opt "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
)

var (
	// DefaultClient is qdgo default client
	DefaultClient = NewClient()
)

// NewClient creates default qdgo client
func NewClient(opts ...client.Option) client.Client {
	opts = append(opts,
		client.Retries(0),
		client.Registry(consul.NewRegistry()),
		client.Wrap(opt.NewClientWrapper(nil)),
		client.Wrap(stat.NewDefaultLatencyWrapper()),
		client.WrapCall(qdgrpc.NewClientMetricHandlerWrapper()),
	)
	return grpc.NewClient(opts...)
}
