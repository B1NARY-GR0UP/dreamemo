package redis

import (
	"context"
	"github.com/B1NARY-GR0UP/dreamemo/api"
)

var _ api.Getter = (*RedisGetter)(nil)

// RedisGetter use redis as datasource
type RedisGetter struct {
}

func (r *RedisGetter) Get(ctx context.Context, key string) ([]byte, error) {
	// TODO: refer to hertz-contrib/cache
	// TODO implement me
	panic("implement me")
}
