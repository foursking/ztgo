package http

import (
	"time"

	"github.com/gin-gonic/gin"
)

// ServerOptions server options
type ServerOptions struct {
	Addr          string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	IdleTimeout   time.Duration
	Middlewares   []gin.HandlerFunc
	LogIgnorePath []string
}

// ServerOption changes ServerOptions
type ServerOption func(*ServerOptions)

// DefaultServerOptions default server options
var DefaultServerOptions = ServerOptions{
	Addr:         "0.0.0.0:8080",
	ReadTimeout:  time.Second,
	WriteTimeout: time.Second,
	IdleTimeout:  10 * time.Minute,
}

// ServerAddr indicates server address
func ServerAddr(addr string) ServerOption {
	return func(options *ServerOptions) {
		if addr == "" {
			addr = DefaultServerOptions.Addr
		}
		options.Addr = addr
	}
}

// ServerReadTimeout sets read timeout
func ServerReadTimeout(t time.Duration) ServerOption {
	return func(options *ServerOptions) {
		if t == 0 {
			t = DefaultServerOptions.ReadTimeout
		}
		options.ReadTimeout = t
	}
}

// ServerWriteTimeout sets write timeout
func ServerWriteTimeout(t time.Duration) ServerOption {
	return func(options *ServerOptions) {
		if t == 0 {
			t = DefaultServerOptions.WriteTimeout
		}
		options.WriteTimeout = t
	}
}

// ServerIdleTimeout sets idle timeout
func ServerIdleTimeout(t time.Duration) ServerOption {
	return func(options *ServerOptions) {
		if t == 0 {
			t = DefaultServerOptions.IdleTimeout
		}
		options.IdleTimeout = t
	}
}

// ServerMiddlewares sets server middlewares
func ServerMiddlewares(middlewares ...gin.HandlerFunc) ServerOption {
	return func(options *ServerOptions) {
		for _, m := range middlewares {
			options.Middlewares = append(options.Middlewares, m)
		}
	}
}

// LogIgnorePath path is the ignore list will hidden in the logs
func LogIgnorePath(path ...string) ServerOption {
	return func(options *ServerOptions) {
		options.LogIgnorePath = append(options.LogIgnorePath, path...)
	}
}
