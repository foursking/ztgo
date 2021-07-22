package udp

import (
	"context"

	"github.com/micro/go-micro/v2/broker"
	"github.com/micro/go-micro/v2/codec"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
)

type (
	handleFuncKey     struct{}
	readBufKey        struct{}
	readerNumKey      struct{}
	handlerNumKey     struct{}
	writerNumKey      struct{}
	handleChanSizeKey struct{}
	writeChanSizeKey  struct{}
)

func newOptions(opts ...server.Option) server.Options {
	options := server.Options{
		Codecs:   make(map[string]codec.NewCodec),
		Metadata: map[string]string{},
		Context:  context.Background(),
	}
	for _, o := range opts {
		o(&options)
	}
	if options.Broker == nil {
		options.Broker = broker.DefaultBroker
	}
	if options.Registry == nil {
		options.Registry = registry.DefaultRegistry
	}
	if len(options.Address) == 0 {
		options.Address = server.DefaultAddress
	}
	if len(options.Name) == 0 {
		options.Name = server.DefaultName
	}
	if len(options.Id) == 0 {
		options.Id = server.DefaultId
	}
	if len(options.Version) == 0 {
		options.Version = server.DefaultVersion
	}
	return options
}

func HandleFunc(fn UDPHandleFunc) server.Option {
	return func(opts *server.Options) {
		opts.Context = context.WithValue(opts.Context, handleFuncKey{}, fn)
	}
}

func ReadBufferSize(size int) server.Option {
	return func(opts *server.Options) {
		opts.Context = context.WithValue(opts.Context, readBufKey{}, size)
	}
}

func ReaderNum(num int) server.Option {
	return func(opts *server.Options) {
		opts.Context = context.WithValue(opts.Context, readerNumKey{}, num)
	}
}

func HandlerNum(num int) server.Option {
	return func(opts *server.Options) {
		opts.Context = context.WithValue(opts.Context, handlerNumKey{}, num)
	}
}

func WriterNum(num int) server.Option {
	return func(opts *server.Options) {
		opts.Context = context.WithValue(opts.Context, writerNumKey{}, num)
	}
}

func HandleChanSize(size int) server.Option {
	return func(opts *server.Options) {
		opts.Context = context.WithValue(opts.Context, handleChanSizeKey{}, size)
	}
}

func WriteChanSize(size int) server.Option {
	return func(opts *server.Options) {
		opts.Context = context.WithValue(opts.Context, writeChanSizeKey{}, size)
	}
}
