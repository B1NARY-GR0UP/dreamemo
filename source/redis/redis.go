// Copyright 2023 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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

// NewSource return a redis source getter
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

// Get from redis source
func (r *Source) Get(ctx context.Context, key string) ([]byte, error) {
	return r.cli.Get(ctx, key).Bytes()
}
