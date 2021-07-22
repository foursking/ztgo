package http

import (
	"time"

	"git.code.oa.com/qdgo/core/net/util/breaker"
)

// ClientOptions is http client options
type ClientOptions struct {
	DialTimeout time.Duration
	Timeout     time.Duration
	KeepAlive   time.Duration
	Breaker     *breaker.Config
	URL         map[string]*ClientOptions
	Host        map[string]*ClientOptions
}

// ClientOption changes ClientOptions
type ClientOption func(*ClientOptions)

// DefaultClientOptions default client options
var DefaultClientOptions = ClientOptions{
	DialTimeout: 500 * time.Millisecond,
	Timeout:     time.Second,
	KeepAlive:   time.Minute,
}

// ClientDialTimeout sets client dial timeout
func ClientDialTimeout(t time.Duration) ClientOption {
	return func(options *ClientOptions) {
		if t == 0 {
			t = DefaultClientOptions.DialTimeout
		}
		options.DialTimeout = t
	}
}

// ClientTimeout sets client timeout
func ClientTimeout(t time.Duration) ClientOption {
	return func(options *ClientOptions) {
		if t == 0 {
			t = DefaultClientOptions.Timeout
		}
		options.Timeout = t
	}
}

// ClientKeepAlive sets client keepAlive duration
func ClientKeepAlive(t time.Duration) ClientOption {
	return func(options *ClientOptions) {
		if t == 0 {
			t = DefaultClientOptions.KeepAlive
		}
		options.KeepAlive = t
	}
}
