package loadbalance

import "context"

type LoadBalancer interface {
	Pick(key string) (Instance, bool)
}

type Instance interface {
	Get(ctx context.Context) error
}
