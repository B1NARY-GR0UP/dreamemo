package redis

import (
	"context"

	"github.com/B1NARY-GR0UP/dreamemo/source"
	"github.com/redis/go-redis/v9"
)

var _ source.Getter = (*Source)(nil)

// Source use redis as datasource
type Source struct {
	cli *redis.Client
}

func NewSource(opts ...Option) *Source {
	redisOpts := &redis.Options{
		Addr:     ":6379",
		Password: "",
		DB:       0,
	}
	for _, opt := range opts {
		opt(redisOpts)
	}
	rdb := redis.NewClient(redisOpts)
	return &Source{
		cli: rdb,
	}
}

func (r *Source) Get(ctx context.Context, key string) ([]byte, error) {
	return r.cli.Get(ctx, key).Bytes()
}
