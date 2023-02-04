package loadbalance

import (
	"context"
	"github.com/B1NARY-GR0UP/dreamemo/protocol"
)

type LoadBalancer interface {
	Pick(key string) (Instance, bool)
}

type Instance interface {
	Get(ctx context.Context, in protocol.GetRequest, out protocol.GetResponse) error
}
