package lock

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v7"
)

var (
	// ErrLockIsOccupied 锁被占用中
	ErrLockIsOccupied = errors.New("lock is occupied")

	// ErrLockNotFound 在尝试释放一把不存在的锁
	ErrLockNotFound = errors.New("lock not found")
)

// Lock represents a distribute lock by redis `setnx`
func Lock(ctx context.Context, rds *redis.Client, name string, opts ...OptionFunc) error {
	options := DefaultOptions
	for _, o := range opts {
		o(&options)
	}
	success, err := rds.WithContext(ctx).SetNX(name, 1, options.Expiration).Result()
	if err != nil {
		return err
	}
	if !success {
		return ErrLockIsOccupied
	}
	return nil
}

// Unlock releases a lock
func Unlock(ctx context.Context, rds *redis.Client, name string) error {
	deleted, err := rds.WithContext(ctx).Del(name).Result()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return ErrLockNotFound
	}
	return nil
}
