package tracing

import (
	"io"
	"sync"

	"git.code.oa.com/qdgo/core/config/env"
	"git.code.oa.com/qdgo/core/log"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var (
	tracerOnce   sync.Once
	tracerCloser io.Closer
)

func MustInitTracer() {
	tracerOnce.Do(func() {
		cfg, err := config.FromEnv()
		if err != nil {
			log.Errorf("jaeger: get config from env error(%v)", err)
			return
		}
		// 如果未指定采样方式，默认按概率采样 50%
		if cfg.Sampler.Type == "" {
			cfg.Sampler = &config.SamplerConfig{
				Type:  jaeger.SamplerTypeProbabilistic,
				Param: 0.5,
			}
		}
		if tracerCloser, err = cfg.InitGlobalTracer(env.AppNameWithEnv()); err != nil {
			log.Errorf("jaeger: init tracer error(%v)", err)
			return
		}
		log.Infof("jaeger: init tracer with disable(%v) sampler(%+v)", cfg.Disabled, cfg.Sampler)
	})
}

func CloseTracer() error {
	if tracerCloser == nil {
		return nil
	}
	return tracerCloser.Close()
}
