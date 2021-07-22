package redis

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"git.code.oa.com/qdgo/core/config/dsn"

	"github.com/go-redis/redis/v7"
)

type dsnOpt struct {
	Scheme       string        `dsn:"scheme"`
	Address      string        `dsn:"address"`
	Password     string        `dsn:"password"`
	DialTimeout  time.Duration `dsn:"query.dialTimeout,500ms"`
	ReadTimeout  time.Duration `dsn:"query.readTimeout,1s"`
	WriteTimeout time.Duration `dsn:"query.writeTimeout,1s"`
	PoolSize     int           `dsn:"query.poolSize,50"`
	MinIdle      int           `dsn:"query.minIdle,10"`
	IdleTimeout  time.Duration `dsn:"query.idleTimeout,10m"`
	KeepAlive    time.Duration `dsn:"query.keepAlive,10m"`
}

// Options gets redis options
func Options(dsnstr string) (*redis.Options, error) {
	d, err := dsn.Parse(dsnstr)
	if err != nil {
		return nil, err
	}
	o := dsnOpt{}
	if err = d.Bind(&o); err != nil {
		return nil, err
	}
	opts := redis.Options{
		Network:      o.Scheme,
		Addr:         o.Address,
		Password:     o.Password,
		DialTimeout:  o.DialTimeout,
		ReadTimeout:  o.ReadTimeout,
		WriteTimeout: o.WriteTimeout,
		PoolSize:     o.PoolSize,
		MinIdleConns: o.MinIdle,
		IdleTimeout:  o.IdleTimeout,
	}
	return &opts, nil
}
