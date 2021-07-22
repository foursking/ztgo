package env

import (
	"os"

	"git.code.oa.com/qdgo/core/net/ip"
)

const (
	// Dev development env
	Dev string = "dev"
	// OA test env
	OA string = "oa"
	// Gray gray env
	Gray string = "gray"
	// OL production env
	OL string = "ol"
)

var (
	// DefaultDeployEnv is default deploy env
	DefaultDeployEnv = Dev

	// AppName app name
	AppName string
	// ServerIP server ip
	ServerIP string
	// Hostname server hostname
	Hostname string
	// DeployEnv deploy env
	DeployEnv string
	// PrometheusAddr is prometheus http addr
	PrometheusAddr string
)

func init() {
	ServerIP = ip.InternalIP()
	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		Hostname = hostname
	}
}

// Env gets environment value, def will be return when key is not set
func Env(name string, defaults ...string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// AppNameWithEnv returns app name with deploy env
func AppNameWithEnv() string {
	return AppName + "." + DeployEnv
}
