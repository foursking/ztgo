package core

import (
	_ "github.com/micro/go-plugins/broker/kafka/v2"
	_ "github.com/micro/go-plugins/broker/redis/v2"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	//_ "github.com/micro/go-plugins/registry/etcdv3/v2"
	_ "github.com/micro/go-plugins/transport/grpc/v2"
	//_ "github.com/micro/go-plugins/transport/quic/v2"
	_ "github.com/micro/go-plugins/transport/tcp/v2"
)
