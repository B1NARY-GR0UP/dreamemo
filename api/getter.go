package api

import "context"

type Getter interface {
	Get(ctx context.Context, key string) ([]byte, error)
}

// GetterFunc uses the same concept as http.HandlerFunc
type GetterFunc func(ctx context.Context, key string) ([]byte, error)

func (f GetterFunc) Get(ctx context.Context, key string) ([]byte, error) {
	return f(ctx, key)
}
