package lock

import (
	"context"
	"testing"

	"github.com/foursking/ztgo/cache/redis"

	"github.com/stretchr/testify/assert"
)

func TestLock(t *testing.T) {
	dsn := "tcp://127.0.0.1:6379"
	rds, err := redis.NewClient(dsn)
	assert.Nil(t, err)
	lockKey := "test-lock"
	err = Lock(context.TODO(), rds, lockKey)
	assert.Nil(t, err)
	// 尝试去锁一个已经存在的锁，应错误
	err = Lock(context.TODO(), rds, lockKey)
	assert.Equal(t, ErrLockIsOccupied, err)
	err = Unlock(context.TODO(), rds, lockKey)
	assert.Nil(t, err)
}

func TestUnlock(t *testing.T) {
	dsn := "tcp://127.0.0.1:6379"
	rds, err := redis.NewClient(dsn)
	assert.Nil(t, err)
	lockKey := "test-unlock"
	// 尝试释放一把不存在的锁，应错误
	err = Unlock(context.TODO(), rds, lockKey)
	assert.Equal(t, ErrLockNotFound, err)
	// 锁上再释放，应正确
	err = Lock(context.TODO(), rds, lockKey)
	assert.Nil(t, err)
	err = Unlock(context.TODO(), rds, lockKey)
	assert.Nil(t, err)
}
