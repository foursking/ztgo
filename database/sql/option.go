package sql

import (
	"time"

	"git.code.oa.com/qdgo/core/log"
	"git.code.oa.com/qdgo/core/net/util/breaker"
)

// Options is sql config
type Options struct {
	DSN          string          // write data source name
	ReadDSN      []string        // read data source name
	Active       int             // pool
	Idle         int             // pool
	IdleTimeout  time.Duration   // connect max life time
	QueryTimeout time.Duration   // query sql timeout
	ExecTimeout  time.Duration   // execute sql timeout
	TranTimeout  time.Duration   // transaction sql timeout
	Breaker      *breaker.Config // breaker
}

// Option is sql option func
type Option func(*Options)

func applyOptions(opts ...Option) *Options {
	options := Options{
		Active:       50,
		Idle:         5,
		IdleTimeout:  10 * time.Minute,
		QueryTimeout: 500 * time.Millisecond,
		ExecTimeout:  500 * time.Millisecond,
		TranTimeout:  500 * time.Millisecond,
	}
	for _, o := range opts {
		o(&options)
	}
	if options.DSN == "" {
		log.Fatalf("mysql must be set dsn, options(%+v)", options)
	}
	if options.QueryTimeout == 0 || options.ExecTimeout == 0 || options.TranTimeout == 0 {
		log.Fatalf("mysql must be set query/execute/transaction timeout, options(%+v)", options)
	}
	return &options
}

func DSN(dsn string) Option {
	return func(opts *Options) {
		opts.DSN = dsn
	}
}

func ReadDSN(dsns []string) Option {
	return func(opts *Options) {
		opts.ReadDSN = dsns
	}
}

func Active(i int) Option {
	return func(opts *Options) {
		opts.Active = i
	}
}

func Idle(i int) Option {
	return func(opts *Options) {
		opts.Active = i
	}
}

func QueryTimeout(d time.Duration) Option {
	return func(opts *Options) {
		opts.QueryTimeout = d
	}
}

func IdleTimeout(d time.Duration) Option {
	return func(opts *Options) {
		opts.IdleTimeout = d
	}
}

func ExecTimeout(d time.Duration) Option {
	return func(opts *Options) {
		opts.ExecTimeout = d
	}
}

func TranTimeout(d time.Duration) Option {
	return func(opts *Options) {
		opts.TranTimeout = d
	}
}

func Breaker(bc *breaker.Config) Option {
	return func(opts *Options) {
		opts.Breaker = bc
	}
}
