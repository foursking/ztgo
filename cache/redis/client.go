package redis

import "github.com/go-redis/redis/v7"

// NewClient create redis client
func NewClient(dsn string) (*redis.Client, error) {
	opts, err := Options(dsn)
	if err != nil {
		return nil, err
	}
	cli := redis.NewClient(opts)
	cli.AddHook(&clientHook{})
	return cli, err
}
