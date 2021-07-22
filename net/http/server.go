package http

import (
	"net/http"

	"git.code.oa.com/qdgo/core/config/env"
	"git.code.oa.com/qdgo/core/log"

	"github.com/gin-gonic/gin"
)

// Server http server based on gin
type Server struct {
	*gin.Engine
	*http.Server

	options *ServerOptions
}

// NewServer new http server
func NewServer(opts ...ServerOption) *Server {
	if env.DeployEnv != env.Dev {
		gin.SetMode(gin.ReleaseMode)
	}
	options := DefaultServerOptions
	for _, o := range opts {
		o(&options)
	}
	engine := gin.New()
	engine.Use(Recovery(), Logger(), Tracing())
	for _, m := range options.Middlewares {
		engine.Use(m)
	}
	if len(options.LogIgnorePath) > 0 {
		for _, path := range options.LogIgnorePath {
			ignorePaths[path] = true
		}
	}
	srv := http.Server{
		Addr:         options.Addr,
		Handler:      engine,
		ReadTimeout:  options.ReadTimeout,
		WriteTimeout: options.WriteTimeout,
		IdleTimeout:  options.IdleTimeout,
	}
	return &Server{Engine: engine, Server: &srv, options: &options}
}

// Run runs http serve
func (s *Server) Run() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("http server panic(%v)", err)
		}
	}()
	log.Infof("http start listen addr(%s)", s.options.Addr)
	go func() {
		if err := s.Server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
}

// Ping ping http server
func (s *Server) Ping(ping func(ctx *gin.Context)) {
	s.GET("/health/check", ping)
}
