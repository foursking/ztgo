package core

import (
	"github.com/foursking/ztgo/config"
	"github.com/foursking/ztgo/config/env"
	"github.com/foursking/ztgo/log"

	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2/debug/profile/http"
)

var (
	DefaultFlags = []cli.Flag{
		&cli.StringFlag{
			Name:        "qdgo_prometheus_addr",
			EnvVars:     []string{"QDGO_PROMETHEUS_ADDR"},
			Value:       "0.0.0.0:8088",
			Usage:       "prometheus address",
			Destination: &env.PrometheusAddr,
		},
		&cli.StringFlag{
			Name:    "qdgo_config_path",
			EnvVars: []string{"QDGO_CONFIG_PATH"},
			Usage:   "qdgo config file path or directory",
		},
		&cli.StringFlag{
			Name:        "qdgo_deploy_env",
			EnvVars:     []string{"QDGO_DEPLOY_ENV"},
			Value:       env.DefaultDeployEnv,
			Usage:       "deploy env",
			Destination: &env.DeployEnv,
		},
		&cli.StringFlag{
			Name:        "qdgo_pprof_addr",
			Usage:       "pprof address",
			EnvVars:     []string{"QDGO_PPROF_ADDR"},
			Value:       http.DefaultAddress,
			Destination: &http.DefaultAddress,
		},
	}
)

func cliAction(ctx *cli.Context) (err error) {
	if filepath := ctx.String("qdgo_config_path"); filepath != "" {
		if err = config.InitFileConfig(filepath); err != nil {
			log.Fatalf("init file config(%s) error(%v)", filepath, err)
		}
	}
	env.AppName = ctx.String("server_name")
	return
}
