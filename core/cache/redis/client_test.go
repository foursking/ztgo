package redis

import (
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
)

var (
	cli *redis.Client
)

func TestMain(m *testing.M) {
	// 9.134.197.160 是 devnet 公用开发机
	var err error
	cli, err = NewClient("tcp://9.134.197.160:6379")
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestSet(t *testing.T) {
	err := cli.Set("name", "www", time.Hour).Err()
	assert.Nil(t, err)
	cmd := cli.Get("name")
	assert.Nil(t, cmd.Err())
	assert.Equal(t, "www", cmd.Val())
}
