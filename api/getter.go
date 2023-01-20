package api

import "context"

// Getter specifies how to get data from a datasource
type Getter interface {
	Get(ctx context.Context, key string) ([]byte, error)
}

// GetterFunc uses the same concept as http.HandlerFunc
type GetterFunc func(ctx context.Context, key string) ([]byte, error)

// Get dat from datasource
func (f GetterFunc) Get(ctx context.Context, key string) ([]byte, error) {
	return f(ctx, key)
}
