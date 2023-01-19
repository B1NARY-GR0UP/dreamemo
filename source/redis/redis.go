package redis

import (
	"context"
)

// RedisGetter use redis as datasource
type RedisGetter struct {
}

func (r *RedisGetter) Get(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}
